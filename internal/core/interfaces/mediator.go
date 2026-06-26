package interfaces

import "context"

type Mediator interface {
	RegisterQuery(queryPrototype Query, handler QueryHandler)
	HandleQuery(query Query) (Result, error)
	RegisterCommand(commandPrototype Command, handlers ...CommandHandler)
	HandleCommand(ctx context.Context, command Command) ([]Result, error)
	RegisterEvent(eventPrototype Event, handlers ...EventHandler)
	HandleEvent(ctx context.Context, event Event) error
	PublishEvents(ctx context.Context, events ...Event) error
}
