package errors

import "net/http"

func InternalServerError(details string) *CustomError {
	return &CustomError{
		Code:    http.StatusInternalServerError,
		Message: "Internal Server Error",
		Details: details,
	}
}
