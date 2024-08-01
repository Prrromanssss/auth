package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
)

// Config holds the configuration for the application.
type Config struct {
	GRPC     Server   `validate:"required" yaml:"grpc"`
	Postgres Database `validate:"required" yaml:"postgres"`
}

// Server holds the configuration for the gRPC server.
type Server struct {
	Host string `validate:"required" yaml:"host"`
	Port string `validate:"required" yaml:"port"`
}

// Database holds the configuration for the PostgreSQL database.
type Database struct {
	Host     string `validate:"required" yaml:"host"`
	Port     string `validate:"required" yaml:"port"`
	User     string `validate:"required" yaml:"user"`
	Password string `validate:"required" yaml:"password"`
	DBName   string `validate:"required" yaml:"dbname"`
	SSLMode  string `validate:"required" yaml:"sslmode"`
}

// LoadConfig reads and parses the configuration from a file specified by the CONFIG_PATH environment variable.
func LoadConfig() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, errors.Wrap(err, "cannot read config")
	}

	return &cfg, nil
}
