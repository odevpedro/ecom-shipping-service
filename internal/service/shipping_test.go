package service

import (
	"testing"

	"github.com/odevpedro/ecom-shipping-service/internal/model"
)

func TestCalculateReturnsValues(t *testing.T) {
	svc := NewShippingService()
	input := model.CalculateInput{
		FromCEP:  "01001000",
		ToCEP:    "20020000",
		WeightKg: 2.5,
	}
	result := svc.Calculate(input)

	if result.Carrier == "" {
		t.Error("expected carrier to be set")
	}
	if result.PriceCents <= 0 {
		t.Errorf("expected price > 0, got %d", result.PriceCents)
	}
	if result.EstimatedDays <= 0 {
		t.Errorf("expected days > 0, got %d", result.EstimatedDays)
	}
	if result.Currency != "BRL" {
		t.Errorf("expected BRL, got %s", result.Currency)
	}
}

func TestCalculatePricing(t *testing.T) {
	svc := NewShippingService()

	light := svc.Calculate(model.CalculateInput{FromCEP: "01001000", ToCEP: "02002000", WeightKg: 1})
	heavy := svc.Calculate(model.CalculateInput{FromCEP: "01001000", ToCEP: "02002000", WeightKg: 10})

	if heavy.PriceCents <= light.PriceCents {
		t.Error("expected heavy to cost more than light")
	}
}

func TestCalculateDistance(t *testing.T) {
	svc := NewShippingService()

	nearby := svc.Calculate(model.CalculateInput{FromCEP: "01001000", ToCEP: "01002000", WeightKg: 1})
	far := svc.Calculate(model.CalculateInput{FromCEP: "01001000", ToCEP: "30000000", WeightKg: 1})

	if far.EstimatedDays <= nearby.EstimatedDays {
		t.Error("expected far delivery to take longer than nearby")
	}
}

func TestTrackReturnsEvents(t *testing.T) {
	svc := NewShippingService()
	result := svc.Track("order-123")

	if result.OrderID != "order-123" {
		t.Errorf("expected order-123, got %s", result.OrderID)
	}
	if len(result.Events) == 0 {
		t.Error("expected at least one tracking event")
	}
	if result.Status != "in_transit" {
		t.Errorf("expected in_transit, got %s", result.Status)
	}
}
