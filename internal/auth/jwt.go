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
	accessTokenMaxAge  = 15 * time.Minute
	refreshTokenMaxAge = 90 * (24 * time.Hour) // 90 days
	issuer             = "easymirror.io"
)

// Generate JWT Token based in the email and in the role as input.
// Creates a token by the algorithm signing method (HS256) and the user's ID, role, and exp into claims.
// Claims are pieces of info added into the tokens.
func GenerateJWT(userID string) (*AuthToken, error) {
	// Add the signingkey and convert it to an array of bytes
	accessSecret := []byte(os.Getenv("JWT_ACCESS_SECRET"))
	refreshSecret := []byte(os.Getenv("JWT_REFRESH_SECRET"))

	// Generate access token
	// The JWT library defines a struct with the MapClaims for define the different claims
	// to include in our token payload content in key-value format
	accessTokenClaims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(accessTokenMaxAge).Unix(),
		Subject:   userID,
		IssuedAt:  time.Now().Unix(),
		Issuer:    issuer,
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenStr, err := accessToken.SignedString(accessSecret)
	if err != nil {
		return nil, fmt.Errorf("SignedString error: %w", err)
	}

	// Generate refresh token
	refreshTokenClaims := &jwt.StandardClaims{IssuedAt: time.Now().Unix(), Issuer: issuer}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenStr, err := refreshToken.SignedString(refreshSecret)
	if err != nil {
		return nil, fmt.Errorf("SignedString error: %w", err)
	}

	// Return the AuthToken
	return &AuthToken{
		AccessToken:  accessTokenStr,
		RefreshToken: refreshTokenStr,
	}, nil
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
		return []byte(os.Getenv("JWT_ACCESS_SECRET")), nil
	})
	if err != nil {
		log.Println("Error validating JWT:", err)
		return nil, fmt.Errorf("ValidateJWT error: %w", err)
	}
	return token, nil
}
