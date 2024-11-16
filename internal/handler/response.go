package handler

import "github.com/kaium123/order/internal/errors"

// ResponseData is the response structure for the application.
type ResponseData struct {
	// Data is the response data.
	Data interface{} `json:"data,omitempty"`
}

// ResponseError is the response structure for the application.
type ResponseError struct {
	// Errors is the response errors.
	Errors []Error `json:"errors,omitempty"`
}

// Error is the error structure for the application.
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (re *ResponseError) GetErrorResponse(code int, err error) (int, *ResponseError) {
	re.Errors = append(re.Errors, Error{
		Code:    errors.ErrorCodeDescriptions[code],
		Message: err.Error(),
	})

	return code, re
}
