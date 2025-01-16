package postgres

import (
	"fmt"
	_ "github.com/lib/pq"
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
	engine, err := xorm.NewEngine("postgres", fmt.Sprintf(s.config.Dns, s.config.User,
		s.config.Password, s.config.Host, s.config.Port, s.config.Name))
	if err != nil {
		log.Fatal(err)
	}
	engine.SetMaxIdleConns(s.config.MaxIdleCons)
	engine.SetMaxOpenConns(s.config.MaxOpenCons)
	engine.SetConnMaxLifetime(time.Minute * time.Duration(s.config.ConnMaxLifetime))
	return engine
}
