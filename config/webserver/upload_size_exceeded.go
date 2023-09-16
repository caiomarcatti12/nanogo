package webserver

import (
	"github.com/caiomarcatti12/nanogo/v2/config/errors"
	"net/http"
	"strconv"
)

func UploadSizeExceededException(sizeLimit int64) *errors.CustomError {
	return &errors.CustomError{
		Code:    http.StatusBadRequest,
		Message: "Upload size exceeded, limit: " + strconv.FormatInt(sizeLimit, 10),
	}
}
