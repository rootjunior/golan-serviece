package use_cases

import (
	"context"
	"fmt"
	"go-service/internal/core/events"
	. "go-service/internal/core/interfaces"
)

type CreateUserUseCase struct {
	mediator IMediator
}

func NewCreateUserUseCase(mediator IMediator) *CreateUserUseCase {
	return &CreateUserUseCase{mediator: mediator}
}

func (uc *CreateUserUseCase) Execute(ctx context.Context, command Command) (Result, error) {
	_ = command
	err := uc.mediator.PublishEvents(ctx, events.UserCreateEvent{}, events.UserCreateEvent{}, events.UserCreateEvent{}, events.UserCreateEvent{}, events.UserCreateEvent{})
	if err != nil {
		return nil, err
	}
	fmt.Println("Была выполнена команда CreateUserCommand")
	return nil, nil
}
