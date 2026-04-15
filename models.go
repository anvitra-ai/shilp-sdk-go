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

// IndexType represents the type of index for a column
type IndexType string

const (
	IndexTypeHNSW     IndexType = "hnsw"
	IndexTypeInverted IndexType = "inverted"
	IndexTypeMetadata IndexType = "metadata"
)

// Collection represents a collection in the database
type Collection struct {
	Name                 string                 `json:"name"`
	IsLoaded             bool                   `json:"is_loaded"`
	Fields               []string               `json:"fields"`
	SearchableFields     []string               `json:"searchable_fields"`
	FieldConfig          map[string]IndexType   `json:"field_config"`
	Metadata             []MetadataColumnSchema `json:"metadata,omitempty"`
	HasMetadataEnabled   bool                   `json:"has_metadata_enabled"`
	NoReferenceStorage   bool                   `json:"no_reference_storage"`
	StorageType          StorageBackendType     `json:"storage_type"`
	ReferenceStorageType StorageBackendType     `json:"reference_storage_type"`
	IsPQEnabled          bool                   `json:"is_pq_enabled"`
	IsNLIEnabled         bool                   `json:"is_nli_enabled,omitempty"`
	NLIDomain            string                 `json:"nli_domain,omitempty"`
	TotalNoofDocuments   int                    `json:"total_no_of_documents"`
}

// StorageBackendType represents the type of storage backend available to store the data for persistance
type StorageBackendType int

const (
	StorageBackendDoesnotExist StorageBackendType = -1
	StorageBackendFile         StorageBackendType = iota
	StorageBackendS3
)

func (s StorageBackendType) String() string {
	switch s {
	case StorageBackendFile:
		return "filesystem"
	case StorageBackendS3:
		return "s3"
	default:
		return "unknown"
	}
}

func (s StorageBackendType) IsValid() bool {
	switch s {
	case StorageBackendFile, StorageBackendS3:
		return true
	default:
		return false
	}
}

type EnableMetadataStoreRequest struct {
	Fields []MetadataColumnSchema `json:"fields,omitempty"`
}

type EnableMetadataStoreResponse struct {
	Success        bool   `json:"success"`
	Message        string `json:"message"`
	RecordsIndexed int    `json:"records_indexed,omitempty"`
}

type MetadataColumnSchema struct {
	Name string   `json:"name"`
	Type AttrType `json:"type"`
}

// ListCollectionsResponse represents the response for listing collections
type ListCollectionsResponse struct {
	Success        bool                  `json:"success"`
	Message        string                `json:"message"`
	Data           []Collection          `json:"data"`
	MetadataInfo   []MetadataSupportInfo `json:"metadata_info"`
	IsNliSupported bool                  `json:"is_nli_supported"`
}

type MetadataSupportInfo struct {
	SupportMetadata bool               `json:"support_metadata"`
	Name            string             `json:"name"`
	Type            StorageBackendType `json:"type"`
	IsDefault       bool               `json:"is_default"`
}

// AddCollectionRequest represents the request to add a new collection
type AddCollectionRequest struct {
	Name                 string             `json:"name"`
	NoReferenceStorage   bool               `json:"no_reference_storage"`
	HasMetadataStorage   bool               `json:"has_metadata_storage"`
	StorageType          StorageBackendType `json:"storage_type"`
	ReferenceStorageType StorageBackendType `json:"reference_storage_type"`
	EnablePQ             bool               `json:"enable_pq"`
}

type GetCollectionDataResponse struct {
	Data    []CollectionDataRecord `json:"data"`
	Total   int                    `json:"total"`
	Success bool                   `json:"success"`
	Message string                 `json:"message"`
}

type CollectionDataRecord struct {
	ID      string                 `json:"id"`
	Data    map[string]interface{} `json:"data"`
	Vectors map[string][]float32   `json:"vectors,omitempty"`
}

