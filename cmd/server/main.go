package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"github.com/odevpedro/ecom-shipping-service/internal/handler"
	"github.com/odevpedro/ecom-shipping-service/internal/service"
)

func main() {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3005"
	}

	shippingSvc := service.NewShippingService()
	shippingHandler := handler.NewShippingHandler(shippingSvc)

	r := mux.NewRouter()
	r.Use(handler.RequestIDMiddleware)

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok","service":"shipping"}`))
	}).Methods("GET")
	r.HandleFunc("/live", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"alive"}`))
	}).Methods("GET")
	r.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ready"}`))
	}).Methods("GET")

	r.HandleFunc("/api/shipping/calculate", shippingHandler.Calculate).Methods("POST")
	r.HandleFunc("/api/shipping/{orderId}/track", shippingHandler.Track).Methods("GET")

	c := cors.Default()
	handler := c.Handler(r)

	log.Printf("Shipping Service running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
