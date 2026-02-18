package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Message string `json:"message" example:"not implemented"`
}

func notImplemented(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, ErrorResponse{Message: "not implemented"})
}
