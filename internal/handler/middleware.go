package handler

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Data  interface{} `json:"data"`
	Error *ErrorObj   `json:"error"`
	Meta  Meta        `json:"meta"`
}

type ErrorObj struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details"`
}

type Meta struct {
	RequestID string `json:"requestId"`
}

func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = generateID()
		}
		w.Header().Set("X-Request-ID", requestID)
		next.ServeHTTP(w, r)
	})
}

func WriteError(w http.ResponseWriter, r *http.Request, code string, message string, status int) {
	requestID := r.Header.Get("X-Request-ID")
	if requestID == "" {
		requestID = generateID()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Data:  nil,
		Error: &ErrorObj{Code: code, Message: message, Details: struct{}{}},
		Meta:  Meta{RequestID: requestID},
	})
}
