package handler

import (
	"net/http"

	"go-order-system/pkg/response"

	"github.com/gin-gonic/gin"
)

func HealthCheck(c *gin.Context) {
	response.Success(c, http.StatusOK, "Server running with custom route abstraction!", nil)
}