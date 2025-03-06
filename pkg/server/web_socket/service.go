package web_socket

import (
	"context"
	"github.com/skolldire/web-simplify/pkg/utilities/log"
	"net/http"
	"strings"

	"github.com/coder/websocket"
)

var _ Service = (*server)(nil)

func NewServer(l log.Service) Service {
	return &server{
		clients: make(map[*Client]bool),
		log:     l,
	}
}

func (s *server) SetProcessingFunc(f ProcessingFunc) {
	s.processor = f
}

func (s *server) HandleNewConnection(w http.ResponseWriter, r *http.Request) {
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

func (s *server) listenForMessages(c *Client) {
	defer func() {
		s.lock.Lock()
		delete(s.clients, c)
		s.lock.Unlock()
		_ = c.Conn.Close(websocket.StatusNormalClosure, "Closing connection")
		s.log.Info(c.Ctx, "Client disconnected", map[string]interface{}{"ip": c.IP})
	}()

	for {
		_, message, err := c.Conn.Read(c.Ctx)
		if err != nil {
			s.log.Error(c.Ctx, err, "error to listen client", nil)
			break
		}
		if s.processor != nil {
			processedMessage, err := s.processor(message)
			if err != nil {
				s.log.Error(c.Ctx, err, "error to receive message", map[string]interface{}{"ip": c.IP, "message": string(message)})
				continue
			}
			s.sendMessageToClient(c, processedMessage)
		}
	}
}

func (s *server) sendMessageToClient(client *Client, message []byte) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if err := client.Conn.Write(client.Ctx, websocket.MessageText, message); err != nil {
		s.log.Error(client.Ctx, err, "error to send message", nil)
		_ = client.Conn.Close(websocket.StatusInternalError, "Error in sending message")
		delete(s.clients, client)
	}
}
