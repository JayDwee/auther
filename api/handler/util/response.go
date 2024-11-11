package util

import (
	"encoding/json"
	"net/http"
)

func JsonResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)

	if errResponse, ok := data.(error); ok {
		// handle error
		data = struct {
			Error string `json:"error"`
		}{Error: errResponse.Error()}
	}
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		panic(err)
	}
}

func JsonResponseNoCache(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	JsonResponse(w, statusCode, data)
}
