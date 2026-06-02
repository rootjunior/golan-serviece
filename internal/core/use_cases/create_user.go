package use_cases

import (
	"context"
	"fmt"
	"go-service/internal/core/commands"
	"go-service/internal/core/events"
	i "go-service/internal/core/interfaces"
)

type CreateUserUseCase struct {
	mediator i.IMediator
}

func NewCreateUserUseCase(mediator i.IMediator) *CreateUserUseCase {
	return &CreateUserUseCase{mediator: mediator}
}

func (uc *CreateUserUseCase) Execute(ctx context.Context, command interface{}) (interface{}, error) {
	cmd := command.(commands.CreateUserCommand)
	_ = cmd
	err := uc.mediator.PublishEvents(ctx, events.UserCreateEvent{}, events.UserCreateEvent{}, events.UserCreateEvent{}, events.UserCreateEvent{}, events.UserCreateEvent{})
	if err != nil {
		return nil, err
	}
	fmt.Println("Была выполнена команда CreateUserCommand")
	return nil, nil
}
