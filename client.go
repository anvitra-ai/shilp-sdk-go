package shilp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// Client is the client for the Shilp API
type Client struct {
	baseURL    string
	httpClient *http.Client
	authToken  string
}

// ClientOption is a function that configures the Client
type ClientOption func(*Client)

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = client
	}
}

// WithAuth sets the authentication token for the client. Use https://www.goiam.dev/api-spec#tag/auth/post/auth/v1/client to get a token. The token is a JWT that should be included in the Authorization header of each request.
func WithAuth(token string) ClientOption {
	return func(c *Client) {
		c.authToken = token
	}
}

// NewClient creates a new Shilp API client
func NewClient(baseURL string, opts ...ClientOption) *Client {
	c := &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// doRequest performs an HTTP request
func (c *Client) doRequest(method, path string, body interface{}, result interface{}, queryParams map[string]string) error {
	req, err := c.prepareRequest(method, path, body, queryParams)
	if err != nil {
		return err
	}

	if c.authToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.authToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error: %s (status: %d)", string(bodyBytes), resp.StatusCode)
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

func (c *Client) doFileRequest(method, path string, filename string) error {
	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	// Create form file field: name MUST be "file"
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return fmt.Errorf("failed to create form file: %w", err)
	}

	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	_, err = f.WriteTo(part)
	if err != nil {
		return fmt.Errorf("failed to write file to form: %w", err)
	}

	// Important: close writer to finalize the boundary
	writer.Close()

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error: %s (status: %d)", string(bodyBytes), resp.StatusCode)
	}

	return nil
}

// doRequest performs an HTTP request
func (c *Client) doRequestWithFileResponse(method, path string, body interface{}, queryParams map[string]string) (io.ReadCloser, error) {
	req, err := c.prepareRequest(method, path, body, queryParams)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s (status: %d)", string(bodyBytes), resp.StatusCode)
	}

	return resp.Body, nil
}

func (c *Client) prepareRequest(method, path string, body interface{}, queryParams map[string]string) (*http.Request, error) {
	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	if len(queryParams) > 0 {
		q := u.Query()
		for k, v := range queryParams {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
	}

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, u.String(), reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}
