package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"

	"github.com/odevpedro/ecom-shipping-service/internal/config"
	"github.com/odevpedro/ecom-shipping-service/internal/handler"
	"github.com/odevpedro/ecom-shipping-service/internal/repository"
	"github.com/odevpedro/ecom-shipping-service/internal/service"
)

func main() {
	godotenv.Load()
	cfg := config.Load()

	var db *sql.DB
	if cfg.DatabaseURL != "" {
		var err error
		db, err = repository.NewPostgres(cfg.DatabaseURL)
		if err != nil {
			log.Printf("WARNING: could not connect to database: %v", err)
		} else {
			defer db.Close()
			log.Println("connected to PostgreSQL")
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
	handler := c.Handler(r)

	log.Printf("Shipping Service running on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, handler))
}
