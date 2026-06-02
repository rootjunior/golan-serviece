package config

import (
	"context"
	"log"

	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	WorkersCount         int    `env:"WORKERS_COUNT,default=1"`
	MediatorEventBufSize int    `env:"MEDIATOR_EVENT_BUF_SIZE,default=256"`
	ServerRESTAddress    string `env:"SERVER_REST_ADDRESS,default=:8080"`
	ServerGRPCAddress    string `env:"SERVER_GRPC_ADDRESS,default=:50051"`
	EventBusSize         int    `env:"EVENT_BUS_SIZE,default=256"`
}

func NewConfig() *Config {
	cfg := &Config{}
	if err := envconfig.Process(context.Background(), cfg); err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}
	return cfg
}
