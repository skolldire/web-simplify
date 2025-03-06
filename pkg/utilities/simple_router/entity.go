package simple_router

import (
	"github.com/go-chi/chi/v5"
	"github.com/skolldire/web-simplify/pkg/server/tcp"
	"github.com/skolldire/web-simplify/pkg/server/web_socket"
	"net/http"
	"sync"
)

const (
	appDefaultPort  = "8080"
	appDefaultScope = "local"
)

type Service interface {
	Run() error
	RegisterRoute(pattern string, handler http.HandlerFunc)
	RegisterWebSocket(pattern string, server web_socket.Service, f web_socket.ProcessingFunc)
	StartTCPServer(name string, server tcp.Service, f tcp.ProcessingFunc)
}

type App struct {
	Router      *chi.Mux
	Port        string
	Scope       string
	tcpServers  map[string]tcp.Service
	tcpHandlers map[string]tcp.ProcessingFunc
	wsServers   map[string]web_socket.Service
	wsHandlers  map[string]web_socket.ProcessingFunc
	mu          sync.Mutex
}

type Config struct {
	Port           string
	Scope          string
	LogLevel       string
	LogDestination string
}
