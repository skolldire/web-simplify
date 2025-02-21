package simple_router

import (
	"github.com/go-chi/chi/v5"
)

const (
	appDefaultPort  = "8080"
	appDefaultScope = "local"
)

type Service interface {
	Run() error
}

type App struct {
	Router *chi.Mux
	Port   string
	Scope  string
}

type Config struct {
	Port           string
	Scope          string
	LogLevel       string
	LogDestination string
}
