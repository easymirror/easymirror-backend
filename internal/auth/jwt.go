package auth

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	// How long a JWT is valid for
	session_length = 24 * time.Hour
)

// Generate JWT Token based in the email and in the role as input.
// Creates a token by the algorithm signing method (HS256) and the user's ID, role, and exp into claims.
// Claims are pieces of info added into the tokens.
func GenerateJWT(userID string) (string, error) {
	// Add the signingkey and convert it to an array of bytes
	signingKey := []byte(os.Getenv("JWT_SECRET"))

	// Generate a token with the HS256 as the Signign Method
	token := jwt.New(jwt.SigningMethodHS256)

	// The JWT library defines a struct with the MapClaims for define the different claims
	// to include in our token payload content in key-value format
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(session_length).Unix()

	// Sign the token with the signingkey defined in the step before
	tokenStr, err := token.SignedString(signingKey)
	if err != nil {
		log.Println("Error during the Signing Token:", err.Error())
		return "", err
	}
	return tokenStr, err
}

// ValidateJWT validates the signature of a given JWT token
func ValidateJWT(receivedToken string) (*jwt.Token, error) {
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.
	token, err := jwt.Parse(receivedToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			// Validate the Token and return an error if the signing token is not the proper one
			return nil, fmt.Errorf("unexpected signing method: %v in token of type: %v", t.Header["alg"], t.Header["typ"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		log.Println("Error validating JWT:", err)
		return nil, fmt.Errorf("ValidateJWT error: %w", err)
	}
	return token, nil

}
