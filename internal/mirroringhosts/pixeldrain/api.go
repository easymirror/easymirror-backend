package pixeldrain

import (
	"encoding/base64"
	"os"
)

const (
	baseURL = "https://pixeldrain.com/api"
)

// apiKey Returns the API key from the environment variable
func apiKey() string {
	return os.Getenv("PIXELDRAIN_API_KEY")
}

// apiKeyEncoded Returns a Base64 encoded API key
func apiKeyEncoded() string {
	t := ":" + apiKey()
	return base64.StdEncoding.EncodeToString([]byte(t))
}

// authHeader returns the value for the `Authorization` header
func authHeader() string {
	return "Basic " + apiKeyEncoded()
}
