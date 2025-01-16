package tcp

import (
	"context"
	"fmt"
	"github.com/skolldire/web-simplify/pkg/utilities/log"
	"net"
	"sync"
)

type service struct {
	server net.Listener
	log    log.Service
	port   string
}

var _ Service = (*service)(nil)

func NewService(d Dependencies) *service {
	s := initializeService(d.Config, d.Log)
	return &service{
		server: *s,
		log:    d.Log,
		port:   d.Config.Port,
	}
}

func initializeService(c Config, l log.Service) *net.Listener {
	var wg sync.WaitGroup
	defer wg.Done()
	listener, err := net.Listen("tcp", ":"+c.Port)
	if err != nil {
		l.Error(context.Background(), err, map[string]interface{}{
			"port": c.Port, "message": "Error starting server"})
		return nil
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			l.Error(context.Background(), err,
				map[string]interface{}{"port": c.Port,
					"message": "Error closing listener"})
		}
	}(listener)
	l.Info(context.Background(), fmt.Sprintf("%s Server listening", c.InstanceName),
		map[string]interface{}{"port": c.Port})
	return &listener
}

func (s *service) GetMessage(f ProcessingFunc) {
	for {
		conn, err := s.server.Accept()
		if err != nil {
			s.log.Error(context.Background(), err,
				map[string]interface{}{"port": s.port,
					"message": "Error accepting connection",
					"client":  conn.RemoteAddr().String()})
			continue
		}
		s.log.Info(context.Background(), "New connection accepted on port",
			map[string]interface{}{"port": s.port,
				"client": conn.RemoteAddr().String()})
		go s.handleConnection(conn, f)
	}
}

func (s *service) handleConnection(conn net.Conn, f ProcessingFunc) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			s.log.Error(context.Background(), err,
				map[string]interface{}{"message": "Error closing connection"})
			return
		}
	}(conn)
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			s.log.Error(context.Background(), err,
				map[string]interface{}{"message": "Error reading from connection"})
		}
		message := string(buf[:n])
		response, err := f(message)
		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Error responding:", err)
			return
		}
	}
}
