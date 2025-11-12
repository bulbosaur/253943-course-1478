package db

import (
	"context"
	"fmt"
	"net"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Username string `env:"POSTGRES_USER"     env-default:"postgres"`
	Password string `env:"POSTGRES_PASSWORD" env-default:"postgres"`
	Host     string `env:"POSTGRES_HOST"     env-default:"db"`
	Port     string `env:"POSTGRES_PORT"     env-default:"5432"`
	DBName   string `env:"POSTGRES_DB"       env-default:"postgres"`
}

type Database struct {
	Pool *pgxpool.Pool
}

func NewPostgres(config Config) (*Database, error) {
	hostPort := net.JoinHostPort(config.Host, config.Port)
	dataSource := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable&search_path=public",
		config.Username, config.Password, hostPort, config.DBName)

	pool, err := pgxpool.New(context.Background(), dataSource)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(context.Background())
	if err != nil {
		pool.Close()
		return nil, err
	}

	return &Database{
		Pool: pool,
	}, nil
}
