package helpers

import (
	"encoding/json"
	"log"
	"net/http"
)

func ErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	var response = map[string]string{
		"message": message,
		"status":  "failed",
	}

	if statusCode > 499 {
		log.Println("Internal Server Error", message)
	}

	var jsonResponse, marshalErr = json.Marshal(response)
	if marshalErr != nil {
		ErrorResponse(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)
}

func SuccessResponse(w http.ResponseWriter, data interface{}, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	var responseBody = map[string]any{
		"message": message,
		"status":  "success",
	}

	if data != nil {
		responseBody["data"] = data
	}

	var jsonResponse, marshalErr = json.Marshal(responseBody)
	if marshalErr != nil {
		ErrorResponse(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)
}
