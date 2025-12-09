package shilp

import "fmt"

// ListCollections lists all collections
func (c *Client) ListCollections() (*ListCollectionsResponse, error) {
	var result ListCollectionsResponse
	err := c.doRequest("GET", "/api/collections/v1/", nil, &result, nil)
	return &result, err
}

// AddCollection adds a new collection
func (c *Client) AddCollection(req AddCollectionRequest) (*GenericResponse, error) {
	var result GenericResponse
	err := c.doRequest("POST", "/api/collections/v1/", req, &result, nil)
	return &result, err
}

// DeleteRecord deletes a record from a collection
func (c *Client) DeleteRecord(collectionName, id string) (*GenericResponse, error) {
	var result GenericResponse
	path := fmt.Sprintf("/api/collections/v1/%s/%s", collectionName, id)
	err := c.doRequest("DELETE", path, nil, &result, nil)
	return &result, err
}

// ExpiryCleanup performs expiry cleanup on a collection
func (c *Client) ExpiryCleanup(collectionName string) (*GenericResponse, error) {
	var result GenericResponse
	path := fmt.Sprintf("/api/collections/v1/%s/expiry-cleanup", collectionName)
	err := c.doRequest("POST", path, nil, &result, nil)
	return &result, err
}

// DropCollection drops an existing collection
func (c *Client) DropCollection(name string) (*GenericResponse, error) {
	var result GenericResponse
	path := fmt.Sprintf("/api/collections/v1/%s", name)
	err := c.doRequest("DELETE", path, nil, &result, nil)
	return &result, err
}

// FlushCollection flushes a collection to disk
func (c *Client) FlushCollection(name string) (*GenericResponse, error) {
	var result GenericResponse
	path := fmt.Sprintf("/api/collections/v1/%s/flush", name)
	err := c.doRequest("POST", path, nil, &result, nil)
	return &result, err
}

// LoadCollection loads a collection into memory
func (c *Client) LoadCollection(name string) (*GenericResponse, error) {
	var result GenericResponse
	path := fmt.Sprintf("/api/collections/v1/%s/load", name)
	err := c.doRequest("POST", path, nil, &result, nil)
	return &result, err
}

// UnloadCollection unloads a collection from memory
func (c *Client) UnloadCollection(name string) (*GenericResponse, error) {
	var result GenericResponse
	path := fmt.Sprintf("/api/collections/v1/%s/unload", name)
	err := c.doRequest("POST", path, nil, &result, nil)
	return &result, err
}

// RenameCollection renames an existing collection
func (c *Client) RenameCollection(oldName, newName string) (*GenericResponse, error) {
	var result GenericResponse
	path := fmt.Sprintf("/api/collections/v1/%s/rename/%s", oldName, newName)
	err := c.doRequest("PUT", path, nil, &result, nil)
	return &result, err
}

// ReIndexCollection re-indexes a collection for debug purposes
func (c *Client) ReIndexCollection(collectionName string) (*GenericResponse, error) {
	var result GenericResponse
	path := fmt.Sprintf("/api/collections/v1/%s/reindex", collectionName)
	err := c.doRequest("PUT", path, nil, &result, nil)
	return &result, err
}

// InsertRecord inserts a new record into a collection
func (c *Client) InsertRecord(req InsertRecordRequest) (*InsertRecordResponse, error) {
	var result InsertRecordResponse
	err := c.doRequest("POST", "/api/collections/v1/record", req, &result, nil)
	return &result, err
}
