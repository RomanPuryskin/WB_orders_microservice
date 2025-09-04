package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/orders_api/internal/database/postgres"
	"github.com/orders_api/internal/kafka"
	"github.com/orders_api/internal/logger"
)

type Config struct {
	Postgres   postgres.PostgresConfig
	ServerPort string `env:"SERVER_PORT" envDefault:":3000"`
	Logger     logger.Config
	Kafka      kafka.KafkaConfig
}

func MustLoad() (*Config, error) {
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("[MustLoad|Parse config from env file] , %w", err)
	}

	return cfg, nil
}
