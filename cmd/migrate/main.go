package main

import (
	"errors"
	"flag"
	"fmt"
	"lyceum/config"
	"net"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var migrationsPath string
	var command string

	flag.StringVar(&migrationsPath, "path", "./migrations", "path to migrations directory")
	flag.StringVar(&command, "command", "up", "migration command: up, down, force, version")
	flag.Parse()

	var (
		configDir = "./config"
		envPath   = filepath.Join(configDir, ".env")
		yamlPath  = filepath.Join(configDir, "config.yaml")
	)

	cfg := config.LoadConfig(envPath, yamlPath)

	hostPort := net.JoinHostPort(cfg.PostgreSQL.Host, cfg.PostgreSQL.Port)

	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		cfg.PostgreSQL.Username,
		cfg.PostgreSQL.Password,
		hostPort,
		cfg.PostgreSQL.DBName,
	)

	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		dsn,
	)
	if err != nil {
		fmt.Printf("failed to create migrate instance: %v\n", err)
		os.Exit(1)
	}
	defer m.Close()

	switch command {
	case "up":
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			fmt.Printf("failed to apply migrations: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Migrations applied successfully!")

	case "down":
		if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			fmt.Printf("failed to rollback migrations: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Migrations rolled back successfully!")

	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			fmt.Printf("failed to get version: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Current version: %d, Dirty: %v\n", version, dirty)

	default:
		fmt.Printf("unknown command: %s\n", command)
		fmt.Println("Available commands: up, down, version")
		os.Exit(1)
	}
}
