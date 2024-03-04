package pixeldrain

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type FolderPayload struct {
	Title     string `json:"title"`
	Anonymous bool   `json:"anonymous"`
	Files     []File `json:"files"`
}
type File struct {
	PixelDrainID string `json:"id"`
	Description  string `json:"description"`
}

// newFolder creates a list of files that can be viewed together on the file viewer page.
// It returns the URI to the PixelDrain folder
func newFolder(ctx context.Context, mirrorID string, ids []string) (string, error) {
	// Create Payload
	files := make([]File, len(ids))
	for i, id := range ids {
		files[i] = File{PixelDrainID: id, Description: fmt.Sprintf("File %v", id)}
	}
	p := &FolderPayload{
		Title:     fmt.Sprintf("Mirror %v files", mirrorID),
		Anonymous: true,
		Files:     files,
	}
	payload, _ := json.Marshal(p)

	// Make request
	req, _ := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		baseURL+"/list",
		bytes.NewBuffer(payload),
	)
	req.Header.Set("Authorization", authHeader())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request error: %w", err)
	}

	// Parse Response
	return parseFolderResponse(resp)
}

// parseFolderResponse parse's the response from the newFolder method
func parseFolderResponse(resp *http.Response) (string, error) {
	// Read the response
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	log.Println(string(body))

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
