package tcp

import (
	"context"
	"fmt"
	"github.com/skolldire/web-simplify/pkg/utilities/log"
	"net"
)

type service struct {
	server net.Listener
	log    log.Service
	port   string
}

var _ Service = (*service)(nil)

func NewService(d Dependencies) *service {
	listener, err := net.Listen("tcp", ":"+d.Config.Port)
	if err != nil {
		d.Log.Error(context.Background(), err, map[string]interface{}{
			"port": d.Config.Port, "message": "Error starting server"})
		return nil
	}
	d.Log.Info(context.Background(), fmt.Sprintf("%s Server listening", d.Config.InstanceName),
		map[string]interface{}{"port": d.Config.Port})

	return &service{
		server: listener,
		log:    d.Log,
		port:   d.Config.Port,
	}
}

func (s *service) GetMessage(f ProcessingFunc) {
	for {
		conn, err := s.server.Accept()
		if err != nil {
			s.log.Error(context.Background(), err,
				map[string]interface{}{"port": s.port, "message": "Error accepting connection"})
			continue
		}
		s.log.Info(context.Background(), "New connection accepted",
			map[string]interface{}{"port": s.port, "client": conn.RemoteAddr().String()})
		go s.handleConnection(conn, f)
	}
}

func (s *service) handleConnection(conn net.Conn, f ProcessingFunc) {
	defer func() {
		err := conn.Close()
		if err != nil {
			s.log.Error(context.Background(), err,
				map[string]interface{}{"message": "Error closing connection"})
		}
	}()

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err.Error() == "EOF" {
				s.log.Info(context.Background(), "Client disconnected",
					map[string]interface{}{"client": conn.RemoteAddr().String()})
				break
			}
			s.log.Error(context.Background(), err,
				map[string]interface{}{"message": "Error reading from connection"})
			break
		}
		message := string(buf[:n])
		response, err := f(message)
		if err != nil {
			s.log.Error(context.Background(), err,
				map[string]interface{}{"message": "Error processing message"})
			continue
		}
		_, err = conn.Write([]byte(response))
		if err != nil {
			s.log.Error(context.Background(), err,
				map[string]interface{}{"message": "Error writing response"})
			break
		}
	}
}
