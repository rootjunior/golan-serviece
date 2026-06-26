package handlers

import (
	"context"
	"fmt"
	"go-service/internal/core/events"
	. "go-service/internal/core/interfaces"
)

type CreateUserHandler struct {
	mediator Mediator
}

func NewCreateUserHandler(mediator Mediator) *CreateUserHandler {
	return &CreateUserHandler{mediator: mediator}
}

func (uc *CreateUserHandler) Execute(ctx context.Context, command Command) (Result, error) {
	_ = command
	err := uc.mediator.PublishEvents(ctx, events.UserCreateEvent{}, events.UserCreateEvent{}, events.UserCreateEvent{}, events.UserCreateEvent{}, events.UserCreateEvent{})
	if err != nil {
		return nil, err
	}
	fmt.Println("Была выполнена команда CreateUserCommand")
	return nil, nil
}
