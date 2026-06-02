package rest

import (
	"context"
	"errors"
	_ "go-service/docs"
	"go-service/internal/config"
	"go-service/internal/presentation/rest/controllers"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	cfg        config.Config
	controller *controllers.UserController
	srv        *http.Server
}

func NewServer(cfg *config.Config, c *controllers.UserController) *Server {
	return &Server{cfg: *cfg, controller: c}
}

func (s *Server) run(ctx context.Context) error {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders: []string{"Content-Length"},
		//AllowCredentials: true,
	}))
	// Каждый запрос получает контекст сервера
	r.Use(func(c *gin.Context) {
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.URL("/swagger/doc.json"),
	))
	r.POST("/posts", AuthMiddleware(), s.controller.CreateUser)

	s.srv = &http.Server{
		Addr:    s.cfg.ServerRESTAddress,
		Handler: r,
	}
	return s.srv.ListenAndServe()
}

func (s *Server) Start(ctx context.Context) {
	go func() {
		go func() {
			<-ctx.Done()
			log.Printf("Done server context")

			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := s.srv.Shutdown(shutdownCtx); err != nil {
				log.Printf("server shutdown error: %v\n", err)
			}
		}()

		if err := s.run(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("server error: %v\n", err)
		}
	}()
}
