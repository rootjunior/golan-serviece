package interfaces

import "context"

type IMediator interface {
	RegisterQuery(queryPrototype Query, handler QueryUseCase)
	ExecuteQuery(query Query) (Result, error)
	RegisterCommand(commandPrototype Command, handlers ...CommandUseCase)
	ExecuteCommand(ctx context.Context, command Command) ([]Result, error)
	RegisterEvent(eventPrototype Event, handlers ...EventHandler)
	HandleEvent(ctx context.Context, event Event) error
	PublishEvents(ctx context.Context, events ...Event) error
}
