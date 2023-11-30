package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type RequestError struct {
	Message    string
	StatusCode int
}

func (err *RequestError) Error() string {
	return err.Message
}

func DecodeJSONBody(w http.ResponseWriter, body io.ReadCloser, dst interface{}) error {
	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()

	decodeErr := decoder.Decode(&dst)

	if decodeErr != nil {
		var syntaxError *json.SyntaxError
		var unmarshalErr *json.UnmarshalTypeError

		switch {
		case errors.As(decodeErr, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &RequestError{Message: msg, StatusCode: http.StatusBadRequest}

		case errors.Is(decodeErr, io.ErrUnexpectedEOF):
			msg := "Request body contains badly-formed JSON"
			return &RequestError{Message: msg, StatusCode: http.StatusBadRequest}

		case errors.As(decodeErr, &unmarshalErr):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalErr.Field, unmarshalErr.Offset)
			return &RequestError{Message: msg, StatusCode: http.StatusBadRequest}

		case strings.HasPrefix(decodeErr.Error(), "json: unknown field"):
			fieldName := strings.TrimPrefix(decodeErr.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &RequestError{Message: msg, StatusCode: http.StatusBadRequest}

		case errors.Is(decodeErr, io.EOF):
			msg := "Request body must not be empty"
			return &RequestError{Message: msg, StatusCode: http.StatusBadRequest}

		default:
			log.Print(decodeErr.Error())
			return &RequestError{Message: http.StatusText(http.StatusInternalServerError), StatusCode: http.StatusInternalServerError}
		}
	}

	return nil
}

func ErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := make(map[string]string)
	response["message"] = message
	response["status"] = "failed"
	jsonResponse, _ := json.Marshal(response)

	w.Write(jsonResponse)
}

func SuccessResponse(w http.ResponseWriter, data interface{}, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := make(map[string]any)
	response["message"] = message
	response["status"] = "success"
	response["data"] = data
	jsonResponse, _ := json.Marshal(response)

	w.Write(jsonResponse)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 16)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
