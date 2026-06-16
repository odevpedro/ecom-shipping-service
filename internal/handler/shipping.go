package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/odevpedro/ecom-shipping-service/internal/model"
	"github.com/odevpedro/ecom-shipping-service/internal/service"
)

type ShippingHandler struct {
	svc *service.ShippingService
}

func NewShippingHandler(svc *service.ShippingService) *ShippingHandler {
	return &ShippingHandler{svc: svc}
}

func (h *ShippingHandler) Calculate(w http.ResponseWriter, r *http.Request) {
	var input model.CalculateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	result := h.svc.Calculate(input)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *ShippingHandler) Track(w http.ResponseWriter, r *http.Request) {
	orderID := mux.Vars(r)["orderId"]
	if orderID == "" {
		http.Error(w, `{"error":"orderId is required"}`, http.StatusBadRequest)
		return
	}

	result := h.svc.Track(orderID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
