package db_connection

const (
	Oracle   = "oracle"
	Mysql    = "mysql"
	Postgres = "postgres"
	SQLite   = "sqlite"
)

type Config struct {
	Motor           string `json:"driver"`
	Name            string `json:"name"`
	User            string `json:"user"`
	Password        string `json:"password"`
	Host            string `json:"host"`
	Port            string `json:"port"`
	Dns             string `json:"dns"`
	MaxOpenCons     int    `json:"max-open-cons"`
	MaxIdleCons     int    `json:"max-idle-cons"`
	ConnMaxLifetime uint   `json:"conn-max-lifetime"`
}
