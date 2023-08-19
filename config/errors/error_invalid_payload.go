package errors

import "net/http"

type CustomError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e *CustomError) Error() string {
	return e.Message
}

func InvalidPayload(details string) *CustomError {
	return &CustomError{
		Code:    http.StatusBadRequest,
		Message: "Invalid payload format",
		Details: details,
	}
}