type GetCollectionSchemaResponse struct {
	Message string            `json:"message,omitempty"`
	Success bool              `json:"success"`
	Data    *CollectionSchema `json:"data,omitempty"`
}
type CollectionSchema struct {
	Attributes  []Attribute      `json:"attributes,omitempty"`
	ValueSchema []CategorySchema `json:"value_schema,omitempty"`
}

type Attribute struct {
	Name       string        `json:"name,omitempty"`
	Type       AttributeType `json:"type,omitempty"`
	IndexType  IndexType     `json:"index_type,omitempty"`
	IsMetadata bool          `json:"is_metadata,omitempty"`
}

type AttributeType int

const (
	AttributeTypeNumerical AttributeType = 1
	AttributeTypeString    AttributeType = 2
)

type CategorySchema struct {
	Name      string          `json:"name,omitempty"`
	IndexType IndexType       `json:"index_type,omitempty"`
	Values    []CategoryValue `json:"values,omitempty"`
	Synonyms  []string        `json:"synonyms,omitempty"`
}

type CategoryValue struct {
	Value string `json:"value,omitempty"`
	Count int    `json:"count,omitempty"`
}

// InsertRecordRequest represents the request to insert a record
type InsertRecordRequest struct {
	Collection        string                 `json:"collection"`
	Expiry            int64                  `json:"expiry,omitempty"`
	ID                string                 `json:"id,omitempty"`
	Record            map[string]interface{} `json:"record"`
	MetadataFields    map[string]AttrType    `json:"metadata_fields,omitempty"`
	EmbeddingProvider string                 `json:"embedding_provider,omitempty"`
	// If a vector field is present, it will be used instead of embedding generation
	Fields        []string `json:"fields,omitempty"`
	KeywordFields []string `json:"keyword_fields,omitempty"`
	// This is the map of field name to vector data
	Vectors      map[string][]float32          `json:"vectors,omitempty"`
	Model        string                        `json:"model,omitempty"`
	VectorConfig map[string]VectorCreateConfig `json:"vector_config,omitempty"`
	// ArrayFields specifies which fields in the record are comma separated string that have individual meaning. Separate embeddings are generated and centroid is stored for the entire field.
	ArrayFields []string `json:"array_fields,omitempty"`
}

type VectorCreateConfig struct {
	EfConstruction int `json:"ef_construction,omitempty"`
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
	// Source configuration - use either FilePath OR MongoDB settings
	//FilePath has the path to the file to be ingested or `database/collection` for MongoDB
	FilePath   string           `json:"file_path,omitempty"`
	SourceType IngestSourceType `json:"source_type,omitempty"` // "file" or "mongodb"

	// MongoDB source configuration
	DatabaseName        string         `json:"database_name,omitempty"`
	MongoCollection     string         `json:"mongo_collection,omitempty"`
	Query               map[string]any `json:"query,omitempty"`
	MongoFetchBatchSize int            `json:"mongo_fetch_batch_size,omitempty"`

	// Common configuration
	CollectionName        string              `json:"collection_name"`
	KeywordFields         []string            `json:"keyword_fields,omitempty"`
	MetadataFields        map[string]AttrType `json:"metadata_fields,omitempty"`
	Fields                []string            `json:"fields"`
	ArrayFields           []string            `json:"array_fields,omitempty"`
	IdField               string              `json:"id_field,omitempty"`
	ExpiryField           string              `json:"expiry_field,omitempty"`
	EmbeddingProviderName string              `json:"embedding_provider,omitempty"`
	EmbeddingModel        string              `json:"embedding_model,omitempty"`
	IngestionBatchSize    int                 `json:"ingestion_batch_size,omitempty"`

	VectorConfig map[string]VectorCreateConfig `json:"vector_config,omitempty"`
}

type IngestSourceType string

const (
	IngestSourceTypeFile    IngestSourceType = "file"
	IngestSourceTypeMongoDB IngestSourceType = "mongodb"
	IngestSourceTypeAnvitra IngestSourceType = "anvitra"
)

func (i IngestSourceType) IsValid() bool {
	switch i {
	case IngestSourceTypeFile, IngestSourceTypeMongoDB, IngestSourceTypeAnvitra:
		return true
	default:
		return false
	}
}

