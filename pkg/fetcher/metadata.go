package fetcher

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"golang.org/x/net/html"
)

// targetMetadata is a map of HTML tags whose attributes (e.g., "src" or "href") should be extracted.
var targetMetadata = map[string]bool{
	"img":    true,
	"script": true,
	"link":   true,
}

// targetAttr is a map defines a set of HTML tag attributes (keys) that the function ExtractMetadata is interested in.
var targetAttr = map[string]bool{
	"src":  true,
	"href": true,
}

// Metadata data structure for the metadata web page.
type Metadata struct {
	Site      string    `json:"site"`
	NumLinks  int64     `json:"num_links"`
	Images    int64     `json:"images"`
	Assets    []string  `json:"assets"`
	LastFetch time.Time `json:"last_fetch"`
}

// ExtractMetadata parses an HTML document from the given io.Reader
// and returns a slice of the values of "src" or "href" attributes for the HTML tags specified in targetMetadata.
func (c *client) ExtractMetadata(url, filePath string, file io.Reader) (*Metadata, error) {
	metadata := &Metadata{
		Site: url,
	}

	// TODO:
	// Need to research about this code.
	// Alternative solution to fetch more assets.
	/*
		doc, err := html.Parse(file)
		if err != nil {
			return nil, fmt.Errorf("failed to parse html: %w", err)
		}

		var links []string
		var f func(*html.Node)
		f = func(n *html.Node) {
			if n.Type == html.ElementNode && targetMetadata[n.Data] {
				for _, a := range n.Attr {
					if strings.Contains(a.Val, "base64") {
						continue
					}

					if (n.Data == "img" || n.Data == "script") && a.Key == "src" {
						links = append(links, a.Val)
					}

					if n.Data == "link" && a.Key == "href" {
						links = append(links, a.Val)
					}
				}
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
		f(doc)

		metadata.Assets = links
	*/

	// Create a new HTML tokenizer.
	tokenizer := html.NewTokenizer(file)

	// Loop through the tokens in the HTML document.
	for {
		tokenType := tokenizer.Next()

		switch {
		case tokenType == html.ErrorToken:
			err := tokenizer.Err()

			// If end of file reached.
			if errors.Is(err, io.EOF) {
				// Save metadata to JSON file.
				err = c.saveMetadataJSON(metadata, filePath)
				if err != nil {
					return nil, err
				}

				return metadata, nil
			}

			// Return the error.
			return nil, fmt.Errorf("failed to extract metadata: %w", tokenizer.Err())
		case tokenType == html.StartTagToken:
			token := tokenizer.Token()

			// Skips if the token data is not the target metadata.
			if !targetMetadata[token.Data] {
				continue
			}

			// Count assets like link and images.
			c.countAssets(metadata, token.Data)

			// This is a tag that we are interested in, extract its "src" or "href" attribute.
			for _, attr := range token.Attr {
				// Skips if the attr key is not the target attribute.
				if !targetAttr[attr.Key] {
					continue
				}

				// Collects the assets.
				metadata.Assets = append(metadata.Assets, attr.Val)
			}
		}
	}
}

// countAssets counts the number of link and image asset.
func (c *client) countAssets(metadata *Metadata, tokenData string) {
	if tokenData == "img" {
		metadata.Images++
	}

	if tokenData == "script" || tokenData == "link" {
		metadata.NumLinks++
	}
}

// saveMetadataJSON saves the metadata to JSON file.
func (c *client) saveMetadataJSON(metadata *Metadata, filePath string) error {
	var lastFetch time.Time

	// Checks if the metadata already exists.
	// If exists, gets the last fetch time
	metadataFile, err := os.ReadFile(filePath)
	if err == nil {
		var lastMetadata *Metadata

		err = json.Unmarshal(metadataFile, &lastMetadata)
		if err != nil {
			return fmt.Errorf("failed to unmarshal metadata: %w", err)
		}

		lastFetch = lastMetadata.LastFetch
	}

	// Updates last fetch.
	metadata.LastFetch = time.Now().UTC()

	metadataJSON, err := json.MarshalIndent(metadata, "", " ")
	if err != nil {
		return fmt.Errorf("failed to marshal indent metadata: %w", err)
	}

	// Save the new JSON.
	err = c.SavePage(filePath, metadataJSON)
	if err != nil {
		return fmt.Errorf("failed to save metadata json: %w", err)
	}

	// Update metadata last fetch.
	metadata.LastFetch = lastFetch

	return nil
}

// StringMetadata build the string message of metadata.
func (c *client) StringMetadata(metadata *Metadata) string {
	lastFetch := "-"

	// If last fetch isn't empty, change to `Tue Mar 16 2021 15:46 UTC` format.
	if !metadata.LastFetch.IsZero() {
		lastFetch = metadata.LastFetch.Format("Mon Jan 02 2006 15:04 MST")
	}

	return fmt.Sprintf("site: %s\nnum_links: %d\nimages: %d\nlast_fetch: %s\n\n",
		metadata.Site,
		metadata.NumLinks,
		metadata.Images,
		lastFetch,
	)
}
