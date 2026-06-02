package worker

import (
	"context"
	"fmt"
	"go-service/internal/config"
	i "go-service/internal/core/interfaces"
	"go-service/internal/infrastructure/bus"
	"sync"
)

type WorkerPool struct {
	mediator     i.IMediator
	bus          *bus.EventBus
	mu           sync.RWMutex
	workersCount int
}

func NewWorkerPool(cfg *config.Config, mediator i.IMediator, bus *bus.EventBus) *WorkerPool {
	return &WorkerPool{
		mediator:     mediator,
		bus:          bus,
		workersCount: cfg.WorkersCount,
	}
}

func (w *WorkerPool) StartProcessEvents(ctx context.Context) {
	var wg sync.WaitGroup

	for range w.workersCount {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case event, ok := <-w.bus.Receive():
					if !ok {
						return
					}
					if err := w.mediator.HandleEvent(ctx, event); err != nil {
						fmt.Printf("dispatch error: %v\n", err)
					}
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		w.bus.Close()
	}()
}
