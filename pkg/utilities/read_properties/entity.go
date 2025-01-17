package read_properties

import (
	"github.com/skolldire/web-simplify/pkg/server/tcp"
	"github.com/skolldire/web-simplify/pkg/utilities/db_connect_sql/orm"
	"github.com/skolldire/web-simplify/pkg/utilities/db_connect_sql/simple"
	"github.com/skolldire/web-simplify/pkg/utilities/simple_router"
)

type LoadProperties interface {
	Apply() (Config, error)
}

type Config struct {
	Router       simple_router.Config       `json:"router"`
	DBOrm        []map[string]orm.Config    `json:"db_orm"`
	DBSimple     []map[string]simple.Config `json:"db_simple"`
	TCPServer    []map[string]tcp.Config    `json:"tcp_server"`
	Repositories map[string]interface{}     `json:"repositories"`
	UsesCases    map[string]interface{}     `json:"uses_cases"`
	Endpoints    map[string]interface{}     `json:"endpoints"`
}
