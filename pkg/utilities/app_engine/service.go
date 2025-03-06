package app_engine

import (
	"database/sql"
	"github.com/skolldire/web-simplify/pkg/client/rest"
	"github.com/skolldire/web-simplify/pkg/server/tcp"
	"github.com/skolldire/web-simplify/pkg/server/web_socket"
	"github.com/skolldire/web-simplify/pkg/utilities/db_connect_sql/orm"
	"github.com/skolldire/web-simplify/pkg/utilities/db_connect_sql/simple"
	"github.com/skolldire/web-simplify/pkg/utilities/log"
	"github.com/skolldire/web-simplify/pkg/utilities/read_properties"
	"github.com/skolldire/web-simplify/pkg/utilities/simple_router"
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
		App:                 simple_router.NewService(c.Router),
		Tracer:              tracer,
		RestClients:         creteRequester(c.RestClients, tracer),
		DBOrmConnections:    createDBOrmConnections(c.DBOrm),
		DBSimpleConnections: createDBSimpleConnections(c.DBSimple),
		TCPServer:           createTCPServer(c.TCPServer, tracer),
		WebSockets:          createWebSockets(c.WebSockets, tracer),
		RepositoriesConfig:  c.Repositories,
		UsesCasesConfig:     c.UsesCases,
		HandlerConfig:       c.Endpoints,
	}
}

func creteTracer(c simple_router.Config) log.Service {
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

func createDBOrmConnections(cs []map[string]orm.Config) map[string]*xorm.Engine {
	connections := make(map[string]*xorm.Engine)
	for _, c := range cs {
		for k, v := range c {
			connections[k] = orm.NewService(v).Init()
		}
	}
	return connections
}

func createDBSimpleConnections(cs []map[string]simple.Config) map[string]*sql.DB {
	connections := make(map[string]*sql.DB)
	for _, c := range cs {
		for k, v := range c {
			connections[k] = simple.NewService(v).Init()
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
