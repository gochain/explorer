package main

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type H map[string]interface{}

func errorResponse(w http.ResponseWriter, code int, err error) {
	bodyMap := H{"error": H{"message": err.Error()}}
	writeJSON(w, code, bodyMap)
}

func writeJSON(w http.ResponseWriter, code int, obj interface{}) {
	jsonValue, err := json.Marshal(obj)
	if err != nil {
		logger.Error("Failed to marshal JSON response", zap.Error(err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(jsonValue)
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
