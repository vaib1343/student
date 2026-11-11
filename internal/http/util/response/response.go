package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

const (
	StatusOK    = "ok"
	StatusError = "error"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func GeneralError(err error) Response {
	return Response{
		Status: StatusError,
		Error:  err.Error(),
	}
}

func ValidationError(errs validator.ValidationErrors) Response {
	var errMsg []string
	for _, err := range errs {
		switch err.Tag() {
		case "required":
			errMsg = append(errMsg, fmt.Sprintf("%s is required", err.Field()))
		default:
			errMsg = append(errMsg, fmt.Sprintf("%s is not valid", err.Field()))
		}
	}

	return Response{
		Status: StatusError,
		Error:  strings.Join(errMsg, ", "),
	}
}
