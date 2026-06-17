package service

import (
	"errors"

	"github.com/odevpedro/ecom-shipping-service/internal/model"
)

type CorreiosCarrier struct{}

func NewCorreiosCarrier() *CorreiosCarrier {
	return &CorreiosCarrier{}
}

func (c *CorreiosCarrier) Calculate(input model.CalculateInput) (model.CalculateOutput, error) {
	return model.CalculateOutput{}, errors.New("Correios integration not implemented yet")
}

func (c *CorreiosCarrier) Track(orderID string) (model.TrackingOutput, error) {
	return model.TrackingOutput{}, errors.New("Correios integration not implemented yet")
}
