package app_engine

import (
	"database/sql"
	"github.com/skolldire/web-simplify/pkg/client/rest"
	"github.com/skolldire/web-simplify/pkg/server/tcp"
	"github.com/skolldire/web-simplify/pkg/server/web_socket"
	"github.com/skolldire/web-simplify/pkg/simplify/simple_router"
	"github.com/skolldire/web-simplify/pkg/utilities/log"
	"xorm.io/xorm"
)

type Engine struct {
	App                 simple_router.Service
	Tracer              log.Service
	RestClients         map[string]rest.Service
	DBOrmConnections    map[string]*xorm.Engine
	DBSimpleConnections map[string]*sql.DB
	TCPServer           map[string]tcp.Service
	WebSockets          map[string]web_socket.Service
	RepositoriesConfig  map[string]interface{}
	UsesCasesConfig     map[string]interface{}
	HandlerConfig       map[string]interface{}
}
