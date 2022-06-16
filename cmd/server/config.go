package main

import (
	"context"

	"github.com/sethvargo/go-envconfig"
)

// Config stores app configuration
type Config struct {
	DBDriver          string `env:"DB_DRIVER,default=postgres"`
	DBSource          string `env:"DB_SOURCE,default=postgres://dev:pass@127.0.0.1:5432/devdb?sslmode=disable"`
	HTTPServerAddress string `env:"HTTP_SERVER_ADDRESS,default=0.0.0.0:8080"`
	ReadTimeout       int    `env:"READ_TIMEOUT,default=5"`
	IdleTimeout       int    `env:"IDLE_TIMEOUT,default=30"`
	ShutdownTimeout   int    `env:"SHUTDOWN_TIMEOUT,default=10"`
	CacheSize         int    `env:"CACHE_SIZE,default=1024"`
	NatsURL           string `env:"NATS_URL,default=nats://127.0.0.1:4222"`
	ClusterID         string `env:"CLUSTER_ID,default=test-cluster"`
	ClientID          string `env:"CLIENT_ID,default=client-456"`
	DurableName       string `env:"DURABLE_NAME,default=my-durable"`
}

// NewConfig reads config from env and creates config struct
func NewConfig() (*Config, error) {
	ctx := context.Background()
	var c Config
	if err := envconfig.Process(ctx, &c); err != nil {
		return nil, err
	}

	return &c, nil
}
