package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"lyceum/config"
	"lyceum/pkg/db"
	"net"
	"path/filepath"

	lg "lyceum/logger"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
)

func main() {
	var command string

	flag.StringVar(&command, "command", "up", "migration command: up, down, force, version")
	flag.Parse()

	var (
		configDir = "./config"
		envPath   = filepath.Join(configDir, ".env")
		yamlPath  = filepath.Join(configDir, "config.yaml")
	)

	cfg := config.LoadConfig(envPath, yamlPath)

	logger := lg.NewLogger(cfg.Env.LogLevel)
	defer logger.Sync() //nolint:errcheck // error checking is redundant here

	ctx := lg.WithRequestID(context.Background(), "")
	ctx = lg.WithLogger(ctx, logger)

	m, err := makeMigrations(cfg.PostgreSQL)
	if err != nil {
		logger.Error(ctx, "failed to create migrate instance", zap.Any("error", err))
	}
	defer m.Close()

	switch command {
	case "up":
		if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			logger.Error(ctx, "failed to apply migrations", zap.Any("error", err))
		}
		logger.Info(ctx, "Migrations applied successfully!")

	case "down":
		if err = m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			logger.Error(ctx, "failed to rollback migrations", zap.Any("error", err))
		}
		logger.Info(ctx, "Migrations rolled back successfully!")

	case "version":
		version, dirty, err := m.Version() //nolint:govet // can't declarate err with "="
		if err != nil {
			logger.Error(ctx, "failed to get version", zap.Any("error", err))
		}
		logger.Info(ctx, "About version", zap.Any("Current version", version), zap.Any("Dirty", dirty))

	default:
		logger.Error(ctx, "unknown command", zap.Any("command", command))
		logger.Info(ctx, "Available commands: up, down, version")
	}
}

func makeMigrations(cfg db.Config) (*migrate.Migrate, error) {
	var migrationsPath string
	flag.StringVar(&migrationsPath, "path", "./migrations", "path to migrations directory")

	hostPort := net.JoinHostPort(cfg.Host, cfg.Port)

	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		cfg.Username,
		cfg.Password,
		hostPort,
		cfg.DBName,
	)

	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		dsn,
	)

	return m, err
}
