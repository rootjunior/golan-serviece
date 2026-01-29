package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Get posts
// @Description Получить список постов (требуется авторизация)
// @Tags posts
// @Produce json
// @Success 200 {array} PostDB "Успешный ответ"
// @Failure 400 {object} ErrorSchema "Неверный запрос"
// @Failure 401 {object} ErrorSchema "Неавторизован"
// @Failure 403 {object} ErrorSchema "Запрещено"
// @Failure 500 {object} ErrorSchema "Внутренняя ошибка сервера"
// @Security BearerAuth
// @Router /posts [get]
func getPostsHandler(c *gin.Context) {
	m := c.MustGet("appState").(*AppState).Mediator
	posts, err := m.Execute(GetPostsQuery{})
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, posts)
}

// @Summary Create post (background processing)
// @Tags posts
// @Accept json
// @Produce json
// @Param post body PostRequest true "Post body"
// @Success 200 {object} PostRequest
// @Router /posts [post]
func createPostHandler(pool *WorkerPool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req PostRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorSchema{
				Code: http.StatusBadRequest,
				Text: err.Error(),
			})
			return
		}

		fmt.Printf("📬 POST получен: %+v\n", req)

		// 🔥 ФОНОВАЯ ЗАДАЧА
		pool.Enqueue(req)

		// ⚡ Ответ сразу
		c.JSON(http.StatusOK, gin.H{
			"status": "post accepted",
			"async":  true,
		})
	}
}
