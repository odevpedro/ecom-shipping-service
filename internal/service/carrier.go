package service

import (
	"errors"

	"github.com/odevpedro/ecom-shipping-service/internal/model"
)

var ErrNotFound = errors.New("not found")

type Carrier interface {
	Calculate(input model.CalculateInput) (model.CalculateOutput, error)
	Track(orderID string) (model.TrackingOutput, error)
}
