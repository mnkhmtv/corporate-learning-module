package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string         `yaml:"env" env:"ENV" env-default:"development"`
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Auth     AuthConfig     `yaml:"auth"`
}

type ServerConfig struct {
	Port         string        `yaml:"port" env:"SERVER_PORT" env-default:"8080"`
	ReadTimeout  time.Duration `yaml:"read_timeout" env-default:"10s"`
	WriteTimeout time.Duration `yaml:"write_timeout" env-default:"10s"`
}

type DatabaseConfig struct {
	Host        string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	Port        int    `yaml:"port" env:"DB_PORT" env-default:"5432"`
	User        string `yaml:"user" env:"DB_USER" env-default:"postgres"`
	Password    string `yaml:"password" env:"DB_PASSWORD" env-required:"true"`
	DBName      string `yaml:"dbname" env:"DB_NAME" env-default:"training_db"`
	SSLMode     string `yaml:"sslmode" env:"DB_SSLMODE" env-default:"disable"`
	AutoMigrate bool   `yaml:"auto_migrate" env:"AUTO_MIGRATE" env-default:"true"`
}

type AuthConfig struct {
	JWTSecret string        `yaml:"jwt_secret" env:"JWT_SECRET" env-required:"true"`
	TokenTTL  time.Duration `yaml:"token_ttl" env:"TOKEN_TTL" env-default:"24h"`
}

// Load reads configuration from YAML file and environment variables
func Load(configPath string) (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
