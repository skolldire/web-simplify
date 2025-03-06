package orm

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/godror/godror"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	converter "github.com/skolldire/web-simplify/pkg/utilities/data_converter"
	"github.com/skolldire/web-simplify/pkg/utilities/log"
	"time"
	"xorm.io/xorm"
)

type service struct {
	config Config
	log    log.Service
}

var _ Service = (*service)(nil)

func NewService(c Config, log log.Service) *service {
	return &service{config: c, log: log}
}

func (s *service) Init() *xorm.Engine {
	connLine := fmt.Sprintf(s.config.Dns, s.config.User, s.config.Password,
		s.config.Host, s.config.Port, s.config.Name)
	engine, err := xorm.NewEngine(converter.DBToDriverMap(s.config.Motor), connLine)
	if err != nil {
		s.log.FatalError(context.Background(), err, map[string]interface{}{"message": "error in xorm creation"})
	}
	err = engine.DB().Ping()
	if err != nil {
		s.log.FatalError(context.Background(), err, map[string]interface{}{"message": "error pinging db by xorm"})
	}
	engine.SetMaxIdleConns(s.config.MaxIdleCons)
	engine.SetMaxOpenConns(s.config.MaxOpenCons)
	engine.SetConnMaxLifetime(time.Minute * time.Duration(s.config.ConnMaxLifetime))
	return engine
}
