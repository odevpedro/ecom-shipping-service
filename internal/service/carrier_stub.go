package service

import (
	"math"
	"strconv"

	"github.com/odevpedro/ecom-shipping-service/internal/model"
)

type StubCarrier struct{}

func NewStubCarrier() *StubCarrier {
	return &StubCarrier{}
}

func (s *StubCarrier) Calculate(input model.CalculateInput) (model.CalculateOutput, error) {
	distance := estimateDistance(input.FromCEP, input.ToCEP)
	weightCharge := int(input.WeightKg * 50)
	distanceCharge := int(distance * 1)
	priceCents := weightCharge + distanceCharge
	days := estimateDays(distance)

	return model.CalculateOutput{
		Carrier:       "Correios",
		ServiceName:   "PAC",
		PriceCents:    priceCents,
		EstimatedDays: days,
		Currency:      "BRL",
	}, nil
}

func (s *StubCarrier) Track(orderID string) (model.TrackingOutput, error) {
	if orderID != "order-123" {
		return model.TrackingOutput{}, ErrNotFound
	}

	return model.TrackingOutput{
		OrderID: orderID,
		Carrier: "Correios",
		Status:  "in_transit",
		Events: []model.TrackingEvent{
			{Date: "2026-06-14 08:30", Location: "São Paulo, SP", Description: "Objeto postado"},
			{Date: "2026-06-15 14:15", Location: "Curitiba, PR", Description: "Em trânsito para unidade de distribuição"},
			{Date: "2026-06-16 09:00", Location: "Curitiba, PR", Description: "Saiu para entrega ao destinatário"},
		},
	}, nil
}

func estimateDistance(from, to string) float64 {
	f := extractPrefix(from)
	t := extractPrefix(to)
	diff := math.Abs(float64(f - t))
	if diff < 1 {
		diff = 1
	}
	return diff * 50
}

func estimateDays(distance float64) int {
	days := int(distance / 200)
	if days < 1 {
		return 1
	}
	if days > 15 {
		return 15
	}
	return days
}

func extractPrefix(cep string) int {
	if len(cep) < 5 {
		return 0
	}
	v, err := strconv.Atoi(cep[:5])
	if err != nil {
		return 0
	}
	return v
}
