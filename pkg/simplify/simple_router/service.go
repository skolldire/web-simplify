package simple_router

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/skolldire/web-simplify/pkg/server/tcp"
	"github.com/skolldire/web-simplify/pkg/server/web_socket"
	"github.com/skolldire/web-simplify/pkg/simplify/simple_router/docsify"
	"github.com/skolldire/web-simplify/pkg/simplify/simple_router/ping"
	"github.com/skolldire/web-simplify/pkg/simplify/simple_router/swagger"
	"github.com/skolldire/web-simplify/pkg/utilities/app_profile"
	"github.com/skolldire/web-simplify/pkg/utilities/log"
	"net/http"
	"net/http/pprof"
	"os"
)

var _ Service = (*App)(nil)

func NewService(c Config, l log.Service) Service {
	routes := initRoutes()
	if routes == nil {
		panic("Router initialization failed")
	}

	return &App{
		Router:      routes,
		Port:        setPort(c.Port),
		tcpServers:  make(map[string]tcp.Service),
		tcpHandlers: make(map[string]tcp.ProcessingFunc),
		wsServers:   make(map[string]web_socket.Service),
		wsHandlers:  make(map[string]web_socket.ProcessingFunc),
		log:         l,
	}
}

func (a *App) Run() error {
	return http.ListenAndServe(a.Port, a.Router)
}

func (a *App) RegisterRoute(pattern string, handler http.HandlerFunc) {
	a.Router.HandleFunc(pattern, handler)
}

func (a *App) StartTCPServer(name string, server tcp.Service, f tcp.ProcessingFunc) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if _, exists := a.tcpServers[name]; exists {
		fmt.Printf("TCP Server '%s' ya está registrado.\n", name)
		return
	}

	a.tcpServers[name] = server
	a.tcpHandlers[name] = f

	go server.Start(context.Background(), f)

	fmt.Printf("TCP Server '%s' iniciado correctamente en segundo plano.\n", name)
}

func (a *App) RegisterWebSocket(pattern string, server web_socket.Service, f web_socket.ProcessingFunc) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if _, exists := a.wsServers[pattern]; exists {
		fmt.Printf("WebSocket en '%s' ya está registrado.\n", pattern)
		return
	}

	server.SetProcessingFunc(f)

	a.Router.HandleFunc(pattern, server.HandleNewConnection)

	a.wsServers[pattern] = server
	a.wsHandlers[pattern] = f

	fmt.Printf("WebSocket registrado en '%s'.\n", pattern)
}

func registerPprofRoutes(router chi.Router) {
	pprofMux := http.NewServeMux()
	pprofMux.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	pprofMux.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	pprofMux.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	pprofMux.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	pprofMux.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))

	router.Mount("/debug/pprof", http.StripPrefix("/debug/pprof", pprofMux))
}

func initRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/ping", ping.NewService().Apply())

	if app_profile.IsStageProfile() || app_profile.IsTestProfile() {
		r.Mount("/swagger", swagger.NewService().Apply())
		r.Handle("/documentation-tech/*", docsify.NewService().Apply())
	}

	registerPprofRoutes(r)
	return r
}

func setPort(p string) string {
	if p != "" {
		return ":" + p
	}
	if envPort := os.Getenv("PORT"); envPort != "" {
		return ":" + envPort
	}
	return ":" + appDefaultPort
}
