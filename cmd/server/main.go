package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"github.com/odevpedro/ecom-shipping-service/internal/handler"
)

func main() {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3005"
	}

	r := mux.NewRouter()
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods("GET")

	shippingHandler := handler.NewShippingHandler()
	r.HandleFunc("/api/shipping/calculate", shippingHandler.Calculate).Methods("POST")
	r.HandleFunc("/api/shipping/{orderId}/track", shippingHandler.Track).Methods("GET")

	c := cors.Default()
	handler := c.Handler(r)

	log.Printf("Shipping Service running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
