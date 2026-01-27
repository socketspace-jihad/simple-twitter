package postgresql

import (
	"database/sql"
	"fmt"
	"simple_twitter/internal/models"

	_ "github.com/lib/pq"
)

type PostgreSQL struct {
	Hostname string
	User     string
	Passwd   string
	Database string
	Port     string
	DB       *sql.DB
}

type PostgreSQLConfig func(*PostgreSQL)

func WithHostname(hostname string) PostgreSQLConfig {
	return func(ps *PostgreSQL) {
		ps.Hostname = hostname
	}
}

func NewPostgreSQL(configs ...PostgreSQLConfig) *PostgreSQL {
	p := &PostgreSQL{
		Hostname: "localhost",
		User:     "postgres",
		Passwd:   "test1234",
		Database: "simple-twitter",
		Port:     "5432",
	}
	for _, config := range configs {
		config(p)
	}
	return p
}

func (p *PostgreSQL) Connect() error {
	return nil
}

func (p *PostgreSQL) Disconnect() error {
	return nil
}

func init() {
	p := NewPostgreSQL()
	db, err := sql.Open(
		"postgres",
		fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			p.Hostname,
			p.Port,
			p.User,
			p.Passwd,
			p.Database,
		))
	if err != nil {
		panic(err)
	}
	p.DB = db
	models.SetUserDB(p)
	models.SetPostDB(p)
}
