package handlers

import (
	"context"
	"fmt"
	"go-service/internal/core/events"
	. "go-service/internal/core/interfaces"
)

type UserCreatedEventHandler struct {
	mediator Mediator
}

func NewUserCreatedEventHandler(mediator Mediator) *UserCreatedEventHandler {
	return &UserCreatedEventHandler{mediator: mediator}
}

func (eh *UserCreatedEventHandler) Execute(ctx context.Context, event Event) error {
	e := event.(events.UserCreateEvent)
	fmt.Println("Было обработано событе")
	_ = e
	return nil
}
