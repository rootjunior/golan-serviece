package app

import (
	"go-service/internal/application/bootstrap"
	"go-service/internal/application/mediator"
	"go-service/internal/application/worker"

	"go-service/internal/config"
	"go-service/internal/core/handlers"
	. "go-service/internal/core/interfaces"
	"go-service/internal/infrastructure/adapters"
	"go-service/internal/infrastructure/bus"
	"go-service/internal/presentation/grpc"
	grpccontrollers "go-service/internal/presentation/grpc/controllers"
	"go-service/internal/presentation/rest"
	restV1 "go-service/internal/presentation/rest/v1"

	"go.uber.org/fx"
)

func NewApplication() App {
	return fx.New(
		fx.Provide(
			restV1.NewController,
			grpccontrollers.NewUserController,
			fx.Annotate(bus.NewEventBus, fx.As(new(EventBus))),
			fx.Annotate(mediator.NewMediator, fx.As(new(Mediator))),
			fx.Annotate(worker.NewWorkerPool, fx.As(new(WorkerPool))),
			config.NewConfig,
			adapters.NewPostClient,
			handlers.NewCreateUserHandler,
			grpc.NewServer,
			rest.NewServer,
			handlers.NewUserCreatedEventHandler,
		),

		fx.Invoke(
			bootstrap.RegisterMediator,
			bootstrap.RunHooks,
		),
	)
}
