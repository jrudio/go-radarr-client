package radarr

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// utils.go holds network utils and function helpers

func (c Client) get(query string, params url.Values) (*http.Response, error) {
	relativeURL, err := url.Parse(query)

	if err != nil {
		return &http.Response{}, err
	}

	endpointURL := c.URL.ResolveReference(relativeURL)

	if params == nil {
		params = endpointURL.Query()
	}

	endpointURL.RawQuery = params.Encode()

	client := http.Client{
		Timeout: time.Duration(c.Timeout) * time.Second,
	}

	fmt.Printf("radarr GET request: %s\n", endpointURL)

	req, err := http.NewRequest("GET", endpointURL.String(), nil)

	if err != nil {
		return &http.Response{}, err
	}

	req.Header.Set("x-api-key", c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	return client.Do(req)
}

func (c Client) post(query string, body []byte) (*http.Response, error) {
	relativeURL, err := url.Parse(query)

	if err != nil {
		return &http.Response{}, err
	}

	endpointURL := c.URL.ResolveReference(relativeURL)

	client := http.Client{
		Timeout: time.Duration(c.Timeout) * time.Second,
	}

	req, err := http.NewRequest("POST", endpointURL.String(), bytes.NewBuffer(body))

	if err != nil {
		return &http.Response{}, err
	}

	req.Header.Set("x-api-key", c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	return client.Do(req)
}

func (c Client) delete(query string, params url.Values) (*http.Response, error) {
	relativeURL, err := url.Parse(query)

	if err != nil {
		return &http.Response{}, err
	}

	endpointURL := c.URL.ResolveReference(relativeURL)

	if params == nil {
		params = endpointURL.Query()
	}

	endpointURL.RawQuery = params.Encode()

	client := http.Client{
		Timeout: time.Duration(c.Timeout) * time.Second,
	}

	req, err := http.NewRequest("DELETE", endpointURL.String(), nil)

	if err != nil {
		return &http.Response{}, err
	}

	req.Header.Set("x-api-key", c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	return client.Do(req)
}

func encodeURL(str string) (string, error) {
	u, err := url.Parse(str)

	if err != nil {
		return "", err
	}

	return u.String(), nil
}
