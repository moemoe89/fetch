package fetcher

//go:generate rm -f ./fetcher_mock.go
//go:generate mockgen -destination fetcher_mock.go -package fetcher -mock_names Fetcher=GoMockClient -source fetcher.go

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// errFailedSetHTTPClient represents an error message when the process of setting
// the HTTP client fails.
var errFailedSetHTTPClient = errors.New("failed to set client.http_client")

// Fetcher is an interface that defines the methods for fetching a page from a website
// and saving it to disk, as well as extracting metadata about the page.
type Fetcher interface {
	// FetchPage fetches the contents of a web page and returns it as a byte slice.
	// The url argument specifies the web page to be fetched.
	// The ctx argument is a context for fetching the web page.
	FetchPage(ctx context.Context, url string) ([]byte, error)
	// SavePage saves the contents of a web page to a file on disk.
	// The url argument specifies the web page to be saved.
	// The body argument is a byte slice containing the contents of the web page.
	SavePage(filename string, body []byte) error
	// ExtractMetadata extracts metadata about a web page, such as the number of links and images.
	ExtractMetadata(file io.Reader) ([]string, error)
	// Zip zips the given `filePaths` and `dirs` into a single archive file specified by `filename`.
	Zip(filename string, filePaths, dirs []string) error
}

// Compile time interface implementation check.
var _ Fetcher = (*client)(nil)

// HTTPClient is an interface that defines the methods of an HTTP client.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// client is a struct that implements the Fetcher and HTTPClient interfaces.
// It contains an HTTPClient field that is used to make HTTP requests.
type client struct {
	httpClient HTTPClient
}

// New returns an implementation of the Fetcher interface.
// The opts argument is a slice of options that configure the fetcher.
func New(opts ...Option) (Fetcher, error) {
	c := new(client)

	// Apply the options to the fetcher.
	for _, opt := range append(defaultOptions, opts...) {
		if err := opt(c); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	return c, nil
}
