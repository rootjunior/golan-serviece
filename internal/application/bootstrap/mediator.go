package bootstrap

import (
	"go-service/internal/core/commands"
	"go-service/internal/core/events"
	"go-service/internal/core/handlers"
	. "go-service/internal/core/interfaces"
)

func RegisterMediator(m Mediator, createUserHandler *handlers.CreateUserHandler, userCreatedEH *handlers.UserCreatedEventHandler) {
	// Queries
	// Command
	m.RegisterCommand(commands.CreateUserCommand{}, createUserHandler)
	// Events
	m.RegisterEvent(events.UserCreateEvent{}, userCreatedEH)
}
