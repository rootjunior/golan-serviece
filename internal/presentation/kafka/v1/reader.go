package v1

import (
	"context"

	. "go-service/internal/core/interfaces"

	kafkago "github.com/segmentio/kafka-go"
)

type Handler func(ctx context.Context, msg kafkago.Message) error

type Controller struct {
	mediator Mediator
}

func NewController(m Mediator) *Controller {
	return &Controller{mediator: m}
}
