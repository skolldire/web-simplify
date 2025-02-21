package web_socket

import (
	"context"
	"github.com/skolldire/web-simplify/pkg/utilities/log"
	"net/http"
	"strings"

	"github.com/coder/websocket"
)

func NewServer(l log.Service) *Server {
	return &Server{
		clients: make(map[*Client]bool),
		log:     l,
	}
}

func (s *Server) SetProcessingFunc(f ProcessingFunc) {
	s.processor = f
}

func (s *Server) HandleNewConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: true})
	if err != nil {
		s.log.Info(r.Context(), "Failed to accept connection", map[string]interface{}{"error": err.Error()})
		return
	}
	ctx, cancel := context.WithCancel(r.Context())
	clientIP := strings.Split(r.RemoteAddr, ":")[0]
	c := &Client{Conn: conn, Ctx: ctx, Cancel: cancel, IP: clientIP}

	s.lock.Lock()
	s.clients[c] = true
	s.lock.Unlock()

	s.log.Info(ctx, "New client connected", map[string]interface{}{"ip": clientIP})
	go s.listenForMessages(c)
}

func (s *Server) listenForMessages(c *Client) {
	defer func() {
		s.lock.Lock()
		delete(s.clients, c)
		s.lock.Unlock()
		c.Conn.Close(websocket.StatusNormalClosure, "Closing connection")
		s.log.Info(c.Ctx, "Client disconnected", map[string]interface{}{"ip": c.IP})
	}()

	for {
		_, message, err := c.Conn.Read(c.Ctx)
		if err != nil {
			s.log.Error(c.Ctx, err, nil)
			break
		}
		s.log.Info(c.Ctx, "Received message", map[string]interface{}{"ip": c.IP, "message": string(message)})

		if s.processor != nil {
			processedMessage, err := s.processor(message)
			if err != nil {
				s.log.Error(c.Ctx, err, map[string]interface{}{"ip": c.IP, "message": string(message)})
				continue
			}
			s.sendMessageToClient(c, processedMessage)
		}
	}
}

func (s *Server) sendMessageToClient(client *Client, message []byte) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if err := client.Conn.Write(client.Ctx, websocket.MessageText, message); err != nil {
		s.log.Error(client.Ctx, err, nil)
		client.Conn.Close(websocket.StatusInternalError, "Error in sending message")
		delete(s.clients, client)
	}
}
