package tcp

import (
	"context"
	"fmt"
	"net"
)

var _ Service = (*service)(nil)

func NewService(d Dependencies) Service {
	listener, err := net.Listen("tcp", ":"+d.Config.Port)
	if err != nil {
		d.Log.Error(context.Background(), err, "Error starting server", map[string]interface{}{
			"port": d.Config.Port})
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

func (s *service) Start(ctx context.Context, f ProcessingFunc) {
	go s.getMessage(ctx, f)
}

func (s *service) getMessage(ctx context.Context, f ProcessingFunc) {
	for {
		select {
		case <-ctx.Done():
			s.log.Info(context.Background(), "Shutting down TCP server",
				map[string]interface{}{"port": s.port})
			s.server.Close()
			return
		default:
			conn, err := s.server.Accept()
			if err != nil {
				s.log.Error(context.Background(), err, "Error accepting connection",
					map[string]interface{}{"port": s.port})
				continue
			}
			s.log.Info(context.Background(), "New connection accepted",
				map[string]interface{}{"port": s.port, "client": conn.RemoteAddr().String()})
			go s.handleConnection(conn, f)
		}
	}
}

func (s *service) handleConnection(conn net.Conn, f ProcessingFunc) {
	defer func() {
		err := conn.Close()
		if err != nil {
			s.log.Error(context.Background(), err, "Error closing connection", nil)
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
			s.log.Error(context.Background(), err, "Error reading from connection", nil)
			break
		}
		message := string(buf[:n])
		response, err := f(message)
		if err != nil {
			s.log.Error(context.Background(), err, "Error processing message", nil)
			continue
		}
		_, err = conn.Write([]byte(response))
		if err != nil {
			s.log.Error(context.Background(), err, "Error writing response", nil)
			break
		}
	}
}
