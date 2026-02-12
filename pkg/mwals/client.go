package mwals

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Client is a consumer of the MW-ALS API.
type Client struct {
	BaseURL string
	HTTP    *http.Client
}

func NewClient(url string) *Client {
	return &Client{
		BaseURL: url,
		HTTP:    &http.Client{},
	}
}

// Resolve translates an alias to endpoints via the ALS server.
func (c *Client) Resolve(ctx context.Context, alias string) (*ResolutionResponse, error) {
	url := fmt.Sprintf("%s/resolve/%s", c.BaseURL, alias)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ALS resolution failed: status %d", resp.StatusCode)
	}

	var res ResolutionResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

// Register signs up a new alias.
func (c *Client) Register(ctx context.Context, req *RegistrationRequest) error {
	url := fmt.Sprintf("%s/register", c.BaseURL)
	body, _ := json.Marshal(req)

	hReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	hReq.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTP.Do(hReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("ALS registration failed: status %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) Health(ctx context.Context) (bool, error) {
	url := fmt.Sprintf("%s/health", c.BaseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, err
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}
