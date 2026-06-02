package bus

import (
	"go-service/internal/config"
	"sync"
)

type RingChannel struct {
	ch chan interface{}
	mu sync.Mutex
}

func (r *RingChannel) Send(v interface{}) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.ch) == cap(r.ch) {
		select {
		case <-r.ch:
		default:
		}
	}
	r.ch <- v
}

func (r *RingChannel) Close() {
	close(r.ch)
}

func (r *RingChannel) Receive() <-chan interface{} {
	return r.ch
}

type EventBus struct {
	*RingChannel
}

func NewEventBus(cfg *config.Config) *EventBus {
	return &EventBus{
		RingChannel: &RingChannel{
			ch: make(chan interface{}, cfg.EventBusSize),
		},
	}
}
