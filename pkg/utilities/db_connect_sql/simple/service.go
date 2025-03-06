package simple

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/godror/godror"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	converter "github.com/skolldire/web-simplify/pkg/utilities/data_converter"
	"github.com/skolldire/web-simplify/pkg/utilities/log"
	"time"
)

type service struct {
	config Config
	log    log.Service
}

var _ Service = (*service)(nil)

func NewService(cfg Config, log log.Service) *service {
	return &service{
		config: cfg,
		log:    log,
	}
}

func (s service) Init() *sql.DB {
	connLine := fmt.Sprintf(s.config.Dns, s.config.User, s.config.Password,
		s.config.Host, s.config.Port, s.config.Name)
	db, err := sql.Open(converter.DBToDriverMap(s.config.Motor), connLine)
	if err != nil {
		s.log.FatalError(context.Background(), err, map[string]interface{}{"message": "error in sql.Open"})
	}
	err = db.Ping()
	if err != nil {
		s.log.FatalError(context.Background(), err, map[string]interface{}{"message": "error pinging db"})
	}
	db.SetMaxOpenConns(s.config.MaxOpenCons)
	db.SetMaxIdleConns(s.config.MaxIdleCons)
	db.SetConnMaxLifetime(time.Second * time.Duration(s.config.ConnMaxLifetime))
	return db
}
