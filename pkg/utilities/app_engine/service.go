package app_engine

import (
	"database/sql"
	"github.com/skolldire/web-simplify/pkg/server/tcp"
	"github.com/skolldire/web-simplify/pkg/utilities/db_connect_sql/orm"
	"github.com/skolldire/web-simplify/pkg/utilities/db_connect_sql/simple"
	"github.com/skolldire/web-simplify/pkg/utilities/log"
	"github.com/skolldire/web-simplify/pkg/utilities/read_properties/viper"
	"github.com/skolldire/web-simplify/pkg/utilities/simple_router"
	"xorm.io/xorm"
)

func NewApp() *Engine {
	v := viper.NewService()
	c, err := v.Apply()
	if err != nil {
		panic(err)
	}
	return &Engine{
		App:                 simple_router.NewService(c.Router),
		Tracer:              creteTracer(c.Router),
		HttpClient:          nil,
		DBOrmConnections:    createDBOrmConnections(c.DBOrm),
		DBSimpleConnections: createDBSimpleConnections(c.DBSimple),
		TCPServer:           createTCPServer(c.TCPServer, creteTracer(c.Router)),
		RepositoriesConfig:  c.Repositories,
		UsesCasesConfig:     c.UsesCases,
		HandlerConfig:       c.Endpoints,
	}
}

func creteTracer(config simple_router.Config) log.Service {
	return log.NewService(config.LogLevel)
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
	for _, c := range cs {
		for k, v := range c {
			d := tcp.Dependencies{
				Config: v,
				Log:    log,
			}
			servers[k] = tcp.NewService(d)
		}
	}
	return servers
}
