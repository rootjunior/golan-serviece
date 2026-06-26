package v1

import (
	"fmt"
	"go-service/internal/core/commands"
	"go-service/internal/presentation/rest/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateUser @Summary Create user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user body commands.CreateUserCommand true "Create user body"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} errors.ErrorSchema
// @Failure 500 {object} errors.ErrorSchema
// @Router /posts [post]
func (h *Controller) CreateUser(ctx *gin.Context) {
	posts, err := h.mediator.ExecuteCommand(ctx, commands.CreateUserCommand{})
	if err != nil {
		fmt.Println(err)
		errors.HandleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, posts)
}
