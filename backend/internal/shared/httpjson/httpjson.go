package httpjson

import "github.com/gin-gonic/gin"

type ErrorResponse struct {
	Error   string      `json:"error"`
	Code    string      `json:"code,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

func Error(c *gin.Context, status int, code, message string) {
	c.JSON(status, ErrorResponse{
		Error: message,
		Code:  code,
	})
}

func ErrorWithDetails(c *gin.Context, status int, code, message string, details interface{}) {
	c.JSON(status, ErrorResponse{
		Error:   message,
		Code:    code,
		Details: details,
	})
}
