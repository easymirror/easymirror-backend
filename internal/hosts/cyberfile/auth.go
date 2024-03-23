package cyberfile

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

// getAccessToken gets and sets an access token to the account
func (a *account) GetAccessToken(ctx context.Context) (string, error) {
	// Create URL
	u, err := url.Parse(baseURI + "/authorize")
	if err != nil {
		return "", fmt.Errorf("url parse error: %w", err)
	}
	q := u.Query()
	q.Set("username", os.Getenv("CYBERFILE_USERNAME"))
	q.Set("password", os.Getenv("CYBERFILE_PASSWORD"))
	u.RawQuery = q.Encode()

	// Make request
	req, _ := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		u.String(),
		nil,
	)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error with request: %w", err)
	}

	// Parse Response
	return parseAccessToken(resp, a)
}

// parseAccessToken parses the response from the getAccessToken token
func parseAccessToken(resp *http.Response, a *account) (string, error) {
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	response := &struct {
		Data struct {
			AccessToken string `json:"access_token"`
			AccountID   string `json:"account_id"`
		} `json:"data"`
		Status   string `json:"_status"`
		Datetime string `json:"_datetime"`
		Error    string `json:"response"`
	}{}

	if err := json.Unmarshal(body, response); err != nil {
		return "", fmt.Errorf("unmarshal error: %w", err)
	}

	// Parse and validate
	if response.Status != "success" {
		return "", errors.New(response.Error)
	}
	a.accessToken = response.Data.AccessToken
	a.accountID = response.Data.AccountID
	return response.Data.AccessToken, nil
}
