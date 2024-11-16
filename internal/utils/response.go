package utils

import "strconv"

// ResponseData represents the data structure for responses in the application.
type ResponseData struct {
	Message string      `json:"message"`
	Code    string      `json:"code"`
	Type    string      `json:"type"`
	Data    interface{} `json:"data,omitempty"` // Data is omitted if nil.
}

// ResponseError represents the structure for error responses.
type ResponseError struct {
	Message string              `json:"message"`
	Code    string              `json:"code"`
	Type    string              `json:"type"`
	Errors  map[string][]string `json:"errors,omitempty"` // Errors is now a map for structured validation errors.
}

// Error represents a single error detail in the application.
type Error struct {
	Error error
}

// GetErrorResponse constructs and returns a structured error response.
func (re *ResponseError) GetErrorResponse(code int, errors map[string][]string, msg string) (int, *ResponseError) {
	re.Errors = errors
	re.Code = strconv.Itoa(code)
	re.Message = msg
	re.Type = "error"

	return code, re
}

// GetErrorResponse constructs and returns an error response.
func GetResponseData(code int, data interface{}, msg string) *ResponseData {
	rd := &ResponseData{}
	rd.Code = strconv.Itoa(code)
	rd.Message = msg
	rd.Type = "success"
	rd.Data = data
	return rd
}

//// NewResponseError initializes a new ResponseError instance.
//func NewResponseError() *ResponseError {
//	return &ResponseError{
//		Errors: make(map[string][]string),
//	}
//}

// AddValidationError adds a validation error for a specific field.
func (re *ResponseError) AddValidationError(field, message string) {
	re.Errors[field] = append(re.Errors[field], message)
}
