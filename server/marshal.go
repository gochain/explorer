package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type H map[string]interface{}

func errorResponse(w http.ResponseWriter, code int, err error) {
	bodyMap := H{"error": H{"message": err.Error()}}
	writeJSON(w, code, bodyMap)
}

func writeJSON(w http.ResponseWriter, code int, obj interface{}) { // obj map[string]interface{}) {
	jsonValue, _ := json.Marshal(obj)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err := w.Write([]byte(jsonValue))
	if err != nil {
		logger.Error("couldn't write error response", zap.Error(err))
	}
}

func writeFile(w http.ResponseWriter, code int, contentType string, file []byte) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(code)
	_, err := w.Write(file)
	if err != nil {
		logger.Error("couldn't write error response", zap.Error(err))
	}
}

func parseJSON(w http.ResponseWriter, r *http.Request, t interface{}) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(t)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid request body, bad JSON: %v", err))
	}
	return err
}
