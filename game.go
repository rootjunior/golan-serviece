package main

import (
	"context"
	"log"
	"time"
)

func RunGameLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("🛑 Game loop stopped")
			return
		default:
			log.Println("🎮 Game loop running")
			time.Sleep(500 * time.Millisecond)
		}
	}
}
