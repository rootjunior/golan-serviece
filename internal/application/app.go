package application

import (
	"context"
	"go-service/internal/application/bootstrap"
	"go-service/internal/application/mediator"
	"go-service/internal/application/worker"

	"go-service/internal/config"
	"go-service/internal/core/handlers"
	. "go-service/internal/core/interfaces"
	"go-service/internal/infrastructure/bus"
	"go-service/internal/infrastructure/services"
	"go-service/internal/presentation/grpc"
	grpccontrollers "go-service/internal/presentation/grpc/controllers"
	"go-service/internal/presentation/rest"
	restV1 "go-service/internal/presentation/rest/v1"

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

func NewApplication() *fx.App {
	return fx.New(
		fx.Provide(
			restV1.NewController,
			grpccontrollers.NewUserController,
			fx.Annotate(bus.NewEventBus, fx.As(new(EventBus))),
			fx.Annotate(mediator.NewMediator, fx.As(new(Mediator))),
			fx.Annotate(worker.NewWorkerPool, fx.As(new(WorkerPool))),
			config.NewConfig,
			services.NewPostClient,
			handlers.NewCreateUserHandler,
			grpc.NewServer,
			rest.NewServer,
			handlers.NewUserCreatedEventHandler,
		),

		fx.Invoke(
			bootstrap.RegisterMediator,
			RunHooks,
		),
	)
}
