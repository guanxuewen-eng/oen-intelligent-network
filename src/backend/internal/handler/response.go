package handler

import "github.com/gin-gonic/gin"

// ErrorResponse represents a unified error response format
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// SuccessResponse represents a unified success response format
type SuccessResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// OK sends a success response
func OK(c *gin.Context, data interface{}) {
	c.JSON(200, SuccessResponse{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// OKWithMessage sends a success response with a custom message
func OKWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(200, SuccessResponse{
		Code:    0,
		Message: message,
		Data:    data,
	})
}

// BadRequest sends a 400 error response
func BadRequest(c *gin.Context, message string) {
	c.JSON(400, ErrorResponse{
		Code:    400,
		Message: message,
	})
}

// NotFound sends a 404 error response
func NotFound(c *gin.Context, message string) {
	c.JSON(404, ErrorResponse{
		Code:    404,
		Message: message,
	})
}

// InternalError sends a 500 error response
func InternalError(c *gin.Context, message string) {
	c.JSON(500, ErrorResponse{
		Code:    500,
		Message: message,
	})
}

// Error sends a custom error response
func Error(c *gin.Context, code int, message string, details string) {
	c.JSON(code, ErrorResponse{
		Code:    code,
		Message: message,
		Details: details,
	})
}
