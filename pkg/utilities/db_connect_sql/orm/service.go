package orm

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/godror/godror"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	converter "github.com/skolldire/web-simplify/pkg/utilities/data_converter"
	"time"
	"xorm.io/xorm"
)

type service struct {
	config Config
}

var _ Service = (*service)(nil)

func NewService(c Config) *service {
	return &service{config: c}
}

func (s *service) Init() *xorm.Engine {
	connLine := fmt.Sprintf(s.config.Dns, s.config.User, s.config.Password,
		s.config.Host, s.config.Port, s.config.Name)
	engine, err := xorm.NewEngine(converter.DBToDriverMap(s.config.Motor), connLine)
	if err != nil {
		panic(fmt.Errorf("error in xorm creation: %w", err))
	}
	err = engine.DB().Ping()
	if err != nil {
		panic(fmt.Errorf("error pinging db by xorm: %w", err))
	}
	engine.SetMaxIdleConns(s.config.MaxIdleCons)
	engine.SetMaxOpenConns(s.config.MaxOpenCons)
	engine.SetConnMaxLifetime(time.Minute * time.Duration(s.config.ConnMaxLifetime))
	return engine
}
