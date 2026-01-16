package shilp

import (
	"fmt"
	"time"
)

// GenericResponse represents the standard response structure
type GenericResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// MetadataField represents a metadata field definition
type MetadataField struct {
	Name string `json:"name"`
	Type int    `json:"type"`
}

// Collection represents a collection in the database
type Collection struct {
	Name               string                 `json:"name"`
	IsLoaded           bool                   `json:"is_loaded"`
	Fields             []string               `json:"fields"`
	SearchableFields   []string               `json:"searchable_fields"`
	Metadata           []MetadataColumnSchema `json:"metadata,omitempty"`
	HasMetadataEnabled bool                   `json:"has_metadata_enabled"`
	NoReferenceStorage bool                   `json:"no_reference_storage"`
}

type MetadataColumnSchema struct {
	Name string   `json:"name"`
	Type AttrType `json:"type"`
}

// ListCollectionsResponse represents the response for listing collections
type ListCollectionsResponse struct {
	Success         bool         `json:"success"`
	Message         string       `json:"message"`
	Data            []Collection `json:"data"`
	SupportMetadata bool         `json:"support_metadata"`
}

// AddCollectionRequest represents the request to add a new collection
type AddCollectionRequest struct {
	Name               string `json:"name"`
	NoReferenceStorage bool   `json:"no_reference_storage,omitempty"`
	HasMetadataStorage bool   `json:"has_metadata_storage,omitempty"`
}

// InsertRecordRequest represents the request to insert a record
type InsertRecordRequest struct {
	Collection        string                 `json:"collection"`
	Expiry            int64                  `json:"expiry,omitempty"`
	ID                string                 `json:"id,omitempty"`
	Record            map[string]interface{} `json:"record"`
	MetadataFields    map[string]AttrType    `json:"metadata_fields,omitempty"`
	EmbeddingProvider string                 `json:"embedding_provider,omitempty"`
	Fields            []string               `json:"fields,omitempty"`
	KeywordFields     []string               `json:"keyword_fields,omitempty"`
	Model             string                 `json:"model,omitempty"`
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
	ID             string                 `json:"id"`
	Expiry         int64                  `json:"expiry"`
	Fields         map[string]interface{} `json:"fields"`
	KeywordFields  map[string]bool        `json:"keyword_fields,omitempty"`
	MetadataFields map[string]int         `json:"metadata_fields,omitempty"`
}

// IngestRequest represents the request to ingest data
type IngestRequest struct {
	FilePath              string              `json:"file_path"`
	CollectionName        string              `json:"collection_name"`
	KeywordFields         []string            `json:"keyword_fields,omitempty"`
	MetadataFields        map[string]AttrType `json:"metadata_fields,omitempty"`
	Fields                []string            `json:"fields"`
	IdField               string              `json:"id_field,omitempty"`
	ExpiryField           string              `json:"expiry_field,omitempty"`
	EmbeddingProviderName string              `json:"embedding_provider,omitempty"`
	EmbeddingModel        string              `json:"embedding_model,omitempty"`
	IngestionBatchSize    int                 `json:"ingestion_batch_size,omitempty"`
}

// IngestResponse represents the response for data ingestion
type IngestResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}

