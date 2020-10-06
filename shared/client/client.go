package phoneValidation

import (
	"io"
	"kargo-back/shared/environment"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gojektech/heimdall"
	"github.com/gojektech/heimdall/httpclient"
)

var (
	httpClientTimeout           = time.Duration(environment.GetInt64("HTTP_CLIENT_TIMEOUT", 5)) * time.Second
	httpClientRetries           = int(environment.GetInt64("HTTP_CLIENT_RETRIES", 3))
	httpClientBackoffMultiplier = int(environment.GetInt64("HTTP_CLIENT_BACKOFF_MULTIPLIER", 10))
)

type Client struct {
	*httpclient.Client
}

// retries returns duration for linear backoff client interface
func retrier(retry int) time.Duration {
	if retry <= 0 {
		return 0 * time.Millisecond
	}

	return time.Duration(httpClientBackoffMultiplier*retry) * time.Millisecond
}

// NewClient returns httpclient instance with default config
func NewClient() *Client {
	return &Client{
		Client: httpclient.NewClient(
			httpclient.WithHTTPTimeout(httpClientTimeout),
			httpclient.WithRetryCount(httpClientRetries),
			httpclient.WithRetrier(heimdall.NewRetrierFunc(retrier)),
		),
	}
}

func (client *Client) PostWithURLEncodedParams(url string, params url.Values, headers http.Header) (*http.Response, error) {
	var body io.Reader

	if len(params) != 0 {
		body = strings.NewReader(params.Encode())
	}

	// Needed header for url encoded request
	headers.Set("Content-Type", "application/x-www-form-urlencoded")

	return client.Post(url, body, headers)
}
