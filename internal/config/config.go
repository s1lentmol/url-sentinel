package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config holds application configuration
type Config struct {
	Env        string     `yaml:"env" env:"ENV" env-default:"local"`
	Database   Database   `yaml:"database"`
	HTTPServer HTTPServer `yaml:"http_server"`
}

// Database holds database configuration
type Database struct {
	Host     string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	Port     string `yaml:"port" env:"DB_PORT" env-default:"5432"`
	User     string `yaml:"user" env:"DB_USER" env-default:"postgres"`
	Password string `yaml:"password" env:"DB_PASSWORD" env-default:"postgres"`
	DBName   string `yaml:"dbname" env:"DB_NAME" env-default:"url_sentinel"`
	SSLMode  string `yaml:"sslmode" env:"DB_SSLMODE" env-default:"disable"`
}

// DSN returns the database connection string
func (d Database) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.DBName, d.SSLMode,
	)
}

// HTTPServer holds HTTP server configuration
type HTTPServer struct {
	Address         string        `yaml:"address" env:"HTTP_ADDRESS" env-default:"0.0.0.0:8080"`
	ReadTimeout     time.Duration `yaml:"read_timeout" env:"HTTP_READ_TIMEOUT" env-default:"5s"`
	WriteTimeout    time.Duration `yaml:"write_timeout" env:"HTTP_WRITE_TIMEOUT" env-default:"10s"`
	IdleTimeout     time.Duration `yaml:"idle_timeout" env:"HTTP_IDLE_TIMEOUT" env-default:"60s"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env:"HTTP_SHUTDOWN_TIMEOUT" env-default:"10s"`
}

// MustLoad loads configuration from file and environment variables
func MustLoad() *Config {
	var cfg Config

	// Try to load from config file if CONFIG_PATH is set
	configPath := os.Getenv("CONFIG_PATH")
	if configPath != "" {
		if _, err := os.Stat(configPath); err == nil {
			if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
				log.Fatalf("failed to read config file: %v", err)
			}
		} else {
			log.Printf("config file not found at %s, using environment variables", configPath)
		}
	}

	// Read from environment variables (will override file config)
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("failed to read environment variables: %v", err)
	}

	return &cfg
}
