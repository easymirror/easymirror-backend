package cyberfile

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type folderResponse struct {
	Data     folderData `json:"data"`
	Status   string     `json:"_status"`
	Datetime string     `json:"_datetime"`
	Error    string     `json:"response"`
}
type folderData struct {
	ID               string `json:"id"`
	ParentID         string `json:"parentId"`
	FolderName       string `json:"folderName"`
	TotalSize        string `json:"totalSize"`
	IsPublic         string `json:"isPublic"`
	AccessPassword   string `json:"accessPassword"`
	DateAdded        string `json:"date_added"`
	DateUpdated      string `json:"date_updated"`
	URLFolder        string `json:"url_folder"`
	TotalDownloads   int    `json:"total_downloads"`
	ChildFolderCount int    `json:"child_folder_count"`
	FileCount        int    `json:"file_count"`
}

// createFolder creates a new folder in Cyberfile's API.
// If successful, it returns the ID of the folder.
func createFolder(ctx context.Context, a Account, mirrorID string) (string, error) {
	// Create URL
	u, err := url.Parse(baseURI + "/folder/create")
	if err != nil {
		return "", fmt.Errorf("url parse error: %w", err)
	}
	q := u.Query()
	q.Set("access_token", a.AccessToken())
	q.Set("account_id", a.AccountID())
	q.Set("folder_name", fmt.Sprintf("Mirror %v files", mirrorID))
	q.Set("is_public", "1") // Whether a folder is available publicly or private only. 0 = Private, 1 = Unlisted, 2 = Public in site search.
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

	// Parse response
	return parseCreateFolder(resp)
}

// createFolder parses the response from the createFolder function
func parseCreateFolder(resp *http.Response) (string, error) {
	// Read the body into a JSON struct
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	response := &folderResponse{}
	if err := json.Unmarshal(body, response); err != nil {
		return "", fmt.Errorf("umarshal error: %w", err)
	}

	// Parse & Validate response
	if response.Status != "success" {
		return "", errors.New(response.Error)
	}

	return response.Data.ID, nil
}
