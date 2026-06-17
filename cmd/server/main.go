package main

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"

	"github.com/odevpedro/ecom-shipping-service/internal/config"
	"github.com/odevpedro/ecom-shipping-service/internal/handler"
	"github.com/odevpedro/ecom-shipping-service/internal/repository"
	"github.com/odevpedro/ecom-shipping-service/internal/service"
)

func autoMigrate(db *sql.DB) error {
	ddl := `
	CREATE TABLE IF NOT EXISTS shipping_quotes (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		from_cep VARCHAR(8) NOT NULL,
		to_cep VARCHAR(8) NOT NULL,
		weight_kg DOUBLE PRECISION,
		price_cents INTEGER NOT NULL,
		estimated_days INTEGER NOT NULL,
		carrier VARCHAR(100),
		service_name VARCHAR(100),
		created_at TIMESTAMPTZ DEFAULT NOW()
	);
	CREATE TABLE IF NOT EXISTS tracking_events (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		order_id VARCHAR(100) NOT NULL,
		location VARCHAR(255),
		description TEXT,
		event_date TIMESTAMPTZ,
		created_at TIMESTAMPTZ DEFAULT NOW()
	);`
	_, err := db.Exec(ddl)
	return err
}

func main() {
	godotenv.Load()
	cfg := config.Load()
	logger := config.NewLogger()

	var db *sql.DB
	if cfg.DatabaseURL != "" {
		var err error
		db, err = repository.NewPostgres(cfg.DatabaseURL)
		if err != nil {
			logger.Warn("could not connect to database", "error", err)
		} else {
			defer db.Close()
			logger.Info("connected to PostgreSQL")

			if err := autoMigrate(db); err != nil {
				logger.Warn("auto-migration failed", "error", err)
			}
		}
	}
	_ = db

	carrier := service.NewStubCarrier()
	shippingSvc := service.NewShippingService(carrier)
	shippingHandler := handler.NewShippingHandler(shippingSvc)

	r := mux.NewRouter()
	r.Use(handler.RequestIDMiddleware)

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok", "service": "shipping"})
	}).Methods("GET")
	r.HandleFunc("/live", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "alive"})
	}).Methods("GET")
	r.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ready"})
	}).Methods("GET")

	r.HandleFunc("/api/shipping/calculate", shippingHandler.Calculate).Methods("POST")
	r.HandleFunc("/api/shipping/{orderId}/track", shippingHandler.Track).Methods("GET")

	c := cors.Default()
	h := c.Handler(r)

	logger.Info("Shipping Service started", "port", cfg.Port)
	slog.SetDefault(logger)
	if err := http.ListenAndServe(":"+cfg.Port, h); err != nil {
		logger.Error("server failed to start", "error", err)
		os.Exit(1)
	}
}
