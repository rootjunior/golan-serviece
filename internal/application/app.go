package application

import (
	"context"
	"go-service/internal/application/mediator"
	"go-service/internal/config"
	"go-service/internal/core/commands"
	"go-service/internal/core/event_handlers"
	"go-service/internal/core/events"
	i "go-service/internal/core/interfaces"
	"go-service/internal/core/use_cases"
	"go-service/internal/infrastructure/bus"
	"go-service/internal/infrastructure/services"
	"go-service/internal/presentation/grpc"
	grpccontrollers "go-service/internal/presentation/grpc/controllers"
	"go-service/internal/presentation/rest"
	restV1 "go-service/internal/presentation/rest/v1"

	"go-service/internal/infrastructure/worker"
	"log"

	"go.uber.org/fx"
)

func RegisterMediator(m i.IMediator, createUserUC *use_cases.CreateUserUseCase, userCreatedEH *event_handlers.UserCreatedEventHandler) {
	// Queries
	// Command
	m.RegisterCommand(commands.CreateUserCommand{}, createUserUC)
	// Events
	m.RegisterEvent(events.UserCreateEvent{}, userCreatedEH)
}

func RunHooks(
	lifecycle fx.Lifecycle,
	worker i.IWorkerPool,
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
			bus.NewEventBus,
			restV1.NewController,
			grpccontrollers.NewUserController,
			fx.Annotate(mediator.NewMediator, fx.As(new(i.IMediator))),
			fx.Annotate(worker.NewWorkerPool, fx.As(new(i.IWorkerPool))),
			config.NewConfig,
			services.NewPostClient,
			use_cases.NewCreateUserUseCase,
			grpc.NewServer,
			rest.NewServer,
			event_handlers.NewUserCreatedEventHandler,
		),

		fx.Invoke(
			RegisterMediator,
			RunHooks,
		),
	)
}
