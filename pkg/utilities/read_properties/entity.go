package read_properties

import (
	"github.com/skolldire/cash-manager-toolkit/pkg/kit/db_connection"
	"github.com/skolldire/web-simplify/pkg/server/tcp"
	"github.com/skolldire/web-simplify/pkg/utilities/simple_router"
)

type LoadProperties interface {
	Apply() (Config, error)
}

type Config struct {
	Router                    simple_router.Config              `json:"router"`
	DBOrmConnectionsConfig    []map[string]db_connection.Config `json:"db_orm_connections"`
	DBSimpleConnectionsConfig []map[string]db_connection.Config `json:"db_simple_connections"`
	TCPServerConfig           []map[string]tcp.Config           `json:"tcp_server_config"`
	RepositoriesConfig        map[string]interface{}            `json:"repositories_config"`
	UsesCasesConfig           map[string]interface{}            `json:"uses_cases_config"`
	HandlerConfig             map[string]interface{}            `json:"handlers_config"`
}
