package grpc

import (
	"context"
	"go-service/internal/config"
	"go-service/internal/presentation/grpc/controllers"
	userpb "go-service/internal/presentation/grpc/proto"
	"log"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	cfg        config.Config
	server     *grpc.Server
	controller *controllers.UserController
}

func NewServer(cfg *config.Config, c *controllers.UserController) *Server {
	return &Server{
		cfg:        *cfg,
		controller: c,
	}
}

func (s *Server) run() error {
	lis, err := net.Listen("tcp", s.cfg.ServerGRPCAddress)
	if err != nil {
		return err
	}

	s.server = grpc.NewServer()

	userpb.RegisterUserServiceServer(s.server, s.controller)

	return s.server.Serve(lis)
}

func (s *Server) Start(ctx context.Context) {
	go func() {
		<-ctx.Done()

		log.Println("stopping grpc server")

		s.server.GracefulStop()
	}()

	go func() {
		if err := s.run(); err != nil {
			log.Printf("grpc error: %v", err)
		}
	}()
}
