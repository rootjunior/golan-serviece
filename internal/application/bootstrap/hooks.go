package bootstrap

import (
	"context"
	"go-service/internal/config"
	. "go-service/internal/core/interfaces"
	"go-service/internal/presentation/grpc"
	"go-service/internal/presentation/rest"
	"log"

	"go.uber.org/fx"
)

func RunHooks(
	lifecycle fx.Lifecycle,
	worker WorkerPool,
	serverREST *rest.Server,
	serverGRPC *grpc.Server,
	cfg *config.Config,
) {
	ctx, cancel := context.WithCancel(context.Background())

	lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			log.Println("Starting workers and server")
			worker.StartProcessEvents(ctx)
			serverREST.Start(ctx)
			log.Printf("REST Server started: %s", cfg.ServerRESTAddress)
			serverGRPC.Start(ctx)
			log.Printf("GRPC Server started: %s", cfg.ServerGRPCAddress)
			return nil
		},
		OnStop: func(_ context.Context) error {
			cancel()
			log.Println("application stopped")
			return nil
		},
	})
}
