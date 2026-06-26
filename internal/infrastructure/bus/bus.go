package bus

import (
	"go-service/internal/config"
	i "go-service/internal/core/interfaces"

	"sync"
)

type RingChannel struct {
	ch chan i.Event
	mu sync.Mutex
}

func (r *RingChannel) Send(v i.Event) {
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

func (r *RingChannel) Receive() <-chan i.Event {
	return r.ch
}

type EventBus struct {
	*RingChannel
}

func NewEventBus(cfg *config.Config) *EventBus {
	return &EventBus{
		RingChannel: &RingChannel{
			ch: make(chan i.Event, cfg.EventBusSize),
		},
	}
}
