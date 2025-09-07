package utils

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/bytedance/sonic"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

type H map[string]any

type ErrorResponse struct {
	Errors []string `json:"errors"`
}

func BindJSON(r *http.Request, v any) error {
	return sonic.ConfigDefault.NewDecoder(r.Body).Decode(v)
}

func Write(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := sonic.ConfigDefault.NewEncoder(w).Encode(v); err != nil {
		log.Error().Err(err).Msg("Failed to encode response")
	}
}

func WriteError(w http.ResponseWriter, status int, msg string) {
	WriteJSON(w, status, ErrorResponse{Errors: []string{msg}})
}

func CloseBody(body io.ReadCloser) {
	if body == nil {
		return
	}

	if err := body.Close(); err != nil {
		log.Error().Err(err).Msg("Error closing body")
	}
}

func Validate(v any) *ErrorResponse {
	err := validate.Struct(v)
	if err == nil {
		return nil
	}

	errResponse := &ErrorResponse{}

	var validateErrs validator.ValidationErrors
	if !errors.As(err, &validateErrs) {
		errResponse.Errors = append(errResponse.Errors, err.Error())
		return errResponse
	}

	for _, e := range validateErrs {
		field := e.Field()
		tag := e.Tag()
		errResponse.Errors = append(errResponse.Errors, fmt.Sprintf("%s failed on %s", field, tag))
	}

	return errResponse
}
