package main

import (
	"github.com/gorilla/mux"
)

// SetupRouter initialise et configure le routeur avec les différents gestionnaires de routes.
func SetupRouter() *mux.Router {
    r := mux.NewRouter()

    // Application du middleware CORS à toutes les routes
    r.Use(enableCORS)

    // Configuration des gestionnaires de routes pour différentes fonctionnalités
    r.HandleFunc("/register", registerHandler).Methods("POST", "OPTIONS") 
    r.HandleFunc("/login", loginHandler).Methods("POST", "OPTIONS")     
    r.HandleFunc("/shorten", shortenUrlHandler).Methods("POST", "OPTIONS", "PUT", "DELETE", "GET") 
    r.HandleFunc("/resolve", resolveShortUrlHandler).Methods("POST", "OPTIONS", "PUT", "DELETE", "GET") 
    r.HandleFunc("/link-stats", getAllUrlsHandler).Methods("GET")

    return r
}
