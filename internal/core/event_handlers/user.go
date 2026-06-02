package event_handlers

import (
	"context"
	"fmt"
	"go-service/internal/core/events"
	i "go-service/internal/core/interfaces"
)

type UserCreatedEventHandler struct {
	mediator i.IMediator
}

func NewUserCreatedEventHandler(mediator i.IMediator) *UserCreatedEventHandler {
	return &UserCreatedEventHandler{mediator: mediator}
}

func (eh *UserCreatedEventHandler) Execute(ctx context.Context, event interface{}) error {
	e := event.(events.UserCreateEvent)
	fmt.Println("Было обработано событе")
	_ = e
	return nil
}
