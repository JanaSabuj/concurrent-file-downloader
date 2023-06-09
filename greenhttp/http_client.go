package greenhttp

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

type HTTPClient struct {
	client *http.Client
}

// NewHTTPClient creates a new instance of HTTPClient.
func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		client: &http.Client{},
	}
}

// DoRequest performs an HTTP request using the provided http.Request object and returns the response.
func (c *HTTPClient) DoRequest(req *http.Request) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// NewRequest creates a new HTTP request with the specified method, URL, headers, and body.
func (c *HTTPClient) NewRequest(method, url string, headers map[string]string, body []byte) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	if body != nil {
		req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	}

	return req, nil
}
