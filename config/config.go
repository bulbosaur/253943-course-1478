package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Env EnvConfig `mapstructure:"env"`
	GRPC GRPCConfig `mapstructure:"grpc"`
}

type EnvConfig struct {
	LogLevel string `mapstructure:"loglevel"`
}

type GRPCConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

func LoadConfig(envPath, yamlPath string) (*Config, error) {
	viper.SetConfigFile(envPath)
	err := viper.MergeInConfig()
	if err != nil {
		log.Printf("config.LoadConfig: can't find .env: %v", err)
	} else {
		log.Println("Loaded .env file")
	}

	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return nil, fmt.Errorf("config.LoadConfig: error reading config file: %w", err)
	}

	expandedData := os.ExpandEnv(string(data))

	viper.SetConfigType("yaml")
	err = viper.ReadConfig(strings.NewReader(expandedData))
	if err != nil {
		return nil, fmt.Errorf("config.LoadConfig: error parsing config: %w", err)
	}

	viper.AutomaticEnv()

	viper.BindEnv("grpc.host", "GRPC_HOST")
	viper.BindEnv("grpc.port", "GRPC_PORT")

	var cfg Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("config.LoadConfig: unable to decode config into struct: %w", err)
	}

	return &cfg, nil
}
