package shilp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// IngestData ingests data into a collection
func (c *Client) IngestData(req IngestRequest) (*IngestResponse, error) {
	var result IngestResponse
	err := c.doRequest(http.MethodPost, "/api/data/v1/ingest", req, &result, nil)
	return &result, err
}

// SearchDataPost searches for data in a collection using POST request.
// This method supports field-specific weights via the SearchRequest.Weights field,
// allowing fine-tuned control over search relevance scoring.
func (c *Client) SearchData(req SearchRequest) (*SearchResponse, error) {
	if req.Collection == "" {
		return nil, fmt.Errorf("collection name cannot be empty")
	}
	if len(req.VectorQuery) == 0 && len(req.Query) == 0 {
		return nil, fmt.Errorf("both vector_query and query cannot be empty")
	}
	var result SearchResponse
	err := c.doRequest(http.MethodPost, "/api/data/v1/search", req, &result, nil)
	return &result, err
}

// ListStorage lists contents of a directory in uploads storage
// if the source is mongodb, then empty path lists all DBs. If path is a DB, lists all collections in that DB.
func (c *Client) ListStorage(path string, source IngestSourceType) (*ListStorageResponse, error) {
	var result ListStorageResponse
	queryParams := map[string]string{}
	if path != "" {
		queryParams["path"] = path
	}
	err := c.doRequest(http.MethodGet, "/api/data/v1/storage/list", nil, &result, queryParams)
	return &result, err
}

func (c *Client) ListIngestSources() (*ListIngestionSourcesResponse, error) {
	var result ListIngestionSourcesResponse
	err := c.doRequest(http.MethodGet, "/api/data/v1/ingest/sources", nil, &result, nil)
	return &result, err
}

// ReadDocument reads the first few rows of a CSV document or MongoDB collection
// if the source is mongodb, then path is in the format "database/collection"
// options.query can be used to filter the documents returned incase of mongodb
func (c *Client) ReadDocument(path string, options FileReaderOptions) (*ReadDocumentResponse, error) {
	if path == "" {
		return nil, fmt.Errorf("path cannot be empty")
	}
	rows := options.Limit
	if rows < 0 {
		return nil, fmt.Errorf("rows cannot be negative")
	}
	skip := options.Skip
	if skip < 0 {
		return nil, fmt.Errorf("skip cannot be negative")
	}

	if options.Source == IngestSourceTypeMongoDB && len(strings.Split(path, "/")) != 2 {
		return nil, fmt.Errorf("for mongodb source, path must be in the format 'database/collection'")
	}

	if !options.Source.IsValid() {
		return nil, fmt.Errorf("invalid source type - %s", options.Source)
	}

	var result ReadDocumentResponse
	queryParams := map[string]string{
		"path":   path,
		"source": string(options.Source),
	}
	if rows > 0 {
		queryParams["rows"] = strconv.Itoa(rows)
	}
	if skip > 0 {
		queryParams["skip"] = strconv.Itoa(skip)
	}
	if options.Source == IngestSourceTypeMongoDB {
		filterStr, err := json.Marshal(options.MongoFilter)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal mongo filter: %w", err)
		}
		queryParams["mongo_filter"] = string(filterStr)
	}
	err := c.doRequest(http.MethodGet, "/api/data/v1/storage/read", nil, &result, queryParams)
	return &result, err
}

// UploadDataFile uploads a data file to the uploads storage which can be used for ingestion
func (c *Client) UploadDataFile(filename string) (*GenericResponse, error) {
	var result GenericResponse
	err := c.doFileRequest(http.MethodPost, "/api/data/v1/storage/upload", filename)
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

// ListEmbeddingModels lists all available embedding providers and their models
func (c *Client) ListEmbeddingModels() (*ListEmbeddingModelsResponse, error) {
	var result ListEmbeddingModelsResponse
	err := c.doRequest(http.MethodGet, "/api/data/v1/embedding/models", nil, &result, nil)
	return &result, err
}

func (c *Client) ListNLIVerticals() (*ListNLIVerticalsResponse, error) {
	var result ListNLIVerticalsResponse
	err := c.doRequest(http.MethodGet, "/api/data/v1/nli/verticals", nil, &result, nil)
	return &result, err
}
