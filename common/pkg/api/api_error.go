package api

import "fmt"

type APIError struct {
	API    string
	Method string
	Code   int
	Msg    string
}

func (e APIError) Error() string {
	return fmt.Sprintf("error when %s API %s, code: %d, msg: %s", e.Method, e.API, e.Code, e.Msg)
}
