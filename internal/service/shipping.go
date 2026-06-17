package service

import (
	"github.com/odevpedro/ecom-shipping-service/internal/model"
)

type ShippingService struct {
	carrier Carrier
}

func NewShippingService(carrier Carrier) *ShippingService {
	return &ShippingService{carrier: carrier}
}

func (s *ShippingService) Calculate(input model.CalculateInput) (model.CalculateOutput, error) {
	return s.carrier.Calculate(input)
}

func (s *ShippingService) Track(orderID string) (model.TrackingOutput, error) {
	return s.carrier.Track(orderID)
}
