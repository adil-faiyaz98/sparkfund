package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sparkfund/pkg/errors"
)

// Client represents an HTTP client
type Client struct {
	client  *http.Client
	baseURL string
	headers map[string]string
}

// NewClient creates a new HTTP client
func NewClient(baseURL string) *Client {
	return &Client{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: baseURL,
		headers: make(map[string]string),
	}
}

// SetHeader sets a header for all requests
func (c *Client) SetHeader(key, value string) {
	c.headers[key] = value
}

// Get makes a GET request
func (c *Client) Get(path string, result interface{}) error {
	req, err := http.NewRequest("GET", c.baseURL+path, nil)
	if err != nil {
		return errors.ErrInternalServer(err)
	}

	for key, value := range c.headers {
		req.Header.Set(key, value)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return errors.ErrServiceUnavailable(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.NewAppError(resp.StatusCode, "Request failed", nil)
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return errors.ErrInternalServer(err)
		}
	}

	return nil
}

// Post makes a POST request
func (c *Client) Post(path string, body interface{}, result interface{}) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return errors.ErrBadRequest(err)
	}

	req, err := http.NewRequest("POST", c.baseURL+path, bytes.NewBuffer(jsonBody))
	if err != nil {
		return errors.ErrInternalServer(err)
	}

	for key, value := range c.headers {
		req.Header.Set(key, value)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return errors.ErrServiceUnavailable(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return errors.NewAppError(resp.StatusCode, "Request failed", nil)
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return errors.ErrInternalServer(err)
		}
	}

	return nil
}

// Put makes a PUT request
func (c *Client) Put(path string, body interface{}, result interface{}) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return errors.ErrBadRequest(err)
	}

	req, err := http.NewRequest("PUT", c.baseURL+path, bytes.NewBuffer(jsonBody))
	if err != nil {
		return errors.ErrInternalServer(err)
	}

	for key, value := range c.headers {
		req.Header.Set(key, value)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return errors.ErrServiceUnavailable(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.NewAppError(resp.StatusCode, "Request failed", nil)
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return errors.ErrInternalServer(err)
		}
	}

	return nil
}

// Delete makes a DELETE request
func (c *Client) Delete(path string) error {
	req, err := http.NewRequest("DELETE", c.baseURL+path, nil)
	if err != nil {
		return errors.ErrInternalServer(err)
	}

	for key, value := range c.headers {
		req.Header.Set(key, value)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return errors.ErrServiceUnavailable(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return errors.NewAppError(resp.StatusCode, "Request failed", nil)
	}

	return nil
}

// Download downloads a file
func (c *Client) Download(path string) ([]byte, error) {
	req, err := http.NewRequest("GET", c.baseURL+path, nil)
	if err != nil {
		return nil, errors.ErrInternalServer(err)
	}

	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, errors.ErrServiceUnavailable(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.NewAppError(resp.StatusCode, "Download failed", nil)
	}

	return io.ReadAll(resp.Body)
}

// Upload uploads a file
func (c *Client) Upload(path string, reader io.Reader, filename string) error {
	req, err := http.NewRequest("POST", c.baseURL+path, reader)
	if err != nil {
		return errors.ErrInternalServer(err)
	}

	for key, value := range c.headers {
		req.Header.Set(key, value)
	}
	req.Header.Set("Content-Type", "multipart/form-data")
	req.Header.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	resp, err := c.client.Do(req)
	if err != nil {
		return errors.ErrServiceUnavailable(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.NewAppError(resp.StatusCode, "Upload failed", nil)
	}

	return nil
}
