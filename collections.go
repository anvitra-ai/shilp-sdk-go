package shilp

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
)

// ListCollections lists all collections
func (c *Client) ListCollections() (*ListCollectionsResponse, error) {
	var result ListCollectionsResponse
	err := c.doRequest(http.MethodGet, "/api/collections/v1/", nil, &result, nil)
	return &result, err
}

// AddCollection adds a new collection
func (c *Client) AddCollection(req AddCollectionRequest) (*GenericResponse, error) {
	var result GenericResponse
	err := c.doRequest(http.MethodPost, "/api/collections/v1/", req, &result, nil)
	return &result, err
}

// DeleteRecord deletes a record from a collection
func (c *Client) DeleteRecord(collectionName, id string) (*GenericResponse, error) {
	var result GenericResponse
	path := fmt.Sprintf("/api/collections/v1/%s/%s", collectionName, id)
	err := c.doRequest(http.MethodDelete, path, nil, &result, nil)
	return &result, err
}

// ExpiryCleanup performs expiry cleanup on a collection
func (c *Client) ExpiryCleanup(collectionName string) (*GenericResponse, error) {
	var result GenericResponse
	path := fmt.Sprintf("/api/collections/v1/%s/expiry-cleanup", collectionName)
	err := c.doRequest(http.MethodPost, path, nil, &result, nil)
	return &result, err
}

// DropCollection drops an existing collection
func (c *Client) DropCollection(name string) (*GenericResponse, error) {
	var result GenericResponse
	path := fmt.Sprintf("/api/collections/v1/%s", name)
	err := c.doRequest(http.MethodDelete, path, nil, &result, nil)
	return &result, err
}

// FlushCollection flushes a collection to disk
func (c *Client) FlushCollection(name string) (*GenericResponse, error) {
	var result GenericResponse
	path := fmt.Sprintf("/api/collections/v1/%s/flush", name)
	err := c.doRequest(http.MethodPost, path, nil, &result, nil)
	return &result, err
}

// LoadCollection loads a collection into memory
func (c *Client) LoadCollection(name string) (*GenericResponse, error) {
	var result GenericResponse
	path := fmt.Sprintf("/api/collections/v1/%s/load", name)
	err := c.doRequest(http.MethodPost, path, nil, &result, nil)
	return &result, err
}

// UnloadCollection unloads a collection from memory
func (c *Client) UnloadCollection(name string) (*GenericResponse, error) {
	var result GenericResponse
	path := fmt.Sprintf("/api/collections/v1/%s/unload", name)
	err := c.doRequest(http.MethodPost, path, nil, &result, nil)
	return &result, err
}

// ExportCollection exports a collection and returns a ReadCloser for the file
// The caller can io.Copy to store the file locally and is responsible for closing the ReadCloser
func (c *Client) ExportCollection(name string) (io.ReadCloser, error) {
	path := fmt.Sprintf("/api/collections/v1/%s/export", name)
	return c.doRequestWithFileResponse(http.MethodPost, path, nil, nil)
}

// ImportCollection imports a collection from a file
// The filename parameter should be the exported collection file
func (c *Client) ImportCollection(filename string) (*GenericResponse, error) {
	var result GenericResponse
	err := c.doFileRequest(http.MethodPost, "/api/collections/v1/import", filename)
	return &result, err
}

// RenameCollection renames an existing collection
func (c *Client) RenameCollection(oldName, newName string) (*GenericResponse, error) {
	var result GenericResponse
	path := fmt.Sprintf("/api/collections/v1/%s/rename/%s", oldName, newName)
	err := c.doRequest(http.MethodPut, path, nil, &result, nil)
	return &result, err
}

// ReIndexCollection re-indexes a collection for debug purposes
func (c *Client) ReIndexCollection(collectionName string) (*GenericResponse, error) {
	var result GenericResponse
	path := fmt.Sprintf("/api/collections/v1/%s/reindex", collectionName)
	err := c.doRequest(http.MethodPut, path, nil, &result, nil)
	return &result, err
}

// PQTrain performs Product Quantization training for an existing collection
func (c *Client) PQTrain(collectionName string) (*GenericResponse, error) {
	var result GenericResponse
	path := fmt.Sprintf("/api/collections/v1/%s/pq-train", collectionName)
	err := c.doRequest(http.MethodPost, path, nil, &result, nil)
	return &result, err
}

// InsertRecord inserts a new record into a collection
func (c *Client) InsertRecord(req InsertRecordRequest) (*InsertRecordResponse, error) {
	var result InsertRecordResponse
	err := c.doRequest(http.MethodPost, "/api/collections/v1/record", req, &result, nil)
	return &result, err
}

func (c *Client) GetCollectionData(collectionName string, offset, limit int) (*GetCollectionDataResponse, error) {
	var result GetCollectionDataResponse
	path := fmt.Sprintf("/api/collections/v1/%s/data?offset=%d&limit=%d", collectionName, offset, limit)
	err := c.doRequest(http.MethodGet, path, nil, &result, nil)
	return &result, err
}

// EnableNLI enables Natural Language Inference for a collection and vertical. This is an SSE endpoint that streams the progress of enabling NLI.
// The caller is responsible for closing the stop channel to stop the stream.
// vertical should be the vertical supported by the NLI provider, if this need to be a custom vertical, keep it empty.
func (c *Client) EnableNLI(collection string, vertical string, stop <-chan struct{}) (<-chan string, <-chan error) {
	events := make(chan string)
	errs := make(chan error, 1)

	go func() {
		defer close(events)
		defer close(errs)

		url := fmt.Sprintf("%s/api/collections/v1/%s/nli/enable?vertical=%s", c.baseURL, collection, vertical)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			errs <- err
			return
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			errs <- err
			return
		}
		defer resp.Body.Close()

		reader := bufio.NewReader(resp.Body)

		for {
			select {
			case <-stop:
				return
			default:
				line, err := reader.ReadString('\n')
				if err != nil {
					errs <- err
					return
				}
				events <- line
			}
		}
	}()

	return events, errs
}

func (c *Client) GetCollectionSchema(collectionName string) (*GetCollectionSchemaResponse, error) {
	var result GetCollectionSchemaResponse
	path := fmt.Sprintf("/api/collections/v1/%s/schema", collectionName)
	err := c.doRequest(http.MethodGet, path, nil, &result, nil)
	return &result, err
}

func (c *Client) EnableMetadataStore(collectionName string, req EnableMetadataStoreRequest) (*EnableMetadataStoreResponse, error) {
	var result EnableMetadataStoreResponse
	path := fmt.Sprintf("/api/collections/v1/%s/metadata/enable", collectionName)
	err := c.doRequest(http.MethodPost, path, req, &result, nil)
	return &result, err
}
