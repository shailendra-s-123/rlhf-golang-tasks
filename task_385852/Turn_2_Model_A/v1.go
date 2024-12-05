package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// User represents a user in our system.
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var users = map[string]User{
	"admin": {"admin", "password"},
}

// GenerateToken generates a JWT token for a given user.
func GenerateToken(user User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username
	claims["exp"] = int64(time.Now().Add(time.Hour).Unix())

	secret := []byte("your-secret-key") // Replace this with a secure secret key
	signedToken, err := token.SignedString(secret)
	return signedToken, err
}

// AuthenticateHandler handles authentication requests.
func AuthenticateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var user User
	if err := decoder.Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, ok := users[user.Username]; !ok || users[user.Username].Password != user.Password {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	token, err := GenerateToken(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func main() {
	http.HandleFunc("/auth", AuthenticateHandler)
	fmt.Println("Server listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
