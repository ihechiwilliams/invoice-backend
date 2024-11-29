package appbase

import (
	"time"

	"invoice-backend/pkg/config"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS" env-default:"0.0.0.0:3000"`
	ServiceName   string `env:"SERVICE_NAME" env-default:"invoice-backend"`
	Port          string `env:"PORT" env-default:"3000"`
	ServerTimeout int64  `env:"SERVER_TIMEOUT" env-default:"120"`
	Env           string `env:"ENV"`
	LogLevel      string `env:"LOG_LEVEL" env-default:"debug"`
	SentryDSN     string `env:"SENTRY_DSN"`

	// Database
	DatabasePrimaryHost     string `env:"DATABASE_HOST" env-required:"true"`
	DatabaseReadReplicaHost string `env:"DATABASE_HOST_RO" env-required:"true"`
	DatabaseName            string `env:"DATABASE_NAME" env-required:"true"`
	DatabasePassword        string `env:"DATABASE_PASSWORD" env-required:"true"`
	DatabasePort            string `env:"DATABASE_PORT" env-default:"5432"`
	DatabaseUsername        string `env:"DATABASE_USERNAME" env-required:"true"`
}

func LoadConfig() (*Config, error) {
	c := new(Config)

	err := config.LoadConfig(c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Config) HTTPServerTimeout() time.Duration {
	return time.Duration(c.ServerTimeout) * time.Second
}
