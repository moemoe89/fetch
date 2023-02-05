package fetcher

import (
	"errors"
	"io"

	"golang.org/x/net/html"
)

// targetMetadata is a map of HTML tags whose attributes (e.g., "src" or "href") should be extracted.
var targetMetadata = map[string]bool{
	"img":    true,
	"script": true,
	"link":   true,
}

var targetAttr = map[string]bool{
	"src":  true,
	"href": true,
}

// ExtractMetadata parses an HTML document from the given io.Reader
// and returns a slice of the values of "src" or "href" attributes for the HTML tags specified in targetMetadata.
func (c *client) ExtractMetadata(file io.Reader) ([]string, error) {
	var assets []string

	// Create a new HTML tokenizer.
	tokenizer := html.NewTokenizer(file)

	// Loop through the tokens in the HTML document.
	for {
		tokenType := tokenizer.Next()

		switch {
		case tokenType == html.ErrorToken:
			err := tokenizer.Err()
			if errors.Is(err, io.EOF) {
				// End of file reached.
				return assets, nil
			}

			// Return the error.
			return assets, tokenizer.Err()
		case tokenType == html.StartTagToken:
			token := tokenizer.Token()

			if targetMetadata[token.Data] {
				// This is a tag that we are interested in, extract its "src" or "href" attribute.
				for _, attr := range token.Attr {
					if targetAttr[attr.Key] {
						assets = append(assets, attr.Val)
					}
				}
			}
		}
	}
}