// IngestResponse represents the response for data ingestion
type IngestResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}

type ListIngestionSourcesResponse struct {
	Message string             `json:"message"`
	Success bool               `json:"success"`
	Data    []IngestSourceType `json:"data,omitempty"`
}

type VerticalInfo struct {
	Name     string `json:"name,omitempty"`
	Label    string `json:"label,omitempty"`
	IsNative bool   `json:"is_native,omitempty"`
}
type ListNLIVerticalsResponse struct {
	Success bool           `json:"success"`
	Data    []VerticalInfo `json:"data,omitempty"`
	Message string         `json:"message,omitempty"`
}

type FileReaderOptions struct {
	Source      IngestSourceType       `json:"source,omitempty"`
	MongoFilter map[string]interface{} `json:"mongo_filter,omitempty"`
	Skip        int                    `json:"skip,omitempty"`
	Limit       int                    `json:"limit,omitempty"`
}

// SearchRequest represents the request body for POST search
type SearchRequest struct {
	Collection    string                        `json:"collection"`
	Query         string                        `json:"query"`
	Fields        []string                      `json:"fields,omitempty"`
	Limit         int                           `json:"limit,omitempty"`
	Weights       map[string]float64            `json:"weights,omitempty"`
	MaxDistance   *float64                      `json:"max_distance,omitempty"`
	Filters       CompoundFilter                `json:"filters,omitempty"`
	Sort          CompoundSort                  `json:"sort,omitempty"`
	VectorQuery   []float32                     `json:"vector_query,omitempty"`
	UseNli        bool                          `json:"use_nli,omitempty"`
	FieldConfig   map[string]VectorSearchConfig `json:"field_config,omitempty"`
	Queries       map[string]string             `json:"queries,omitempty"`        // For multi-query search, key is the query name and value is the query string
	VectorQueries map[string][]float32          `json:"vector_queries,omitempty"` // For multi-query search, key is the query name and value is the vector query
	FuzzyAlgo     FuzzyAlgo                     `json:"fuzzy_algo,omitempty"`
}

type FuzzyAlgo string

// Fuzzy algorithm identifiers.
const (
	FuzzyAlgoLevenshtein FuzzyAlgo = "levenshtein"
	FuzzyAlgoJaroWinkler FuzzyAlgo = "jaro_winkler"
)

func (f FuzzyAlgo) IsValid() bool {
	return f == FuzzyAlgoLevenshtein || f == FuzzyAlgoJaroWinkler
}

type VectorSearchConfig struct {
	EfSearch int `json:"ef_search,omitempty"`
}

// SearchResponse represents the response for searching data
type SearchResponse struct {
	Success        bool                     `json:"success"`
	Message        string                   `json:"message"`
	Data           []map[string]interface{} `json:"data"`
	Interpretation *Query                   `json:"interpretation,omitempty"`
	Timing         *SearchTiming            `json:"timing,omitempty"`
}

// SearchTiming holds the duration (in milliseconds) for each step of a search request.
type SearchTiming struct {
	InterpretationMs int64 `json:"interpretation_ms,omitempty"`
	EmbeddingMs      int64 `json:"embedding_ms,omitempty"`
	MetadataFilterMs int64 `json:"metadata_filter_ms,omitempty"`
	SearchMs         int64 `json:"search_ms,omitempty"`
	TotalMs          int64 `json:"total_ms"`
}

type Query struct {
	VectorQuery  VectorQuery   `json:"vector_query"`
	Filters      []Filter      `json:"filters"`
	ValueFilters []ValueFilter `json:"value_filters"`
}

type VectorQuery struct {
	ResolvedBy        []string           `json:"resolved_by"`
	VectorQuery       string             `json:"vector_query"`
	VectorQueries     map[string]string  `json:"vector_queries,omitempty"`
	VectorConfidences map[string]float32 `json:"vector_confidences,omitempty"`
}

