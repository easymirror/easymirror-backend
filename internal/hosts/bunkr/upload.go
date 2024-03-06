package bunkr

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/easymirror/easymirror-backend/internal/common"
)

const (
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36"
)

// Upload upload's files to a folder on Bunkr.
// If successful, the URI of the folder is returned
func Upload(ctx context.Context, mirrorID string, sourceURIs []string) (string, error) {
	if len(sourceURIs) == 0 {
		return "", errors.New("no source uri")
	}

	// Create Folder
	folderID, err := createFolder(
		ctx,
		fmt.Sprintf("Mirror %v files", mirrorID),
		true,
		true,
	)
	if err != nil {
		return "", fmt.Errorf("create folder error: %w", err)
	}

	// Upload to folder
	for _, uri := range sourceURIs {
		if _, err := upload(ctx, folderID, uri); err != nil {
			log.Println("Error uploading file:", err)
		}
	}

	// Return URI of the folder
	folder, err := getFolder(ctx, folderID)
	if err != nil {
		return "", fmt.Errorf("error getting folder")
	}
	return folder.Link, nil
}

// UploadTx upload's files to a folder on Bunkr with a given TX. Adds an entry to the database but does not committ the TX.
// If successful, the URI of the folder is returned
func UploadTx(ctx context.Context, tx *sql.Tx, mirrorID string, sourceURIs []string) (string, error) {
	if len(sourceURIs) == 0 {
		return "", errors.New("no source uri")
	}

	// Create Folder
	folderID, err := createFolder(
		ctx,
		fmt.Sprintf("Mirror %v files", mirrorID),
		true,
		true,
	)
	if err != nil {
		return "", fmt.Errorf("create folder error: %w", err)
	}

	// Upload to folder
	for _, uri := range sourceURIs {
		if _, err := upload(ctx, folderID, uri); err != nil {
			log.Println("Error uploading file:", err)
		}
	}

	// Return URI of the folder
	folder, err := getFolder(ctx, folderID)
	if err != nil {
		return "", fmt.Errorf("error getting folder")
	}

	// Add to `host_links` table
	// This statement allows you to insert a new row into a table,
	// or update an existing row if a conflict (e.g., duplicate key violation) occurs.
	// This is known as UPSERT
	statement := `
		INSERT INTO host_links (mirror_id, bunkr)
		VALUES (($1), ($2))
		ON CONFLICT (mirror_id) 
		DO UPDATE
		SET bunkr = EXCLUDED.bunkr;
	`
	if _, err = tx.Exec(statement, mirrorID, folder.Link); err != nil {
		return "", fmt.Errorf("exec tx error: %w", err)
	}
	return folder.Link, nil
}

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

func upload(ctx context.Context, albumID string, presignURI string) (string, error) {
	// Get upload link
	uploadLink, err := getUploadLink(ctx)
	if err != nil {
		return "", fmt.Errorf("getUploadLink error: %w", err)
	}

	// Get the file from the presigned URL.
	req, _ := http.NewRequestWithContext(ctx, "GET", presignURI, nil)
	resp1, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error getting body from presigned URL: %w", err)
	}
	defer resp1.Body.Close()

	// We use an io.Pipe and a goroutine for writing from the file/response body
	// and reading to the request concurrently.
	// By doing this, we don't have to load the entire file into memory/a buffer.
	// This also sends the file as a "transfer encoding chunked"
	r, w := io.Pipe()
	m := multipart.NewWriter(w)
	go func() {
		defer w.Close()
		defer m.Close()

		// Get the filename form the URL
		name, err := common.FilenameFromURI(presignURI)
		if err != nil {
			return
		}

		part, err := m.CreateFormFile("files[]", name)
		if err != nil {
			return
		}
		if _, err = io.Copy(part, resp1.Body); err != nil {
			return
		}
	}()

	// Make request
	req, _ = http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		uploadLink,
		r,
	)
	req.Header = http.Header{
		"sec-ch-ua":          {`"Chromium";v="122", "Not(A:Brand";v="24", "Brave";v="122"`},
		"albumid":            {albumID},
		"sec-ch-ua-mobile":   {"?0"},
		"user-agent":         {userAgent},
		"Content-Type":       {m.FormDataContentType()},
		"accept":             {"application/json"},
		"cache-control":      {"no-cache"},
		"x-requested-with":   {"XMLHttpRequest"},
		"token":              {os.Getenv("BUNKR_API_KEY")},
		"sec-ch-ua-platform": {`"macOS"`},
		"sec-gpc":            {"1"},
		"accept-language":    {"en-US,en;q=0.6"},
		"origin":             {"https://app.bunkrr.su"},
		"sec-fetch-site":     {"cross-site"},
		"sec-fetch-mode":     {"cors"},
		"sec-fetch-dest":     {"empty"},
		"referer":            {"https://app.bunkrr.su/"},
		// "accept-encoding":    {"gzip, deflate, br"},
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request error: %w", err)
	}

	// Parse Response
	return parseUpload(resp)
}

// parseUpload parses the response from the `upload` function
func parseUpload(resp *http.Response) (string, error) {
	// Convert into JSON
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	response := &struct {
		ErrorDescription string `json:"description"`
		Success          bool   `json:"success"`
		Files            []struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"files"`
	}{}
	if err := json.Unmarshal(body, response); err != nil {
		return "", fmt.Errorf("unmarshal error: %w", err)
	}

	// Validate
	if !response.Success {
		return "", errors.New(response.ErrorDescription)
	}

	// Because we are only updating 1 file at a time, we take the first link
	for _, file := range response.Files {
		return file.URL, nil
	}

	// If it gets to this point, means there was no files, which is an error
	return "", errors.New("no file links in response")
}
