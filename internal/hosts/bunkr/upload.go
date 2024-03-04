package bunkr

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/easymirror/easymirror-backend/internal/common"
)

const (
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36"
)

// func Upload(ctx context.Context, mirrorID string, presignURIs []string) (string, error) {
// }

// getUploadLink returns a URI where files can be uploaded to
func getUploadLink(ctx context.Context) (string, error) {
	// Make request
	req, _ := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		baseURI+"/node",
		nil,
	)
	req.Header = http.Header{
		"sec-ch-ua":          {`"Chromium";v="122", "Not(A:Brand";v="24", "Brave";v="122"`},
		"accept":             {"application/json, text/plain, */*"},
		"sec-ch-ua-mobile":   {"?0"},
		"user-agent":         {userAgent},
		"token":              {os.Getenv("BUNKR_API_KEY")},
		"sec-ch-ua-platform": {`"macOS"`},
		"sec-gpc":            {"1"},
		"accept-language":    {"en-US,en;q=0.6"},
		"sec-fetch-site":     {"same-origin"},
		"sec-fetch-mode":     {"cors"},
		"sec-fetch-dest":     {"empty"},
		"referer":            {"https://app.bunkrr.su/"},
		"accept-encoding":    {"gzip, deflate, br"},
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request error: %w", err)
	}

	// Parse response
	return parseGetUploadLink(resp)
}

// parseGetUploadLink parses the response from the `getUploadLink` function.
func parseGetUploadLink(resp *http.Response) (string, error) {
	// Read Response into JSON
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	response := &struct {
		Success   bool   `json:"success"`
		ErrorDesc string `json:"description"`
		URL       string `json:"url"`
	}{}
	if err := json.Unmarshal(body, response); err != nil {
		return "", fmt.Errorf("unmarshal error: %w", err)
	}

	// Validate
	if !response.Success {
		return "", errors.New(response.ErrorDesc)
	}
	return response.URL, nil
}