// SearchRequest represents the request body for POST search
type SearchRequest struct {
	Collection  string             `json:"collection"`
	Query       string             `json:"query"`
	Fields      []string           `json:"fields,omitempty"`
	Limit       int                `json:"limit,omitempty"`
	Weights     map[string]float64 `json:"weights,omitempty"`
	MaxDistance *float64           `json:"max_distance,omitempty"`
	Filters     CompoundFilter     `json:"filters,omitempty"`
	Sort        CompoundSort       `json:"sort,omitempty"`
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
	Success bool                `json:"success"`
	Message string              `json:"message"`
	Data    []map[string]string `json:"data"`
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

// EmbeddingModel represents an embedding model
type EmbeddingModel struct {
	Name      string `json:"name"`
	IsDefault bool   `json:"is_default"`
}

// EmbeddingProvider represents an embedding provider with its models
type EmbeddingProvider struct {
	Name      string           `json:"name"`
	IsDefault bool             `json:"is_default"`
	Models    []EmbeddingModel `json:"models"`
}

// ListEmbeddingModelsResponse represents the response for listing embedding models
type ListEmbeddingModelsResponse struct {
	Success                      bool                `json:"success"`
	Message                      string              `json:"message"`
	Data                         []EmbeddingProvider `json:"data"`
	SupoortsDistributedEmbedding bool                `json:"supports_distributed_embedding"`
}

// AttrType represents the type of a metadata attribute
type AttrType int

const (
	AttrTypeInt64 AttrType = iota
	AttrTypeFloat64
	AttrTypeString
	AttrTypeBool
)

func (t AttrType) String() string {
	switch t {
	case AttrTypeInt64:
		return "int64"
	case AttrTypeFloat64:
		return "float64"
	case AttrTypeString:
		return "string"
	case AttrTypeBool:
		return "bool"
	default:
		return "unknown"
	}
}

// FilterOp represents a filter operation
type FilterOp int

const (
	OpEquals FilterOp = iota
	OpNotEquals
	OpGreaterThan
	OpGreaterThanOrEqual
	OpLessThan
	OpLessThanOrEqual
	OpIn
	OpNotIn
)

func (op FilterOp) String() string {
	switch op {
	case OpEquals:
		return "="
	case OpNotEquals:
		return "!="
	case OpGreaterThan:
		return ">"
	case OpGreaterThanOrEqual:
		return ">="
	case OpLessThan:
		return "<"
	case OpLessThanOrEqual:
		return "<="
	case OpIn:
		return "IN"
	case OpNotIn:
		return "NOT IN"
	default:
		return "unknown"
	}
}

// FilterExpression represents a single filter condition
type FilterExpression struct {
	Attribute string   `json:"attribute,omitempty"`
	Op        FilterOp `json:"op,omitempty"`
	Value     any      `json:"value,omitempty"`
	Values    []any    `json:"values,omitempty"`
}

// Validate checks if the filter expression is valid
func (f *FilterExpression) Validate() error {
	if f.Attribute == "" {
		return fmt.Errorf("attribute name cannot be empty")
	}

	switch f.Op {
	case OpIn, OpNotIn:
		if len(f.Values) == 0 {
			return fmt.Errorf("IN/NOT IN operations require at least one value")
		}
	default:
		if f.Value == nil {
			return fmt.Errorf("value cannot be nil for operation %s", f.Op)
		}
	}

	return nil
}

// CompoundFilter represents a combination of filter expressions
type CompoundFilter struct {
	And []FilterExpression `json:"and,omitempty"`
	// Or  []FilterExpression `json:"or,omitempty"`
}

// SortOrder represents the sort direction
type SortOrder int

const (
	SortAscending SortOrder = iota
	SortDescending
)

func (s SortOrder) String() string {
	switch s {
	case SortAscending:
		return "ASC"
	case SortDescending:
		return "DESC"
	default:
		return "unknown"
	}
}

// SortExpression represents a sort criterion
type SortExpression struct {
	Attribute string    `json:"attribute"`
	Order     SortOrder `json:"order"`
}

// Validate checks if the sort expression is valid
func (s *SortExpression) Validate() error {
	if s.Attribute == "" {
		return fmt.Errorf("sort attribute cannot be empty")
	}
	if s.Order != SortAscending && s.Order != SortDescending {
		return fmt.Errorf("invalid sort order: %d", s.Order)
	}
	return nil
}

type CompoundSort struct {
	Sorts []SortExpression `json:"sorts,omitempty"`
}

// OplogStatusResponse represents oplog status
type OplogStatusResponse struct {
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	LastLSN      uint64 `json:"last_lsn"`
	RetentionLSN uint64 `json:"retention_lsn"`
	ReplicaCount int    `json:"replica_count"`
}

// UpdateReplicaLSNRequest represents replica LSN update
type UpdateReplicaLSNRequest struct {
	Collection string `json:"collection"`
	ReplicaID  string `json:"replica_id"`
	LSN        uint64 `json:"lsn"`
}

// UpdateReplicaLSNResponse represents the update response
type UpdateReplicaLSNResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// RegisterReplicaRequest represents replica registration
type RegisterReplicaRequest struct {
	ReplicaID string `json:"replica_id"`
}

// UnRegisterReplicaRequest represents replica unregistration
type UnRegisterReplicaRequest struct {
	ReplicaID string `json:"replica_id"`
}

// GetOplogRequest represents the query parameters for oplog endpoint
type GetOplogRequest struct {
	AfterLSN uint64 `query:"after_lsn"`
	Limit    int    `query:"limit"`
}

// GetOplogResponse represents the oplog response
type GetOplogResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Entries []*OplogEntry `json:"entries"`
	LastLSN uint64        `json:"last_lsn"`
	Count   int           `json:"count"`
}

// OpType represents the type of operation in the oplog
type OpType string

const (
	OpTypeInsert           OpType = "insert"
	OpTypeUpdate           OpType = "update"
	OpTypeDelete           OpType = "delete"
	OpTypeDropCollection   OpType = "drop_collection"
	OpTypeRenameCollection OpType = "rename_collection"
)

// OplogEntry represents a single entry in the operation log
// Every oplog entry must be self-sufficient for replay
type OplogEntry struct {
	// LSN is the Log Sequence Number - monotonically increasing, gap-free
	LSN uint64 `json:"lsn"`

	// Timestamp when the operation occurred
	Timestamp time.Time `json:"timestamp"`

	// Collection name this operation applies to
	Collection string `json:"collection"`

	// DocID is the document identifier
	DocID string `json:"doc_id"`

	// OpType indicates the operation type
	OpType OpType `json:"op_type"`

	// Vector data (optional, for vector fields)
	Vector []float32 `json:"vector,omitempty"`

	// Metadata (optional, for metadata fields)
	Metadata map[string]any `json:"metadata,omitempty"`

	// Keywords (optional, for keyword search)
	Keywords []string `json:"keywords,omitempty"`

	// FullDoc is required for insert operations
	// Contains the complete document state
	FullDoc *Record `json:"full_doc,omitempty"`

	// Vectors map for multi-vector support
	Vectors map[string][]float32 `json:"vectors,omitempty"`

	// Fields for the document
	Fields map[string]interface{} `json:"fields,omitempty"`

	// KeywordFields tracks which fields are keywords
	KeywordFields map[string]bool `json:"keyword_fields,omitempty"`

	// MetadataFields tracks metadata field types
	MetadataFields map[string]AttrType `json:"metadata_fields,omitempty"`

	// Expiry time for the document (if applicable) - Unix timestamp
	Expiry int64 `json:"expiry,omitempty"`

	// NewName is used for rename operations to track the new collection name
	NewName string `json:"new_name,omitempty"`
}

type Record struct {
	Id             string                 `json:"id"`
	Fields         map[string]interface{} `json:"fields"`
	KeywordFields  map[string]bool        `json:"keyword_fields,omitempty"`
	MetadataFields map[string]AttrType    `json:"metadata_fields,omitempty"`
	Vectors        map[string][]float32   `json:"-"`
	Dist           float32                `json:"-"`
	Nodes          []string               `json:"-"`
	Expiry         int64                  `json:"expiry,omitempty"`
}
