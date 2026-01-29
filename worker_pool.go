package main

import (
	"log"
	"sync"
	"time"
)

type WorkerPool struct {
	queue chan PostRequest
	wg    sync.WaitGroup
}

func NewWorkerPool(workers int) *WorkerPool {
	pool := &WorkerPool{
		queue: make(chan PostRequest, 100),
	}

	for i := 0; i < workers; i++ {
		go func(id int) {
			for post := range pool.queue {
				log.Printf("🔧 Worker %d processing post: %s", id, post.Title)
				time.Sleep(5 * time.Second)
				log.Printf("✅ Worker %d done", id)
			}
		}(i)
	}

	return pool
}

func (p *WorkerPool) Enqueue(post PostRequest) {
	p.queue <- post
}
