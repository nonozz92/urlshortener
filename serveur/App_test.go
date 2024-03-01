package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserSignupLoginAndUrlShortening(t *testing.T) {
	// Simuler une requête POST pour l'inscription
	user := User{Username: "testuser", Password: "testpass"}
	userData, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(userData))
	w := httptest.NewRecorder()
	registerHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %v", w.Code)
	} else {
		t.Log("User registration successful")
	}

	// Simuler une requête POST pour la connexion
	req, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(userData))
	w = httptest.NewRecorder()
	loginHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %v", w.Code)
	} else {
		t.Log("User login successful")
	}

	// Extraire le token JWT de la réponse
	var resp map[string]string
	json.NewDecoder(w.Body).Decode(&resp)
	token, exists := resp["token"]
	if !exists {
		t.Errorf("Expected JWT token in response, got none")
	} else {
		t.Logf("JWT token received: %s", token)
	}

	// Simuler une requête POST pour raccourcir une URL
	urlData := map[string]string{"longUrl": "https://www.youtube.com/watch?v=dQw4w9WgXcQ", "userToken": token}
	urlDataJson, _ := json.Marshal(urlData)
	req, _ = http.NewRequest("POST", "/shorten", bytes.NewBuffer(urlDataJson))
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	shortenUrlHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %v", w.Code)
	}

	// Récupérer et logger l'URL courte
	json.NewDecoder(w.Body).Decode(&resp)
	shortUrl, exists := resp["shortUrl"]
	if !exists {
		t.Errorf("Expected short URL in response, got none")
	} else {
		t.Logf("Short URL created: %s", shortUrl)
	}

	// Simuler une requête GET pour résoudre l'URL courte
	req, _ = http.NewRequest("GET", "/resolve?url="+shortUrl, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	resolveShortUrlHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %v", w.Code)
	} else {
		t.Log("URL resolution successful")
	}
}
