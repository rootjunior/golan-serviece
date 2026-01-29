package main

import (
	"context"
	"database/sql"
	_ "encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	_ "some-go-project/docs"
	_ "sync"
	"syscall"

	_ "github.com/robfig/cron/v3"
	_ "gorm.io/gorm"
)

// @title SIP Bot API
// @version 1.0
// @description API для получения постов с JSONPlaceholder
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	cfg := LoadConfig()
	pool := NewWorkerPool(3)
	state := NewAppState()
	db := InitDB()
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get sqlDB:", err)
	}
	defer func(sqlDB *sql.DB) {
		err := sqlDB.Close()
		if err != nil {
			log.Fatal("Failed to close sqlDB:", err)
		}
	}(sqlDB)

	// Game
	go RunGameLoop(ctx)
	// Cron
	go StartCron(ctx)
	// HTTP
	server := HTTPServer{ctx, cfg, pool, state}
	if err := server.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
	log.Println("✅ Application stopped cleanly")
}