type Filter struct {
	ResolvedBy     []string        `json:"resolved_by"`
	Attribute      []Token         `json:"attribute"`
	Operation      Token           `json:"operation"`
	Operator       FilterOperator  `json:"operator"`
	Value          []Token         `json:"value"`
	IsNumerical    bool            `json:"is_numerical"`
	Grounded       bool            `json:"grounded"`
	NumericalValue *NumericalValue `json:"numerical_value,omitempty"`
}

type FilterOperator string

const (
	FilterOperatorEquals       FilterOperator = "EQ"
	FilterOperatorNotEquals    FilterOperator = "NEQ"
	FilterOperatorGreaterThan  FilterOperator = "GT"
	FilterOperatorLessThan     FilterOperator = "LT"
	FilterOperatorGreaterEqual FilterOperator = "GTE"
	FilterOperatorLessEqual    FilterOperator = "LTE"
	FilterOperatorIn           FilterOperator = "IN"
	FilterOperatorNotIn        FilterOperator = "NOT IN"
)

func (f FilterOperator) ToFilterOp() FilterOp {
	switch f {
	case FilterOperatorEquals:
		return OpEquals
	case FilterOperatorNotEquals:
		return OpNotEquals
	case FilterOperatorGreaterThan:
		return OpGreaterThan
	case FilterOperatorLessThan:
		return OpLessThan
	case FilterOperatorGreaterEqual:
		return OpGreaterThanOrEqual
	case FilterOperatorLessEqual:
		return OpLessThanOrEqual
	case FilterOperatorIn:
		return OpIn
	case FilterOperatorNotIn:
		return OpNotIn
	default:
		return OpUnknown
	}
}

type ValueFilter struct {
	ResolvedBy []string       `json:"resolved_by"`
	Attribute  []Token        `json:"attribute"`
	Values     [][]Token      `json:"values"`
	Grounded   bool           `json:"grounded"`
	Operator   FilterOperator `json:"operator"`
}

type Token struct {
	Text  string `json:"text"`
	Tag   string `json:"tag"`
	Label string `json:"label"`
}

type NumericalValue struct {
	Unit         string  `json:"unit"`
	BaseValue    float64 `json:"base_value"`
	Multiplier   float64 `json:"multiplier"`
	TotalValue   float64 `json:"total_value"`
	OriginalText string  `json:"original_text"`
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

type DebugGetEmbeddingsResponse struct {
	Data    [][]float32 `json:"data"`
	Success bool        `json:"success"`
	Message string      `json:"message"`
}

type DebugGetEmbeddingsRequest struct {
	Texts []string `json:"texts"`
}

// DebugDistanceResponse represents the response for debug distance endpoint
type DebugDistanceResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		Distance              float64   `json:"distance"`
		Vector                []float64 `json:"vector"`
		CustomMatcherDistance float32   `json:"custom_matcher_distance,omitempty"`
		CustomMatcherVector   []float32 `json:"custom_matcher_vector,omitempty"`
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
	AttrTypeCurrency
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
	case AttrTypeCurrency:
		return "currency"
	default:
		return "unknown"
	}
}

// FilterOp represents a filter operation
type FilterOp int

