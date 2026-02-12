package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/roboalchemist/tfc/pkg/jsonapi"
)

// version is injected at build time, used in User-Agent.
var version = "dev"

// Client is a Terraform Cloud API client.
type Client struct {
	httpClient *http.Client
	token      string
	baseURL    string // e.g. https://app.terraform.io/api/v2
	debug      func(string, ...interface{})
}

// NewClient creates a new Terraform Cloud API client.
func NewClient(baseURL, token string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
		token:   token,
		baseURL: baseURL + "/api/v2",
	}
}

// SetDebug sets a debug logging function.
func (c *Client) SetDebug(fn func(string, ...interface{})) {
	c.debug = fn
}

func (c *Client) debugLog(format string, args ...interface{}) {
	if c.debug != nil {
		c.debug(format, args...)
	}
}

func (c *Client) setHeaders(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/vnd.api+json")
	req.Header.Set("User-Agent", "tfc/"+version)
}

// Get performs a GET request and unmarshals the JSON:API response.
func (c *Client) Get(path string, result interface{}) error {
	return c.do("GET", path, nil, result)
}

// Post performs a POST request with a JSON body and unmarshals the response.
func (c *Client) Post(path string, body interface{}, result interface{}) error {
	return c.do("POST", path, body, result)
}

// Patch performs a PATCH request with a JSON body and unmarshals the response.
func (c *Client) Patch(path string, body interface{}, result interface{}) error {
	return c.do("PATCH", path, body, result)
}

// Delete performs a DELETE request.
func (c *Client) Delete(path string) error {
	return c.do("DELETE", path, nil, nil)
}

// GetRaw performs a GET request and returns the raw response body.
func (c *Client) GetRaw(path string) (io.ReadCloser, error) {
	url := c.baseURL + path
	c.debugLog("%s %s (raw)", "GET", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, truncate(string(body), 500))
	}

	return resp.Body, nil
}

func (c *Client) do(method, path string, body interface{}, result interface{}) error {
	url := c.baseURL + path

	var bodyReader io.Reader
	if body != nil {
		var jsonBody []byte
		var err error
		// If body is already []byte, use directly; otherwise marshal
		if b, ok := body.([]byte); ok {
			jsonBody = b
		} else {
			jsonBody, err = json.Marshal(body)
			if err != nil {
				return fmt.Errorf("marshal request: %w", err)
			}
		}
		c.debugLog("%s %s", method, url)
		c.debugLog("Request body: %s", truncate(string(jsonBody), 2000))
		bodyReader = bytes.NewBuffer(jsonBody)
	} else {
		c.debugLog("%s %s", method, url)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	c.debugLog("Response status: %d", resp.StatusCode)
	c.debugLog("Response body: %s", truncate(string(respBody), 2000))

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Try to parse JSON:API errors
		var doc jsonapi.Document
		if json.Unmarshal(respBody, &doc) == nil && len(doc.Errors) > 0 {
			return fmt.Errorf("API error (status %d): %s — %s", resp.StatusCode, doc.Errors[0].Title, doc.Errors[0].Detail)
		}
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, truncate(string(respBody), 500))
	}

	// 204 No Content — nothing to unmarshal
	if resp.StatusCode == http.StatusNoContent || len(respBody) == 0 {
		return nil
	}

	if result != nil {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("parse response: %w", err)
		}
	}

	return nil
}

// GetAllPages walks paginated results, calling collector for each page's data array.
func (c *Client) GetAllPages(path string, collector func([]jsonapi.Resource)) error {
	const maxPages = 100
	for page := 1; page <= maxPages; page++ {
		sep := "?"
		if containsQuery(path) {
			sep = "&"
		}
		url := fmt.Sprintf("%s%spage[number]=%d", path, sep, page)

		var doc jsonapi.Document
		if err := c.Get(url, &doc); err != nil {
			return err
		}

		resources, err := jsonapi.ParseList(&doc)
		if err != nil {
			return err
		}
		collector(resources)

		if doc.Meta == nil || doc.Meta.Pagination == nil || doc.Meta.Pagination.NextPage == 0 || page >= doc.Meta.Pagination.TotalPages {
			break
		}
		c.debugLog("Paginating: page %d/%d", page, doc.Meta.Pagination.TotalPages)
	}
	return nil
}

func containsQuery(path string) bool {
	for _, c := range path {
		if c == '?' {
			return true
		}
	}
	return false
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
