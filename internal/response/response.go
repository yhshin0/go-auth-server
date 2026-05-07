package response

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/yhshin0/go-auth-server/internal/defs"
)

type Response struct {
	Code    string `json:"code"` // e.g. SUCCESS, INVALID_CREDENTIALS
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

func GeneralResponse(w http.ResponseWriter, statusCode int, data any) {
	const success = "SUCCESS"

	write(w, statusCode, Response{Code: success, Data: data})
}

func ErrorResponse(w http.ResponseWriter, err error) {
	const internalServerError = "INTERNAL_SERVER_ERROR"

	statusCode := http.StatusInternalServerError
	code := internalServerError
	message := http.StatusText(http.StatusInternalServerError)

	logLevel := slog.LevelError
	var e defs.CustomError
	if errors.As(err, &e) {
		statusCode = e.StatusCode
		code = e.Code
		message = e.Message
		logLevel = slog.LevelWarn
	}
	slog.Log(context.Background(), logLevel, err.Error())

	write(w, statusCode, Response{Code: code, Message: message})
}

func write(w http.ResponseWriter, statusCode int, resp Response) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	b, _ := json.Marshal(resp)
	w.Write(b)
}

func NoContentResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
