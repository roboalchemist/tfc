package jsonapi

import (
	"encoding/json"
	"fmt"
)

// Document wraps a JSON:API response.
type Document struct {
	Data     json.RawMessage `json:"data"`
	Included []Resource      `json:"included,omitempty"`
	Links    *Links          `json:"links,omitempty"`
	Meta     *Meta           `json:"meta,omitempty"`
	Errors   []APIError      `json:"errors,omitempty"`
}

// Resource is a single JSON:API resource object.
type Resource struct {
	ID            string                     `json:"id"`
	Type          string                     `json:"type"`
	Attributes    json.RawMessage            `json:"attributes"`
	Relationships map[string]json.RawMessage `json:"relationships,omitempty"`
	Links         map[string]interface{}     `json:"links,omitempty"`
}

// Links for pagination.
type Links struct {
	Self  string `json:"self,omitempty"`
	First string `json:"first,omitempty"`
	Last  string `json:"last,omitempty"`
	Next  string `json:"next,omitempty"`
	Prev  string `json:"prev,omitempty"`
}

// Meta for pagination info.
type Meta struct {
	Pagination *Pagination `json:"pagination,omitempty"`
}

// Pagination holds page-based pagination metadata.
type Pagination struct {
	CurrentPage int `json:"current-page"`
	PageSize    int `json:"page-size"`
	PrevPage    int `json:"prev-page"`
	NextPage    int `json:"next-page"`
	TotalPages  int `json:"total-pages"`
	TotalCount  int `json:"total-count"`
}

// APIError represents a JSON:API error object.
type APIError struct {
	Status string `json:"status,omitempty"`
	Title  string `json:"title,omitempty"`
	Detail string `json:"detail,omitempty"`
}

// ParseSingle parses a JSON:API response containing a single resource.
func ParseSingle(doc *Document) (*Resource, error) {
	if len(doc.Errors) > 0 {
		return nil, fmt.Errorf("API error: %s — %s", doc.Errors[0].Title, doc.Errors[0].Detail)
	}
	var res Resource
	if err := json.Unmarshal(doc.Data, &res); err != nil {
		return nil, fmt.Errorf("parse resource: %w", err)
	}
	return &res, nil
}

// ParseList parses a JSON:API response containing a list of resources.
func ParseList(doc *Document) ([]Resource, error) {
	if len(doc.Errors) > 0 {
		return nil, fmt.Errorf("API error: %s — %s", doc.Errors[0].Title, doc.Errors[0].Detail)
	}
	var resources []Resource
	if err := json.Unmarshal(doc.Data, &resources); err != nil {
		return nil, fmt.Errorf("parse resources: %w", err)
	}
	return resources, nil
}

// UnmarshalAttributes decodes a resource's attributes into the given struct.
func UnmarshalAttributes(res *Resource, v interface{}) error {
	return json.Unmarshal(res.Attributes, v)
}

// WrapForCreate builds a JSON:API request body for creating a resource.
func WrapForCreate(resourceType string, attrs interface{}) ([]byte, error) {
	attrBytes, err := json.Marshal(attrs)
	if err != nil {
		return nil, fmt.Errorf("marshal attributes: %w", err)
	}
	envelope := map[string]interface{}{
		"data": map[string]interface{}{
			"type":       resourceType,
			"attributes": json.RawMessage(attrBytes),
		},
	}
	return json.Marshal(envelope)
}

// WrapForUpdate builds a JSON:API request body for updating a resource.
func WrapForUpdate(id, resourceType string, attrs interface{}) ([]byte, error) {
	attrBytes, err := json.Marshal(attrs)
	if err != nil {
		return nil, fmt.Errorf("marshal attributes: %w", err)
	}
	envelope := map[string]interface{}{
		"data": map[string]interface{}{
			"id":         id,
			"type":       resourceType,
			"attributes": json.RawMessage(attrBytes),
		},
	}
	return json.Marshal(envelope)
}
