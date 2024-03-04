package bunkr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

// createFolder creates a new folder on Bunkr's API.
// If successful, it returns the ID of the album
func createFolder(ctx context.Context, name string, downloadable, public bool) (string, error) {
	// Create Payload
	p := &struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Download    bool   `json:"download"`
		Public      bool   `json:"public"`
	}{Name: name, Download: downloadable, Public: public}
	payload, _ := json.Marshal(p)

	// Make request
	req, _ := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		baseURI+"/albums",
		bytes.NewBuffer(payload),
	)
	req.Header = http.Header{
		"sec-ch-ua":          {`"Chromium";v="122", "Not(A:Brand";v="24", "Brave";v="122"`},
		"accept":             {"application/json, text/plain, */*"},
		"content-type":       {"application/json;charset=UTF-8"},
		"sec-ch-ua-mobile":   {"?0"},
		"user-agent":         {userAgent},
		"token":              {os.Getenv("BUNKR_API_KEY")},
		"sec-ch-ua-platform": {`"macOS"`},
		"sec-gpc":            {"1"},
		"accept-language":    {"en-US,en;q=0.6"},
		"origin":             {"https://app.bunkrr.su"},
		"sec-fetch-site":     {"same-origin"},
		"sec-fetch-mode":     {"cors"},
		"sec-fetch-dest":     {"empty"},
		"referer":            {"https://app.bunkrr.su/"},
		"accept-encoding":    {"gzip, deflate, br"},
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error with request: %w", err)
	}

	// Parse Response
	return parseCreateFolder(resp)
}

func parseCreateFolder(resp *http.Response) (string, error) {
	// Convert to JSON
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	response := &struct {
		Success          bool   `json:"success"`
		ErrorDescription string `json:"description"`
		ID               int    `json:"id"`
	}{}
	if err := json.Unmarshal(body, response); err != nil {
		return "", fmt.Errorf("unmarshal error: %w", err)
	}

	// Validate
	if !response.Success {
		return "", fmt.Errorf("error with API: %q", response.ErrorDescription)
	}
	return strconv.Itoa(response.ID), nil
}
