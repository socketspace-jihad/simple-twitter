package postgresql

import (
	"database/sql"
	"fmt"
	"os"
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

func WithUser(user string) PostgreSQLConfig {
	return func(ps *PostgreSQL) {
		ps.User = user
	}
}

func WithPasswd(passwd string) PostgreSQLConfig {
	return func(ps *PostgreSQL) {
		ps.Passwd = passwd
	}
}

func WithDatabase(db string) PostgreSQLConfig {
	return func(ps *PostgreSQL) {
		ps.Database = db
	}
}

func WithPort(port string) PostgreSQLConfig {
	return func(ps *PostgreSQL) {
		ps.Port = port
	}
}

func WithHostnameEnv(env string) PostgreSQLConfig {
	return func(ps *PostgreSQL) {
		if os.Getenv(env) != "" {
			ps.Hostname = os.Getenv(env)
		}
	}
}

func WithUserEnv(env string) PostgreSQLConfig {
	return func(ps *PostgreSQL) {
		if os.Getenv(env) != "" {
			ps.User = os.Getenv(env)
		}
	}
}

func WithPasswdEnv(env string) PostgreSQLConfig {
	return func(ps *PostgreSQL) {
		if os.Getenv(env) != "" {
			ps.Passwd = os.Getenv(env)
		}
	}
}

func WithDatabaseEnv(env string) PostgreSQLConfig {
	return func(ps *PostgreSQL) {
		if os.Getenv(env) != "" {
			ps.Database = os.Getenv(env)
		}
	}
}

func WithPortEnv(env string) PostgreSQLConfig {
	return func(ps *PostgreSQL) {
		if os.Getenv(env) != "" {
			ps.Port = os.Getenv(env)
		}
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
	p := NewPostgreSQL(
		WithHostnameEnv("DB_HOSTNAME"),
		WithPortEnv("DB_PORT"),
		WithUserEnv("DB_USER"),
		WithPasswdEnv("DB_PASSWORD"),
		WithDatabaseEnv("DB_DATABASE"),
	)
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
