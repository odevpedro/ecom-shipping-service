package handler

import "net/http"

type ShippingHandler struct{}

func NewShippingHandler() *ShippingHandler {
	return &ShippingHandler{}
}

func (h *ShippingHandler) Calculate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"calculate - not implemented"}`))
}

func (h *ShippingHandler) Track(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"track - not implemented"}`))
}
