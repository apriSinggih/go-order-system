package response

import (

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func Success(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, Response{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, code int, err string) {
	c.JSON(code, Response{
		Status: "error",
		Error:  err,
	})
}