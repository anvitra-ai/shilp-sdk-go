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
	ID            string                 `json:"id"`
	Expiry        int64                  `json:"expiry"`
	Fields        map[string]interface{} `json:"fields"`
	KeywordFields map[string]bool        `json:"keyword_fields,omitempty"`
}

// IngestRequest represents the request to ingest data
type IngestRequest struct {
	CollectionName string   `json:"collection_name"`
	ExpiryField    string   `json:"expiry_field,omitempty"`
	Fields         []string `json:"fields,omitempty"`
	FilePath       string   `json:"file_path"`
	IDField        string   `json:"id_field,omitempty"`
	KeywordFields  []string `json:"keyword_fields,omitempty"`
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

// DebugDistanceResponse represents the response for debug distance endpoint
type DebugDistanceResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		Distance float64   `json:"distance"`
		Vector   []float64 `json:"vector"`
	} `json:"data"`
}

// DebugNeighbor represents a neighbor node in the graph
type DebugNeighbor struct {
	NodeID   int                    `json:"node_id"`
	VectorID string                 `json:"vector_id"`
	Field    string                 `json:"field"`
	Distance float64                `json:"distance"`
	Metadata map[string]interface{} `json:"metadata"`
}

// DebugNodeInfo represents detailed information about a node
type DebugNodeInfo struct {
	NodeID    int                    `json:"node_id"`
	VectorID  string                 `json:"vector_id"`
	Field     string                 `json:"field"`
	Level     int                    `json:"level"`
	Metadata  map[string]interface{} `json:"metadata"`
	Neighbors []DebugNeighbor        `json:"neighbors"`
}

// DebugNodeInfoResponse represents the response for debug node info endpoint
type DebugNodeInfoResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Data    *DebugNodeInfo `json:"data"`
}

// DebugLevelInfo represents level information
type DebugLevelInfo struct {
	Level     int `json:"level"`
	NodeCount int `json:"node_count"`
}

// DebugLevelsResponse represents the response for debug levels endpoint
type DebugLevelsResponse struct {
	Success bool                        `json:"success"`
	Message string                      `json:"message"`
	Data    map[string][]DebugLevelInfo `json:"data"`
}

// DebugNodesAtLevelResponse represents the response for debug nodes at level endpoint
type DebugNodesAtLevelResponse struct {
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Data    map[string][]int `json:"data"`
}

// DebugVectorNode represents a vector node in the reference node response
type DebugVectorNode struct {
	ID     int       `json:"id"`
	Field  string    `json:"field"`
	Vector []float64 `json:"vector"`
}

// DebugReferenceNode represents a reference node with its metadata and vector nodes
type DebugReferenceNode struct {
	ID       string                 `json:"id"`
	Metadata map[string]interface{} `json:"metadata"`
	Nodes    []DebugVectorNode      `json:"nodes"`
}

// DebugReferenceNodeResponse represents the response for debug reference node endpoint
type DebugReferenceNodeResponse struct {
	Success bool                `json:"success"`
	Message string              `json:"message"`
	Data    *DebugReferenceNode `json:"data"`
}
