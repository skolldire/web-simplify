package read_properties

import (
	"github.com/skolldire/web-simplify/pkg/client/rest"
	"github.com/skolldire/web-simplify/pkg/database/connect_sql/orm"
	"github.com/skolldire/web-simplify/pkg/database/connect_sql/simple"
	"github.com/skolldire/web-simplify/pkg/server/tcp"
	"github.com/skolldire/web-simplify/pkg/simplify/simple_router"
	"sync"
)

type Service interface {
	Apply() (Config, error)
}

type Config struct {
	Router       simple_router.Config       `json:"router"`
	RestClients  []map[string]rest.Config   `json:"rest_clients"`
	DBOrm        []map[string]orm.Config    `json:"db_orm"`
	DBSimple     []map[string]simple.Config `json:"db_simple"`
	TCPServer    []map[string]tcp.Config    `json:"tcp_server"`
	WebSockets   []string                   `json:"web_sockets"`
	Repositories map[string]interface{}     `json:"repositories"`
	UsesCases    map[string]interface{}     `json:"uses_cases"`
	Endpoints    map[string]interface{}     `json:"endpoints"`
}

type service struct {
	propertyFiles []string
	path          string
}

var (
	once     sync.Once
	instance *service
)
