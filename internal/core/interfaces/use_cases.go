package interfaces

import "context"

type CommandUseCase interface {
	Execute(ctx context.Context, command interface{}) (interface{}, error)
}

type QueryUseCase interface {
	Execute(query interface{}) (interface{}, error)
}

type EventHandler interface {
	Execute(ctx context.Context, event interface{}) error
}
