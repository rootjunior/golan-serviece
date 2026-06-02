package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type HTTPServer struct {
	ctx   context.Context
	cfg   Config
	pool  *WorkerPool
	state *AppState
}

func (s *HTTPServer) Start() error {
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("appState", s.state)
		c.Next()
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/posts", AuthMiddleware(), getPostsHandler)
	r.POST("/posts", createPostHandler(s.pool))

	srv := http.Server{
		Addr:    s.cfg.ServerAddress,
		Handler: r,
	}

	go func() {
		<-s.ctx.Done()
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