const (
	OpUnknown                     = -1
	OpEquals             FilterOp = 0
	OpNotEquals                   = 1
	OpGreaterThan                 = 2
	OpGreaterThanOrEqual          = 3
	OpLessThan                    = 4
	OpLessThanOrEqual             = 5
	OpIn                          = 6
	OpNotIn                       = 7
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
	Attribute string          `json:"attribute,omitempty"`
	Op        FilterOp        `json:"op,omitempty"`
	Value     any             `json:"value,omitempty"`
	Values    []any           `json:"values,omitempty"`
	Filters   *CompoundFilter `json:"filters,omitempty"`
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
	Or  []FilterExpression `json:"or,omitempty"`
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

// Status represents the overall status of the registry
type Status struct {
	WriteReplica Replica    `json:"write_replica"`
	ReadReplicas []*Replica `json:"read_replicas"`
	Available    int        `json:"available_count"`
	Total        int        `json:"total_count"`
}

type Replica struct {
	Id        string `json:"id"`
	Address   string `json:"address"`
	IsHealthy bool   `json:"is_healthy"`
	IsSyncing bool   `json:"is_syncing"` // Traffic gate - if true, no traffic sent
}

type ProxyStats struct {
	ActiveProxies int      `json:"active_proxies"`
	Targets       []string `json:"targets"`
}

type DiscoveryStats struct {
	Registry Status     `json:"registry"`
	Proxy    ProxyStats `json:"proxy"`
}

type SyncStatus string

const (
	SyncStatusReady   SyncStatus = "ready"
	SyncStatusSyncing SyncStatus = "syncing"
)

type UpdateSyncStatusRequest struct {
	AccountID string     `json:"account_id"`
	Address   string     `json:"address"`
	Status    SyncStatus `json:"status"`
}

type RegisterToDiscoveryRequest struct {
	AccountID string `json:"account_id"`
	Address   string `json:"address"`
	Id        string `json:"id"`
	IsRead    bool   `json:"is_read"`
	IsWrite   bool   `json:"is_write"`
}

type ReplicaType int

const (
	ReadReplica ReplicaType = iota
	WriteReplica
	SingleNode
)

func (rt ReplicaType) IsRead() bool {
	return rt == ReadReplica
}
func (rt ReplicaType) IsWrite() bool {
	return rt == WriteReplica
}

func (rt ReplicaType) IsSingleNode() bool {
	return rt == SingleNode
}

type GetSettingsResponse struct {
	Success bool      `json:"success"`
	Message string    `json:"message"`
	Data    *Settings `json:"data,omitempty"`
}

type Settings struct {
	Auth           SettingsAuth          `json:"auth"`
	AllowedOrigins []string              `json:"allowedOrigins,omitempty"`
	Integrations   []SettingsIntegration `json:"integrations,omitempty"`
}

type SettingsAuth struct {
	Enable        bool                    `json:"enable"`
	Tested        bool                    `json:"tested"`
	Name          string                  `json:"name,omitempty"`
	Arguments     []ProviderArgumentValue `json:"arguments,omitempty"`
	APIAuthConfig APIAuthConfig           `json:"apiAuthConfig,omitempty"`
}

type APIAuthConfig struct {
	Search      bool `json:"search"`
	Collections bool `json:"collections"`
	Data        bool `json:"data"`
	Explore     bool `json:"explore"`
	Oplog       bool `json:"oplog"`
}

type ProviderArgumentValue struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	IsSecret bool   `json:"is_secret,omitempty"`
}

type SettingsIntegration struct {
	Enable    bool                    `json:"enable"`
	Name      string                  `json:"name,omitempty"`
	Arguments []ProviderArgumentValue `json:"arguments,omitempty"`
}

type SettingsUpdateRequest struct {
	Auth           *SettingsAuth                  `json:"auth,omitempty"`
	Tested         *bool                          `json:"tested,omitempty"`
	AuthConfig     *SettingsAuth                  `json:"authConfig,omitempty"`
	AllowedOrigins *[]string                      `json:"allowedOrigins,omitempty"`
	Integration    map[string]SettingsIntegration `json:"integration,omitempty"`
}

type SettingsAvailableProvidersResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		Auth         []SettingsProviderInfo `json:"auth,omitempty"`
		Integrations []SettingsProviderInfo `json:"integrations,omitempty"`
	} `json:"data,omitempty"`
}

type SettingsProviderArguments struct {
	Label       string `json:"label"`
	Description string `json:"description"`
	IsSecret    bool   `json:"is_secret,omitempty"`
}

type SettingsProviderType string

const (
	AuthProviderType       SettingsProviderType = "auth"
	DataSourceProviderType SettingsProviderType = "data-source"
)

type SettingsProviderInfo struct {
	Name      string                      `json:"name"`
	Type      SettingsProviderType        `json:"type"`
	Arguments []SettingsProviderArguments `json:"arguments,omitempty"`
}
