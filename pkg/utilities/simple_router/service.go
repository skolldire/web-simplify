package simple_router

import (
	"github.com/go-chi/chi/v5"
	"github.com/skolldire/web-simplify/pkg/utilities/app_profile"
	"github.com/skolldire/web-simplify/pkg/utilities/simple_router/docsify"
	"github.com/skolldire/web-simplify/pkg/utilities/simple_router/ping"
	"github.com/skolldire/web-simplify/pkg/utilities/simple_router/swagger"
	"net/http"
)

var _ Service = (*App)(nil)

func NewService(c Config) *App {
	return &App{
		Router: initRoutes(),
		Port:   setPort(c.Port),
	}
}

func (a App) Run() error {
	return http.ListenAndServe(a.Port, a.Router)

}

func initRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/ping", ping.NewService().Apply())
	if app_profile.IsStageProfile() || app_profile.IsTestProfile() {
		r.Mount("/swagger", swagger.NewService().Apply())
		r.Handle("/documentation-tech/*", docsify.NewService().Apply())
	}
	return r
}

func setPort(p string) string {
	if p == "" {
		return appDefaultPort
	}
	return p
}
