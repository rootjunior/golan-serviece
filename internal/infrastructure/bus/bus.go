package bus

import (
	"go-service/internal/config"
	. "go-service/internal/core/interfaces"

	"sync"
)

type RingChannel struct {
	ch chan Event
	mu sync.Mutex
}

func (r *RingChannel) Send(v Event) {
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

func (r *RingChannel) Receive() <-chan Event {
	return r.ch
}

type ImplEventBus struct {
	*RingChannel
}

func NewEventBus(cfg *config.Config) EventBus {
	return &ImplEventBus{
		RingChannel: &RingChannel{
			ch: make(chan Event, cfg.EventBusSize),
		},
	}
}
