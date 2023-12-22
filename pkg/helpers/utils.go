package helpers

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
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

func HashApiKey(apiKey string) string {
	h := sha256.New()
	h.Write([]byte(apiKey))
	bs := h.Sum(nil)

	return base64.URLEncoding.EncodeToString(bs)
}

func GetEnv(key string) string {
	var value = os.Getenv(key)
	if value == "" {
		log.Fatal("Environment variable not set: " + key)
	}

	return value
}
