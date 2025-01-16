package sqlite

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/skolldire/web-simplify/pkg/utilities/db_connect_sql"
	"time"
)

type service struct {
	config db_connection.Config
}

var _ Service = (*service)(nil)

func NewService(cfg db_connection.Config) *service {
	return &service{
		config: cfg,
	}
}

func (s service) Init() *sql.DB {
	db, err := sql.Open("sqlite3", s.config.Dns)
	if err != nil {
		panic(fmt.Errorf("error in sql.Open: %w", err))
	}
	err = db.Ping()
	if err != nil {
		panic(fmt.Errorf("error pinging db: %w", err))
	}
	db.SetMaxOpenConns(s.config.MaxOpenCons)
	db.SetMaxIdleConns(s.config.MaxIdleCons)
	db.SetConnMaxLifetime(time.Second * time.Duration(s.config.ConnMaxLifetime))
	return db
}
