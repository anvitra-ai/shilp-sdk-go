package shilp

// GenericResponse represents the standard response structure
type GenericResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// Collection represents a collection in the database
type Collection struct {
	Name             string   `json:"name"`
	IsLoaded         bool     `json:"is_loaded"`
	Fields           []string `json:"fields"`
	SearchableFields []string `json:"searchable_fields"`
}

// ListCollectionsResponse represents the response for listing collections
type ListCollectionsResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    []Collection `json:"data"`
}

// AddCollectionRequest represents the request to add a new collection
type AddCollectionRequest struct {
	Name               string `json:"name"`
	NoReferenceStorage bool   `json:"no_reference_storage,omitempty"`
}

// InsertRecordRequest represents the request to insert a record
type InsertRecordRequest struct {
	Collection string                 `json:"collection"`
	Expiry     int64                  `json:"expiry,omitempty"`
	ID         string                 `json:"id,omitempty"`
	Record     map[string]interface{} `json:"record"`
}

// InsertRecordResponse represents the response for inserting a record
type InsertRecordResponse struct {
	Success          bool        `json:"success"`
	Message          string      `json:"message"`
	Record           *RecordData `json:"record,omitempty"`
	RemainingRecords int         `json:"remaining_records,omitempty"`
}

// RecordData represents the record data in the response
type RecordData struct {
	ID     string                 `json:"id"`
	Expiry int64                  `json:"expiry"`
	Fields map[string]interface{} `json:"fields"`
}

// IngestRequest represents the request to ingest data
type IngestRequest struct {
	CollectionName string   `json:"collection_name"`
	ExpiryField    string   `json:"expiry_field,omitempty"`
	Fields         []string `json:"fields,omitempty"`
	FilePath       string   `json:"file_path"`
	IDField        string   `json:"id_field,omitempty"`
}

// IngestResponse represents the response for data ingestion
type IngestResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}

// SearchRequest represents the request body for POST search
type SearchRequest struct {
	Collection string             `json:"collection"`
	Query      string             `json:"query"`
	Fields     []string           `json:"fields,omitempty"`
	Limit      int                `json:"limit,omitempty"`
	Weights    map[string]float64 `json:"weights,omitempty"`
}

// SearchResponse represents the response for searching data
type SearchResponse struct {
	Success bool                     `json:"success"`
	Message string                   `json:"message"`
	Data    []map[string]interface{} `json:"data"`
}

// StorageItem represents an item in the storage list
type StorageItem struct {
	Name  string `json:"name"`
	IsDir bool   `json:"isDir"`
}

// ListStorageResponse represents the response for listing storage
type ListStorageResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		Items []StorageItem `json:"items"`
	} `json:"data"`
}

// ReadDocumentResponse represents the response for reading document contents
type ReadDocumentResponse struct {
	Success bool                     `json:"success"`
	Message string                   `json:"message"`
	Data    []map[string]interface{} `json:"data"`
}

// HealthResponse represents the response for health check
type HealthResponse struct {
	Success bool   `json:"success"`
	Version string `json:"version"`
}
