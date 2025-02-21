package app_engine

import (
	"database/sql"
	"github.com/skolldire/web-simplify/pkg/server/tcp"
	"github.com/skolldire/web-simplify/pkg/server/web_socket"
	"github.com/skolldire/web-simplify/pkg/utilities/log"
	"github.com/skolldire/web-simplify/pkg/utilities/requester"
	"github.com/skolldire/web-simplify/pkg/utilities/simple_router"
	"xorm.io/xorm"
)

type Engine struct {
	App                 *simple_router.App
	Tracer              log.Service
	Requesters          map[string]*requester.Requester
	DBOrmConnections    map[string]*xorm.Engine
	DBSimpleConnections map[string]*sql.DB
	TCPServer           map[string]tcp.Service
	WebSockets          map[string]*web_socket.Server
	RepositoriesConfig  map[string]interface{}
	UsesCasesConfig     map[string]interface{}
	HandlerConfig       map[string]interface{}
}
