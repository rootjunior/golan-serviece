package controllers

import (
	"context"
	"go-service/internal/core/commands"
	i "go-service/internal/core/interfaces"
	pb "go-service/internal/presentation/grpc/proto"
)

type UserController struct {
	pb.UnimplementedUserServiceServer
	mediator i.IMediator
}

func NewUserController(mediator i.IMediator) *UserController {
	return &UserController{
		mediator: mediator,
	}
}

func (h *UserController) CreateUser(
	ctx context.Context,
	req *pb.CreateUserRequest,
) (*pb.CreateUserResponse, error) {

	_, err := h.mediator.ExecuteCommand(
		ctx,
		commands.CreateUserCommand{},
	)

	if err != nil {
		return nil, err
	}

	return &pb.CreateUserResponse{
		Message: "user created",
	}, nil
}
