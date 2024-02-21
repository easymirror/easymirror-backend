package auth

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	// How long a JWT is valid for
	accessTokenMaxAge  = 15 * time.Minute
	refreshTokenMaxAge = 90 * (24 * time.Hour) // 90 days
	issuer             = "easymirror.io"
)

// GenerateJWT is a wrapper function that generates valid access and refresh token based on the userID provided.
func GenerateJWT(userID string) (*AuthToken, error) {
	// Generate access token
	accessSecret := []byte(os.Getenv("JWT_ACCESS_SECRET"))
	accessTokenClaims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenMaxAge)),
		Subject:   userID,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    issuer,
	}
	accessTokenStr, err := generateJWT(accessSecret, accessTokenClaims)
	if err != nil {
		return nil, fmt.Errorf("SignedString error: %w", err)
	}

	// Generate refresh token
	refreshSecret := []byte(os.Getenv("JWT_REFRESH_SECRET"))
	refreshTokenClaims := &jwt.RegisteredClaims{IssuedAt: jwt.NewNumericDate(time.Now()), Issuer: issuer}
	refreshTokenStr, err := generateJWT(refreshSecret, refreshTokenClaims)
	if err != nil {
		return nil, fmt.Errorf("SignedString error: %w", err)
	}

	// Return the AuthToken
	return &AuthToken{
		AccessToken:  accessTokenStr,
		RefreshToken: refreshTokenStr,
	}, nil
}

// Generate JWT Token based in the email and in the role as input.
// Creates a token by the algorithm signing method (HS256) and the user's ID, role, and exp into claims.
// Claims are pieces of info added into the tokens.
func generateJWT(secret []byte, claims *jwt.RegisteredClaims) (string, error) {
	// The JWT library defines a struct with the MapClaims for define the different claims
	// to include in our token payload content in key-value format
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
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

// RefreshAccessToken refreshes the access token using the refresh token
func RefreshAccessToken(refreshTokenStr string) (string, error) {
	refreshToken, err := jwt.ParseWithClaims(refreshTokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_REFRESH_SECRET")), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := refreshToken.Claims.(*jwt.RegisteredClaims)
	if !ok || !refreshToken.Valid {
		return "", fmt.Errorf("invalid refresh token")
	}

	// Generate new access token
	accessSecret := []byte(os.Getenv("JWT_ACCESS_SECRET"))
	accessTokenClaims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenMaxAge)),
		Subject:   claims.Subject,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    issuer,
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	newAccessToken, err := accessToken.SignedString(accessSecret)
	if err != nil {
		return "", err
	}

	return newAccessToken, nil
}
