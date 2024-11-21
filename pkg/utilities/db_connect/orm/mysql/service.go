package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/skolldire/web-simplify/pkg/utilities/db_connect"
	"github.com/skolldire/web-simplify/pkg/utilities/db_connect/orm"
	"log"
	"time"
	"xorm.io/xorm"
)

type service struct {
	config db_connection.Config
}

var _ orm.Service = (*service)(nil)

func NewService(c db_connection.Config) *service {
	return &service{config: c}
}

func (s *service) Init() *xorm.Engine {
	engine, err := xorm.NewEngine("mysql", fmt.Sprintf(s.config.Dns, s.config.User,
		s.config.Password, s.config.Host, s.config.Port, s.config.Name))
	if err != nil {
		log.Fatal(err)
	}
	engine.SetMaxIdleConns(s.config.MaxIdleCons)
	engine.SetMaxOpenConns(s.config.MaxOpenCons)
	engine.SetConnMaxLifetime(time.Minute * time.Duration(s.config.ConnMaxLifetime))
	return engine
}
