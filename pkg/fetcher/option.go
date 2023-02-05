package fetcher

import "net/http"

// Option configures client.
type Option func(t *client) error

// defaultOptions is a default configuration for fetcher.
var defaultOptions = []Option{
	WithHTTPClient(http.DefaultClient),
}

// WithHTTPClient returns an option that set the http client.
func WithHTTPClient(httpClient HTTPClient) Option {
	return func(c *client) error {
		if httpClient == nil {
			return errFailedSetHTTPClient
		}

		c.httpClient = httpClient

		return nil
	}
}
