package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
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
		log.Error().Err(err).Msg("couldn't write error response")
	}
}

func writeFile(w http.ResponseWriter, code int, contentType string, file []byte) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(code)
	_, err := w.Write(file)
	if err != nil {
		log.Error().Err(err).Msg("couldn't write error response")
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
