package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, err error) {
	var status int
	var msg string

	switch {
	case errors.Is(err, ErrPostNotFound):
		status = http.StatusNotFound
		msg = "Post not found"
	default:
		status = http.StatusInternalServerError
		msg = "Internal server error"
	}

	c.JSON(status, ErrorSchema{
		Code: status,
		Text: msg,
	})
}

type ErrorSchema struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

func (e *ErrorSchema) Error() string {
	return e.Text
}

var (
	ErrPostNotFound = &ErrorSchema{Code: 404, Text: "Post not found"}
)
