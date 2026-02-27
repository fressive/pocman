package model

// Response represents a standardized JSON response structure.
// Code: indicates the application-level status code (0 for success, non-zero for errors).
// Msg: contains a human-readable message describing the response.
// Data: holds the response payload, supporting any data type.
type Response[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}
