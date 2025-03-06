package app_engine

import (
	"database/sql"
	"github.com/skolldire/web-simplify/pkg/client/rest"
	"github.com/skolldire/web-simplify/pkg/config/read_properties"
	orm2 "github.com/skolldire/web-simplify/pkg/database/connect_sql/orm"
	simple2 "github.com/skolldire/web-simplify/pkg/database/connect_sql/simple"
	"github.com/skolldire/web-simplify/pkg/server/tcp"
	"github.com/skolldire/web-simplify/pkg/server/web_socket"
	simple_router2 "github.com/skolldire/web-simplify/pkg/simplify/simple_router"
	"github.com/skolldire/web-simplify/pkg/utilities/log"
	"xorm.io/xorm"
)

func NewApp() *Engine {
	v := read_properties.NewService()
	c, err := v.Apply()
	if err != nil {
		panic(err)
	}
	tracer := creteTracer(c.Router)
	return &Engine{
		App:                 simple_router2.NewService(c.Router, tracer),
		Tracer:              tracer,
		RestClients:         creteRequester(c.RestClients, tracer),
		DBOrmConnections:    createDBOrmConnections(c.DBOrm, tracer),
		DBSimpleConnections: createDBSimpleConnections(c.DBSimple, tracer),
		TCPServer:           createTCPServer(c.TCPServer, tracer),
		WebSockets:          createWebSockets(c.WebSockets, tracer),
		RepositoriesConfig:  c.Repositories,
		UsesCasesConfig:     c.UsesCases,
		HandlerConfig:       c.Endpoints,
	}
}

func creteTracer(c simple_router2.Config) log.Service {
	return log.NewService(log.Config{
		Level:          c.LogLevel,
		LogDestination: c.LogDestination,
	})
}

func creteRequester(cs []map[string]rest.Config, l log.Service) map[string]rest.Service {
	requesters := map[string]rest.Service{}
	for _, c := range cs {
		for k, v := range c {
			requesters[k] = rest.NewService(v, l)
		}
	}
	return requesters
}

func createDBOrmConnections(cs []map[string]orm2.Config, l log.Service) map[string]*xorm.Engine {
	connections := make(map[string]*xorm.Engine)
	for _, c := range cs {
		for k, v := range c {
			connections[k] = orm2.NewService(v, l).Init()
		}
	}
	return connections
}

func createDBSimpleConnections(cs []map[string]simple2.Config, l log.Service) map[string]*sql.DB {
	connections := make(map[string]*sql.DB)
	for _, c := range cs {
		for k, v := range c {
			connections[k] = simple2.NewService(v, l).Init()
		}
	}
	return connections
}

func createTCPServer(cs []map[string]tcp.Config, log log.Service) map[string]tcp.Service {
	servers := make(map[string]tcp.Service)

	for _, configMap := range cs {
		for key, config := range configMap {
			server := tcp.NewService(tcp.Dependencies{
				Config: config,
				Log:    log,
			})
			if server != nil {
				servers[key] = server
			}
		}
	}

	return servers
}

func createWebSockets(list []string, log log.Service) map[string]web_socket.Service {
	servers := make(map[string]web_socket.Service)
	for _, key := range list {
		server := web_socket.NewServer(log)
		if server != nil {
			servers[key] = server
		}
	}
	return servers
}
