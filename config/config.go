// config/config.go
package config

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Env  EnvConfig  `mapstructure:"env"`
	GRPC GRPCConfig `mapstructure:"grpc"`
	HTTP HTTPConfig `mapstructure:"http"`
}

type EnvConfig struct {
	LogLevel string `mapstructure:"loglevel"`
}

type GRPCConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type HTTPConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Timeout time.Duration `mapstructure:"timeout"`
}

func loadEnvFile(envPath string) error {
	data, err := os.ReadFile(envPath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
			(strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
			value = value[1 : len(value)-1]
		}

		os.Setenv(key, value)
	}
	return nil
}

func LoadConfig(envPath, yamlPath string) (*Config, error) {
	err := loadEnvFile(envPath)
	if err != nil {
		log.Printf("config.LoadConfig: can't read .env file %s: %v", envPath, err)
	}

	data, err := os.ReadFile(yamlPath)
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
	viper.BindEnv("http.host", "HTTP_HOST")
	viper.BindEnv("grpc.port", "GRPC_PORT")
	viper.BindEnv("http.port", "HTTP_PORT")

	var cfg Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("config.LoadConfig: unable to decode config into struct: %w", err)
	}

	return &cfg, nil
}
