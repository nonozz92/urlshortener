package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

var mySigningKey = []byte("secret") // Clé secrète pour signer les JWT.

// registerHandler gère l'inscription des nouveaux utilisateurs.
func registerHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Hachage du mot de passe utilisateur à l'aide de bcrypt.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Erreur lors du hashage du mot de passe", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	// Insertion de l'utilisateur dans la base de données.
	collection := client.Database("urlshortener").Collection("users")
	_, err = collection.InsertOne(context.TODO(), user)
	if err != nil {
		http.Error(w, "Erreur lors de l'enregistrement de l'utilisateur", http.StatusInternalServerError)
		return
	}

	// Réponse indiquant la création réussie.
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Compte créé avec succès!"})
}

// loginHandler gère la connexion des utilisateurs.
func loginHandler(w http.ResponseWriter, r *http.Request) {
	var user, foundUser User
	// Décodage du corps de la requête JSON.
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Recherche de l'utilisateur dans la base de données.
	collection := client.Database("urlshortener").Collection("users")
	err = collection.FindOne(context.TODO(), bson.M{"username": user.Username}).Decode(&foundUser)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password)) != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Création du token JWT pour l'utilisateur authentifié.
	claims := jwt.MapClaims{
		"userId": foundUser.ID,
		"exp":    time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Envoi du token JWT à l'utilisateur.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}


// shortenUrlHandler traite les requêtes pour créer une URL courte.
func shortenUrlHandler(w http.ResponseWriter, r *http.Request) {
    // Définir la structure pour les données de la requête.
    var requestData requestData

    // Décoder le corps de la requête en JSON.
    if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
        http.Error(w, "Données de requête invalides", http.StatusBadRequest)
        return
    }

    // Extraire la base de l'URL longue et générer l'URL courte.
    base := extractBaseURL(requestData.LongUrl)
    idUnique := generateShortUrl()
    shortUrl := fmt.Sprintf("%s/%s", base, idUnique)

    // Valider le token utilisateur et extraire l'ID de l'utilisateur.
    userID := ""
    if requestData.UserToken != "" {
        var err error
        userID, err = validateTokenAndGetUserID(requestData.UserToken)
        if err != nil {
            http.Error(w, "Token invalide", http.StatusUnauthorized)
            return
        }
    }

    // Insérer l'URL raccourcie dans la base de données.
    urlCollection := client.Database("urlshortener").Collection("urls")
    _, err := urlCollection.InsertOne(context.TODO(), bson.M{
        "longUrl":  requestData.LongUrl,
        "shortUrl": shortUrl,
        "userId":   userID,
    })
    if err != nil {
        http.Error(w, "Erreur serveur", http.StatusInternalServerError)
        return
    }

    // Envoyer l'URL raccourcie au client.
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"shortUrl": shortUrl})
}


// extractBaseURL extrait la partie base de l'URL fournie.
func extractBaseURL(urlStr string) string {
    parsedUrl, err := url.Parse(urlStr)
    if err != nil {
        return ""
    }
    return fmt.Sprintf("%s://%s", parsedUrl.Scheme, parsedUrl.Host)
}

// validateTokenAndGetUserID valide le token JWT et extrait l'ID de l'utilisateur.
func validateTokenAndGetUserID(tokenString string) (string, error) {
    if tokenString == "" {
        return "", errors.New("token manquant")
    }

    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("méthode de signature inattendue : %v", token.Header["alg"])
        }
        return mySigningKey, nil
    })

    if err != nil {
        return "", err
    }

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        userID, ok := claims["userId"].(string)
        if !ok {
            return "", errors.New("le champ userId est manquant ou n'est pas une chaîne")
        }
        return userID, nil
    }
    return "", errors.New("token invalide")
}

// resolveShortUrlHandler traite les requêtes pour obtenir l'URL longue originale à partir de l'URL courte.
func resolveShortUrlHandler(w http.ResponseWriter, r *http.Request) {
    urlCollection := client.Database("urlshortener").Collection("urls")

    // Extraire le token d'autorisation de l'en-tête de la requête.
    authHeader := r.Header.Get("Authorization")
    token := strings.TrimPrefix(authHeader, "Bearer ")
    var userID string
    var err error

    // Valider le token et extraire l'ID de l'utilisateur.
    if token != "" {
        userID, err = validateTokenAndGetUserID(token)
        if err != nil {
            log.Printf("Erreur lors de la validation du token : %v", err)
        }
    }

    // Récupérer l'URL courte à partir des paramètres de la requête.
    queryValues := r.URL.Query()
    shortUrl := queryValues.Get("url")
    if shortUrl == "" {
        http.Error(w, "Paramètre d'URL manquant", http.StatusBadRequest)
        return
    }

    // Trouver l'URL longue correspondante dans la base de données.
    var urlData urlData
    
    err = urlCollection.FindOne(context.TODO(), bson.M{"shortUrl": shortUrl}).Decode(&urlData)
    if err != nil {
        http.Error(w, "URL non trouvée", http.StatusNotFound)
        return
    }

    // Vérifier si l'utilisateur a le droit d'accéder à cette URL.
    if urlData.UserID != "" && (userID == "" || urlData.UserID != userID) {
        http.Error(w, "Accès refusé", http.StatusForbidden)
        return
    }

    // Mettre à jour le compteur de visites pour l'URL.
    _, err = urlCollection.UpdateOne(context.TODO(), bson.M{"shortUrl": shortUrl}, bson.M{"$inc": bson.M{"visitCount": 1}})
    if err != nil {
        http.Error(w, "Erreur serveur", http.StatusInternalServerError)
        return
    }

    // Renvoyer l'URL longue au client.
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"longUrl": urlData.LongUrl})
}



// getAllUrlsHandler récupère et renvoie toutes les URLs raccourcies stockées.
func getAllUrlsHandler(w http.ResponseWriter, r *http.Request) {
    urlCollection := client.Database("urlshortener").Collection("urls")

    // Trouver toutes les entrées dans la collection d'URLs.
    cursor, err := urlCollection.Find(context.TODO(), bson.D{{}})
    if err != nil {
        http.Error(w, "Échec de la récupération des URLs", http.StatusInternalServerError)
        return
    }
    defer cursor.Close(context.TODO())

    var urls []struct {
        ShortUrl   string `bson:"shortUrl"`
        VisitCount int    `bson:"visitCount"`
    }

    // Itérer sur le curseur pour extraire les URLs.
    for cursor.Next(context.TODO()) {
        var url struct {
            ShortUrl   string `bson:"shortUrl"`
            VisitCount int    `bson:"visitCount"`
        }
        if err := cursor.Decode(&url); err != nil {
            http.Error(w, "Échec du décodage de l'URL", http.StatusInternalServerError)
            return
        }
        urls = append(urls, url)
    }

    // Vérifier s'il y a des erreurs restantes après l'itération.
    if err := cursor.Err(); err != nil {
        http.Error(w, "Erreur du curseur", http.StatusInternalServerError)
        return
    }

    // Renvoyer la liste des URLs au client.
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(urls)
}
