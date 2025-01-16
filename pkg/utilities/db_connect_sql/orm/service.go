package orm

import (
	"github.com/skolldire/web-simplify/pkg/utilities/db_connect_sql"
	"github.com/skolldire/web-simplify/pkg/utilities/db_connect_sql/orm/mysql"
	"github.com/skolldire/web-simplify/pkg/utilities/db_connect_sql/orm/oracle"
	"github.com/skolldire/web-simplify/pkg/utilities/db_connect_sql/orm/postgres"
	"github.com/skolldire/web-simplify/pkg/utilities/db_connect_sql/orm/sqlite"
	"xorm.io/xorm"
)

type service struct {
	config db_connection.Config
}

var _ Service = (*service)(nil)

func NewService(c db_connection.Config) *service {
	return &service{config: c}
}

func (s *service) Init() *xorm.Engine {
	switch s.config.Motor {
	case db_connection.Mysql:
		return mysql.NewService(s.config).Init()
	case db_connection.Oracle:
		return oracle.NewService(s.config).Init()
	case db_connection.Postgres:
		return postgres.NewService(s.config).Init()
	case db_connection.SQLite:
		return sqlite.NewService(s.config).Init()
	default:
		return nil
	}
}
