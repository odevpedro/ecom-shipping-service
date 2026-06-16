package model

type CalculateInput struct {
	FromCEP       string  `json:"from_cep"`
	ToCEP         string  `json:"to_cep"`
	WeightKg      float64 `json:"weight_kg"`
	HeightCm      float64 `json:"height_cm"`
	WidthCm       float64 `json:"width_cm"`
	LengthCm      float64 `json:"length_cm"`
}

type CalculateOutput struct {
	Carrier       string `json:"carrier"`
	ServiceName   string `json:"service_name"`
	PriceCents    int    `json:"price_cents"`
	EstimatedDays int    `json:"estimated_days"`
	Currency      string `json:"currency"`
}

type TrackingEvent struct {
	Date        string `json:"date"`
	Location    string `json:"location"`
	Description string `json:"description"`
}

type TrackingOutput struct {
	OrderID  string           `json:"order_id"`
	Carrier  string           `json:"carrier"`
	Status   string           `json:"status"`
	Events   []TrackingEvent  `json:"events"`
}
