package interfaces

import "context"

type IMediator interface {
	RegisterQuery(queryPrototype interface{}, handler QueryUseCase)
	ExecuteQuery(query interface{}) (interface{}, error)
	RegisterCommand(commandPrototype interface{}, handlers ...CommandUseCase)
	ExecuteCommand(ctx context.Context, command interface{}) ([]interface{}, error)
	RegisterEvent(eventPrototype interface{}, handlers ...EventHandler)
	HandleEvent(ctx context.Context, event interface{}) error
	PublishEvents(ctx context.Context, events ...interface{}) error
}
