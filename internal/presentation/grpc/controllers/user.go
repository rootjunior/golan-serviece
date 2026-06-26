package controllers

import (
	"context"
	"go-service/internal/core/commands"
	. "go-service/internal/core/interfaces"
	pb "go-service/internal/presentation/grpc/proto"
)

type UserController struct {
	pb.UnimplementedUserServiceServer
	mediator Mediator
}

func NewUserController(mediator Mediator) *UserController {
	return &UserController{
		mediator: mediator,
	}
}

func (h *UserController) CreateUser(
	ctx context.Context,
	req *pb.CreateUserRequest,
) (*pb.CreateUserResponse, error) {

	_, err := h.mediator.HandleCommand(
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
