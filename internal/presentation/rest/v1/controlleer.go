package v1

import i "go-service/internal/core/interfaces"

type Controller struct {
	mediator i.IMediator
}

func NewController(m i.IMediator) *Controller {
	return &Controller{mediator: m}
}
