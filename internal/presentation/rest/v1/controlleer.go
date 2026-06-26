package v1

import . "go-service/internal/core/interfaces"

type Controller struct {
	mediator Mediator
}

func NewController(m Mediator) *Controller {
	return &Controller{mediator: m}
}
