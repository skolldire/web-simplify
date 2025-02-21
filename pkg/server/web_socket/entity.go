package web_socket

import (
	"context"
	"github.com/skolldire/web-simplify/pkg/utilities/log"
	"sync"

	"github.com/coder/websocket"
)

type ProcessingFunc func(msg []byte) ([]byte, error)

type Server struct {
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
