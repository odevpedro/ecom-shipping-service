package handler

import (
	"encoding/json"
	"errors"
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
		WriteError(w, r, "INVALID_REQUEST", "invalid request body", http.StatusBadRequest)
		return
	}

	result, err := h.svc.Calculate(input)
	if err != nil {
		WriteError(w, r, "CALCULATION_ERROR", err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *ShippingHandler) Track(w http.ResponseWriter, r *http.Request) {
	orderID := mux.Vars(r)["orderId"]
	if orderID == "" {
		WriteError(w, r, "INVALID_REQUEST", "orderId is required", http.StatusBadRequest)
		return
	}

	result, err := h.svc.Track(orderID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			WriteError(w, r, "NOT_FOUND", "tracking not found for order", http.StatusNotFound)
			return
		}
		WriteError(w, r, "TRACKING_ERROR", err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
