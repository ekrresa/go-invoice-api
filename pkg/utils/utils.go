package utils

import (
	"encoding/json"
	"net/http"
)

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

	response := make(map[string]string)
	response["message"] = message
	response["status"] = "success"
	jsonResponse, _ := json.Marshal(data)

	w.Write(jsonResponse)
}
