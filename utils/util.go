package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/golang-jwt/jwt"
)

var secretKey = []byte("secret-key")

type User struct {
	Username string
	Password string
}

func isDomainAllowed(domain string, allowedDomains []string) bool {
	for _, allowedDomain := range allowedDomains {
		if domain == allowedDomain {
			return true
		}
	}
	return false
}

func handleFetch(w http.ResponseWriter, r *http.Request) {
	rawURL := r.FormValue("url")

	// Parse the URL to ensure it is valid and safe
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	// Ensure that the URL points to a permitted domain
	allowedDomains := []string{"example.com", "trusteddomain.com"}
	if !isDomainAllowed(parsedURL.Host, allowedDomains) {
		http.Error(w, "Access to the specified domain is not allowed", http.StatusForbidden)
		return
	}

	// Make a request to the provided URL
	response, err := http.Get(parsedURL.String())
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(w, "Response Body: %s", body)
}

func createToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func verifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var u User
	json.NewDecoder(r.Body).Decode(&u)
	fmt.Printf("The user request value %v", u)

	if u.Username == "Chek" && u.Password == "123456" {
		tokenString, err := createToken(u.Username)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Errorf("No username found")
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, tokenString)
		return
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Invalid credentials")
	}
}

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Missing authorization header")
		return
	}

	tokenString = tokenString[len("Bearer "):]

	err := verifyToken(tokenString)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Invalid token")
		return
	}

	fmt.Fprint(w, "Welcome to the the protected area")

}
