package greenhttp

import (
	"bytes"
	"fmt"
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

// DoRequest performs an HTTP request and returns the response body as a byte slice.
func (c *HTTPClient) DoRequest(method, url string, headers map[string]string, body []byte) ([]byte, error) {
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

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

// Do performs an HTTP request and prints the response body to the console.
func (c *HTTPClient) Do(method, url string, headers map[string]string, body []byte) error {
	respBody, err := c.DoRequest(method, url, headers, body)
	if err != nil {
		return err
	}

	fmt.Println(string(respBody))
	return nil
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
