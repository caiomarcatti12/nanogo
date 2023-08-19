package errors

import "net/http"


func InvalidPayload(details string) *CustomError {
	return &CustomError{
		Code:    http.StatusBadRequest,
		Message: "Invalid payload format",
		Details: details,
	}
}
