package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Host        string        `yaml:"host"`
	Port        int           `yaml:"port"`
	Password    string        `yaml:"password"`
	DB          int           `yaml:"db"`
	MaxRetries  int           `yaml:"max_retries"`
	DialTimeout time.Duration `yaml:"dial_timeout"`
	Timeout     time.Duration `yaml:"timeout"`
}

func NewRedisClient(ctx context.Context, cfg Config) (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	db := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
	})

	if err := db.Ping(ctx).Err(); err != nil {
		log.Printf("failed to connect to redis server: %s\n", err.Error())
		return nil, err
	}

	return db, nil
}
