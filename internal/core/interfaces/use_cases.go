package interfaces

import "context"

type CommandUseCase interface {
	Execute(ctx context.Context, command Command) (Result, error)
}

type QueryUseCase interface {
	Execute(query Query) (Result, error)
}

type EventHandler interface {
	Execute(ctx context.Context, event Event) error
}
