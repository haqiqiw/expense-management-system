package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

//go:generate mockery --name=APIClient --structname APIClient --outpkg=mocks --output=./../mocks
type APIClient interface {
	Get(ctx context.Context, path string) (*APIResponse, error)
	Post(ctx context.Context, path string, body interface{}) (*APIResponse, error)
	Put(ctx context.Context, path string, body interface{}) (*APIResponse, error)
	Delete(ctx context.Context, path string) (*APIResponse, error)
}

type APIResponse struct {
	StatusCode int
	Headers    http.Header
	Body       json.RawMessage
}

type Client struct {
	BaseURL string
	Timeout time.Duration
	client  *http.Client
}

func NewClient(baseURL string, timeout time.Duration) *Client {
	return &Client{
		BaseURL: baseURL,
		Timeout: timeout,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) Get(ctx context.Context, path string) (*APIResponse, error) {
	return c.Do(ctx, http.MethodGet, path, nil)
}

func (c *Client) Post(ctx context.Context, path string, body interface{}) (*APIResponse, error) {
	return c.Do(ctx, http.MethodPost, path, body)
}

func (c *Client) Put(ctx context.Context, path string, body interface{}) (*APIResponse, error) {
	return c.Do(ctx, http.MethodPatch, path, body)
}

func (c *Client) Delete(ctx context.Context, path string) (*APIResponse, error) {
	return c.Do(ctx, http.MethodDelete, path, nil)
}

func (c *Client) Do(ctx context.Context, method, path string, body interface{}) (*APIResponse, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body = %w", err)
		}
		bodyReader = bytes.NewReader(jsonBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.BaseURL+path, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request = %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request faile = %w", err)
	}
	defer resp.Body.Close()

	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response = %w", err)
	}

	return &APIResponse{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       rawBody,
	}, nil
}
