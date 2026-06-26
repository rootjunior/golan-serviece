package interfaces

import "context"

type CommandHandler interface {
	Execute(ctx context.Context, command Command) (Result, error)
}

type QueryHandler interface {
	Execute(query Query) (Result, error)
}

type EventHandler interface {
	Execute(ctx context.Context, event Event) error
}
