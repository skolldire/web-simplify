package sqlite

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/skolldire/web-simplify/pkg/utilities/db_connect_sql"
	"log"
	"time"
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
	engine, err := xorm.NewEngine("sqlite3", s.config.Dns)
	if err != nil {
		log.Fatal(err)
	}
	engine.SetMaxIdleConns(s.config.MaxIdleCons)
	engine.SetMaxOpenConns(s.config.MaxOpenCons)
	engine.SetConnMaxLifetime(time.Minute * time.Duration(s.config.ConnMaxLifetime))
	return engine
}
