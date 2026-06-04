package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/arunhiremath92/GoSecureToken/hmac"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var userStore = map[string]string{
	"alice": "pass@123",
	"bob":   "pass@1234",
	"arun":  "xpertscan",
}

func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedBytes), err
}

func checkPasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate presence of username and password
	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Authenticate user...
	value, exists := userStore[req.Username]
	if !exists {
		errorMsg := fmt.Sprintf("no user exists: %v", req.Username)
		http.Error(w, errorMsg, http.StatusBadRequest)
		log.Println(errorMsg)
	}

	if !checkPasswordHash(req.Password, value) {
		errorMsg := fmt.Sprintf("failed to validate the user: %v", req.Username)
		http.Error(w, errorMsg, http.StatusBadRequest)
		log.Println(errorMsg)
	}

	log.Println("Received request to generate token")
	log.Printf("Extracted username from URL: %s\n", req.Username)

	if req.Username == "" {
		http.Error(w, "Username parameter is missing in the URL", http.StatusBadRequest)
		log.Println("Username parameter missing in request")
		return
	}
	token := hmac.GetInstance().GenerateToken(req.Username)
	log.Printf("Generated token for user '%s': %s\n", req.Username, token)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(token))
}

// TokenVerifiyHandler handles token verification requests.
// It expects a username and token in the URL and verifies if the token is valid for the given username.
func TokenVerifiyHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to verify token")
	vars := mux.Vars(r)
	username := vars["username"]
	token := vars["token"]
	log.Printf("Extracted username: %s, token: %s\n", username, token)

	if username == "" || token == "" {
		http.Error(w, "username or token parameter is missing in the URL", http.StatusBadRequest)
		log.Println("username or token parameter missing in request")
		return
	}

	tokenUserName, err := hmac.GetInstance().VerifyToken(token)
	if err != nil {
		errorMsg := fmt.Sprintf("Failed to verify the token: %v", err)
		http.Error(w, errorMsg, http.StatusBadRequest)
		log.Println(errorMsg)
		return
	}

	if tokenUserName == username {
		log.Printf("Token verification successful for user '%s'\n", username)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
		return
	}

	errorMsg := fmt.Sprintf("User name does not match with the token: expected '%s', got '%s'", username, tokenUserName)
	http.Error(w, errorMsg, http.StatusBadRequest)
	log.Println(errorMsg)
}

// main sets up the HTTP server and routes.
func main() {
	log.Println("Starting GoSecureToken API server on port 8080")
	for user, password := range userStore {
		if hashed, error := hashPassword(password); error == nil {
			userStore[user] = hashed
		}
	}
	r := mux.NewRouter()
	r.HandleFunc("/api/login", LoginHandler).Methods("POST")
	r.HandleFunc("/api/protect}", TokenVerifiyHandler).Methods("GET")
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
