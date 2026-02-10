// Package response provides utilities for generating standardized HTTP responses in the Gin framework.
// It defines a Response structure and helper functions to construct consistent JSON responses
// with code, message, and data fields.
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response represents a standardized JSON response structure.
// Code: indicates the application-level status code (0 for success, non-zero for errors).
// Msg: contains a human-readable message describing the response.
// Data: holds the response payload, supporting any data type via interface{}.
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Result sends a JSON response to the client with the specified HTTP status code,
// application code, message, and data payload.
// c: the Gin context for the request.
// httpCode: the HTTP status code (e.g., http.StatusOK).
// code: the application-level status code.
// msg: a message describing the response.
// data: the response payload (can be nil).
func Result(c *gin.Context, httpCode int, code int, msg string, data any) {
	c.JSON(httpCode, Response{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}

// Success sends a successful JSON response with HTTP 200 status, application code 0,
// and the provided data payload.
// c: the Gin context for the request.
// data: the response payload to return to the client.
func Success(c *gin.Context, data interface{}) {
	Result(c, http.StatusOK, 0, "success", data)
}

// Error sends an error JSON response with HTTP 200 status, the specified application code and message.
// c: the Gin context for the request.
// code: the application-level error code.
// msg: a message describing the error.
func Error(c *gin.Context, code int, msg string) {
	Result(c, http.StatusBadRequest, code, msg, nil)
}
