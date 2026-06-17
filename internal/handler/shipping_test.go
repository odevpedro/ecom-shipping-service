package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/odevpedro/ecom-shipping-service/internal/service"
)

func TestCalculateReturns200(t *testing.T) {
	h := newTestHandler()
	body := map[string]interface{}{
		"from_cep":  "01001000",
		"to_cep":    "02002000",
		"weight_kg": 2.5,
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/shipping/calculate", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.Calculate(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp["carrier"] != "Correios" {
		t.Errorf("expected carrier Correios, got %v", resp["carrier"])
	}
}

func TestCalculateReturns400(t *testing.T) {
	h := newTestHandler()
	req := httptest.NewRequest("POST", "/api/shipping/calculate", bytes.NewReader([]byte(`{invalid`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.Calculate(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestTrackReturns200(t *testing.T) {
	h := newTestHandler()
	req := httptest.NewRequest("GET", "/api/shipping/order-123/track", nil)
	rec := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/api/shipping/{orderId}/track", h.Track).Methods("GET")
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp["order_id"] != "order-123" {
		t.Errorf("expected order_id order-123, got %v", resp["order_id"])
	}
}

func TestTrackReturns404(t *testing.T) {
	h := newTestHandler()
	req := httptest.NewRequest("GET", "/api/shipping/nonexistent/track", nil)
	rec := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/api/shipping/{orderId}/track", h.Track).Methods("GET")
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}

func newTestHandler() *ShippingHandler {
	carrier := service.NewStubCarrier()
	svc := service.NewShippingService(carrier)
	return NewShippingHandler(svc)
}
