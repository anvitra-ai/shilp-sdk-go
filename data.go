package shilp

import (
	"bufio"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// IngestData ingests data into a collection
func (c *Client) IngestData(req IngestRequest) (*IngestResponse, error) {
	var result IngestResponse
	err := c.doRequest("POST", "/api/data/v1/ingest", req, &result, nil)
	return &result, err
}

// SearchData searches for data in a collection using GET request.
// For advanced search with field weights and max distance, use SearchDataPost instead.
func (c *Client) SearchData(collection, query string, fields []string, limit int) (*SearchResponse, error) {
	var result SearchResponse
	queryParams := map[string]string{
		"collection": collection,
		"q":          query,
	}
	if len(fields) > 0 {
		queryParams["fields"] = strings.Join(fields, ",")
	}
	if limit > 0 {
		queryParams["limit"] = strconv.Itoa(limit)
	}

	err := c.doRequest("GET", "/api/data/v1/search", nil, &result, queryParams)
	return &result, err
}

// SearchDataPost searches for data in a collection using POST request.
// This method supports field-specific weights via the SearchRequest.Weights field,
// allowing fine-tuned control over search relevance scoring.
func (c *Client) SearchDataPost(req SearchRequest) (*SearchResponse, error) {
	var result SearchResponse
	err := c.doRequest("POST", "/api/data/v1/search", req, &result, nil)
	return &result, err
}

// ListStorage lists contents of a directory in uploads storage
func (c *Client) ListStorage(path string) (*ListStorageResponse, error) {
	var result ListStorageResponse
	queryParams := map[string]string{}
	if path != "" {
		queryParams["path"] = path
	}
	err := c.doRequest("GET", "/api/data/v1/storage/list", nil, &result, queryParams)
	return &result, err
}

// ReadDocument reads the first few rows of a CSV document
func (c *Client) ReadDocument(path string, rows, skip int) (*ReadDocumentResponse, error) {
	if path == "" {
		return nil, fmt.Errorf("path cannot be empty")
	}
	if rows < 0 {
		return nil, fmt.Errorf("rows cannot be negative")
	}
	if skip < 0 {
		return nil, fmt.Errorf("skip cannot be negative")
	}

	var result ReadDocumentResponse
	queryParams := map[string]string{
		"path": path,
	}
	if rows > 0 {
		queryParams["rows"] = strconv.Itoa(rows)
	}
	if skip > 0 {
		queryParams["skip"] = strconv.Itoa(skip)
	}
	err := c.doRequest("GET", "/api/data/v1/storage/read", nil, &result, queryParams)
	return &result, err
}

// StreamIngestStats connects to the SSE endpoint for ingestion statistics
// This returns a channel of strings (events) and an error channel.
// The caller is responsible for closing the stop channel to stop the stream.
func (c *Client) StreamIngestStats(collection string, stop <-chan struct{}) (<-chan string, <-chan error) {
	events := make(chan string)
	errs := make(chan error, 1)

	go func() {
		defer close(events)
		defer close(errs)

		url := fmt.Sprintf("%s/api/data/v1/ingest/stats?collection=%s", c.baseURL, collection)
		req, err := http.NewRequest("GET", url, nil)
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
