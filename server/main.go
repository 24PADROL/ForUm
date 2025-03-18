package main

import (
	"ForUm/auth"
	rate "ForUm/security"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	limiter := rate.NewRateLimiter(200, 60*time.Second)
	mux := http.NewServeMux()

	mux.Handle("/", limiter.Limit(http.HandlerFunc(auth.ServeHTML)))
	mux.Handle("/home", limiter.Limit(http.HandlerFunc(auth.Home)))
	mux.Handle("/register", limiter.Limit(http.HandlerFunc(auth.RegisterUser)))
	mux.Handle("/login", limiter.Limit(http.HandlerFunc(auth.LoginUser)))

	fmt.Println("✅ Serveur lancé sur https://localhost:8080") // Commande Docker :  sudo docker compose up --build

	// Load SSL certificates
	err := http.ListenAndServeTLS(":8080", "localhost.pem", "localhost-key.pem", mux)
	if err != nil {
		log.Fatalf("ListenAndServeTLS: %v", err)
	}
}
