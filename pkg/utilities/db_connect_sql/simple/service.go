package simple

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/godror/godror"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	converter "github.com/skolldire/web-simplify/pkg/utilities/data_converter"
	"time"
)

type service struct {
	config Config
}

var _ Service = (*service)(nil)

func NewService(cfg Config) *service {
	return &service{
		config: cfg,
	}
}

func (s service) Init() *sql.DB {
	connLine := fmt.Sprintf(s.config.Dns, s.config.User, s.config.Password,
		s.config.Host, s.config.Port, s.config.Name)
	db, err := sql.Open(converter.DBToDriverMap(s.config.Motor), connLine)
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
