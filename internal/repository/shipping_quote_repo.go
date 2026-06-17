package repository

import (
	"database/sql"
	"time"
)

type ShippingQuote struct {
	ID            string    `json:"id"`
	FromCEP       string    `json:"from_cep"`
	ToCEP         string    `json:"to_cep"`
	WeightKg      float64   `json:"weight_kg"`
	PriceCents    int       `json:"price_cents"`
	EstimatedDays int       `json:"estimated_days"`
	Carrier       string    `json:"carrier"`
	ServiceName   string    `json:"service_name"`
	CreatedAt     time.Time `json:"created_at"`
}

func SaveQuote(db *sql.DB, q ShippingQuote) error {
	_, err := db.Exec(
		`INSERT INTO shipping_quotes (id, from_cep, to_cep, weight_kg, price_cents, estimated_days, carrier, service_name, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		q.ID, q.FromCEP, q.ToCEP, q.WeightKg, q.PriceCents, q.EstimatedDays, q.Carrier, q.ServiceName, q.CreatedAt,
	)
	return err
}

func GetQuoteByID(db *sql.DB, id string) (*ShippingQuote, error) {
	q := &ShippingQuote{}
	err := db.QueryRow(
		`SELECT id, from_cep, to_cep, weight_kg, price_cents, estimated_days, carrier, service_name, created_at
		 FROM shipping_quotes WHERE id = $1`, id,
	).Scan(&q.ID, &q.FromCEP, &q.ToCEP, &q.WeightKg, &q.PriceCents, &q.EstimatedDays, &q.Carrier, &q.ServiceName, &q.CreatedAt)
	if err != nil {
		return nil, err
	}
	return q, nil
}
