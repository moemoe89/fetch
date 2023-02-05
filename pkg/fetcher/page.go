package fetcher

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
)

// FetchPage makes a GET request to the specified URL and returns the response body.
func (c *client) FetchPage(ctx context.Context, url string) ([]byte, error) {
	// Create a new HTTP request with the given URL.
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Make the GET request.
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	// Read all the data from the response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil
	}

	return body, nil
}

// SavePage writes the provided body to a file with the specified filename.
func (c *client) SavePage(filename string, body []byte) error {
	// Write the body to the file.
	err := os.WriteFile(filename, body, 0644)
	if err != nil {
		return err
	}

	return nil
}
