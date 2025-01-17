package app_engine

import (
	"database/sql"
	"github.com/skolldire/web-simplify/pkg/server/tcp"
	"github.com/skolldire/web-simplify/pkg/utilities/log"
	"github.com/skolldire/web-simplify/pkg/utilities/simple_router"
	"net/http"
	"xorm.io/xorm"
)

type Engine struct {
	App                 *simple_router.App
	Tracer              log.Service
	HttpClient          map[string]http.Client
	DBOrmConnections    map[string]*xorm.Engine
	DBSimpleConnections map[string]*sql.DB
	TCPServer           map[string]tcp.Service
	RepositoriesConfig  map[string]interface{}
	UsesCasesConfig     map[string]interface{}
	HandlerConfig       map[string]interface{}
}
