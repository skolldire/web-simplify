package web_socket

import (
	"context"
	"github.com/skolldire/web-simplify/pkg/utilities/log"
	"net/http"
	"sync"

	"github.com/coder/websocket"
)

type Service interface {
	SetProcessingFunc(f ProcessingFunc)
	HandleNewConnection(w http.ResponseWriter, r *http.Request)
}

type ProcessingFunc func(msg []byte) ([]byte, error)

type server struct {
	clients   map[*Client]bool
	lock      sync.Mutex
	log       log.Service
	processor ProcessingFunc
}

type Client struct {
	Conn   *websocket.Conn
	Ctx    context.Context
	Cancel context.CancelFunc
	IP     string
}
