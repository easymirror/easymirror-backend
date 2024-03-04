package pixeldrain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/easymirror/easymirror-backend/internal/common"
)

// Upload is a wrapper function to upload to PixelDrain's API.
// If successful, it returns a link to the folder with the uploaded files
func Upload(ctx context.Context, mirrorID string, presignURIs []string) (string, error) {
	if len(presignURIs) < 1 {
		return "", errors.New("no presigned URLs")
	}

	// Upload the files to PixelDrain's API
	ids := []string{}
	for _, uri := range presignURIs {
		fileID, err := upload(ctx, uri)
		if err != nil {
			log.Println("Error uploading file:", err)
			continue
		}
		ids = append(ids, fileID)
	}

	// Create a new folder
	folderID, err := newFolder(ctx, mirrorID, ids)
	if err != nil {
		return "", fmt.Errorf("newFolder error: %w", err)
	}
	return folderBaseURL + "/" + folderID, nil
}

// upload upload's a given file to PixelDrain's API.
// If the body is successfully uploaded, the ID of the upload is returned.
func upload(ctx context.Context, presignURI string) (string, error) {
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

		part, err := m.CreateFormFile("file", name)
		if err != nil {
			return
		}
		if _, err = io.Copy(part, resp1.Body); err != nil {
			return
		}
	}()

	// Upload the file to PixelDrain
	req, _ = http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		baseURL+"/file",
		r,
	)
	req.Header = http.Header{
		"Content-Type":  {m.FormDataContentType()},
		"Authorization": {authHeader()},
	}

	client := &http.Client{}
	resp2, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error uploading to pixeldrain: %w", err)
	}
	defer resp2.Body.Close()
	return parseUpload(resp2)
}

// parseUpload parses the response from the upload function
func parseUpload(r *http.Response) (string, error) {
	// Read the response
	defer r.Body.Close()
	body, _ := io.ReadAll(r.Body)

	// Parse the response
	response := &struct {
		Success bool   `json:"success"`
		ID      string `json:"id"`
		ErrorID string `json:"value"`
		Message string `json:"message"`
	}{}
	var err error
	if err = json.Unmarshal(body, response); err != nil {
		return "", fmt.Errorf("unmarshal error: %w", err)
	}
	if !response.Success {
		return "", fmt.Errorf("upload error: %v - %v", response.ErrorID, response.Message)
	}
	return response.ID, nil
}
