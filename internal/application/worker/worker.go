package worker

import (
	"context"
	"fmt"
	"go-service/internal/config"
	. "go-service/internal/core/interfaces"
	"sync"
)

type ImplWorkerPool struct {
	mediator     Mediator
	bus          EventBus
	mu           sync.RWMutex
	workersCount int
}

func NewWorkerPool(cfg *config.Config, mediator Mediator, bus EventBus) WorkerPool {
	return &ImplWorkerPool{
		mediator:     mediator,
		bus:          bus,
		workersCount: cfg.WorkersCount,
	}
}

func (w *ImplWorkerPool) StartProcessEvents(ctx context.Context) {
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
