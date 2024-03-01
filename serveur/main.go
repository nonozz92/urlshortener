package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func init() {
	// Initialisation du client MongoDB.
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("Erreur lors de la connexion à MongoDB:", err)
	}
	// Vérification de la connexion
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("Échec de la connexion à MongoDB:", err)
	}
}

func main() {
	// Configuration et démarrage du serveur HTTP.
	r := SetupRouter()
	log.Fatal(http.ListenAndServe(":8000", r))
}

// enableCORS permet d'activer le support CORS pour les requêtes HTTP.
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// generateShortUrl génère une URL courte aléatoire.
func generateShortUrl() string {
    src := rand.NewSource(time.Now().UnixNano())
    rnd := rand.New(src)

    var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
    b := make([]rune, 8)
    for i := range b {
        b[i] = letters[rnd.Intn(len(letters))]
    }
    return string(b)
}

