package main

import (
	"context"
	_ "encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	_ "some-go-project/docs"
	_ "sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/robfig/cron/v3"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

// ================= DI =================

// @title SIP Bot API
// @version 1.0
// @description API для получения постов с JSONPlaceholder
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func runHTTP(ctx context.Context, cfg Config, pool *WorkerPool, state *AppState) error {
	r := gin.Default()
	// Прокидываем AppState в контекст Gin
	r.Use(func(c *gin.Context) {
		c.Set("appState", state)
		c.Next()
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/posts", AuthMiddleware(), getPostsHandler)
	r.POST("/posts", createPostHandler(pool))

	srv := http.Server{
		Addr:    cfg.ServerAddress,
		Handler: r,
	}

	go func() {
		<-ctx.Done()
		log.Println("🛑 HTTP shutting down...")
		timeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := srv.Shutdown(timeout)
		if err != nil {
			return
		}
	}()

	log.Println("🚀 HTTP server started")
	return srv.ListenAndServe()
}

/* ================= MAIN ================= */

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	cfg := LoadConfig()
	pool := NewWorkerPool(3)
	// Создаем AppState
	state := NewAppState()

	//// Game loop (blocking)
	go func() {
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
	}()
	InitDB()

	// Cron
	go startCron(ctx)

	// HTTP
	if err := runHTTP(ctx, cfg, pool, state); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}

	log.Println("✅ Application stopped cleanly")
}
