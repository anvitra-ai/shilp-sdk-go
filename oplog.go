package shilp

import (
	"fmt"
	"net/http"
)

// GetOplogEntries retrieves oplog entries after a specific LSN for replica synchronization
// Parameters:
//   - collection: Collection name to retrieve oplog entries for (optional, if empty returns entries for all collections)
//   - afterLSN: LSN after which to retrieve oplog entries (required)
//   - limit: Maximum number of oplog entries to retrieve (optional)
func (c *Client) GetOplogEntries(collection string, afterLSN uint64, limit int) (*GetOplogResponse, error) {
	queryParams := map[string]string{
		"after_lsn": fmt.Sprintf("%d", afterLSN),
	}

	if collection != "" {
		queryParams["collection"] = collection
	}

	if limit > 0 {
		queryParams["limit"] = fmt.Sprintf("%d", limit)
	}

	var result GetOplogResponse
	err := c.doRequest(http.MethodGet, "/api/oplog/v1/", nil, &result, queryParams)
	return &result, err
}

// UpdateReplicaLSN updates the last applied LSN for a replica (heartbeat)
// Parameters:
//   - collection: Collection name
//   - replicaID: Replica identifier
//   - lsn: Last applied LSN
func (c *Client) UpdateReplicaLSN(collection, replicaID string, lsn uint64) (*UpdateReplicaLSNResponse, error) {
	req := UpdateReplicaLSNRequest{
		Collection: collection,
		ReplicaID:  replicaID,
		LSN:        lsn,
	}

	var result UpdateReplicaLSNResponse
	err := c.doRequest(http.MethodPost, "/api/oplog/v1/heartbeat", req, &result, nil)
	return &result, err
}

// RegisterReplica registers a replica for oplog retention tracking
// Parameters:
//   - replicaID: Replica identifier
func (c *Client) RegisterReplica(replicaID string) (*GenericResponse, error) {
	req := RegisterReplicaRequest{
		ReplicaID: replicaID,
	}

	var result GenericResponse
	err := c.doRequest(http.MethodPost, "/api/oplog/v1/register", req, &result, nil)
	return &result, err
}

// UnRegisterReplica unregisters a replica for oplog retention tracking
// Parameters:
//   - replicaID: Replica identifier
func (c *Client) UnRegisterReplica(replicaID string) (*GenericResponse, error) {
	req := UnRegisterReplicaRequest{
		ReplicaID: replicaID,
	}

	var result GenericResponse
	err := c.doRequest(http.MethodPost, "/api/oplog/v1/unregister", req, &result, nil)
	return &result, err
}

// GetOplogStatus retrieves current oplog status and statistics for a collection
// Parameters:
//   - collection: Collection name to retrieve oplog status for (required)
func (c *Client) GetOplogStatus(collection string) (*OplogStatusResponse, error) {
	queryParams := map[string]string{
		"collection": collection,
	}

	var result OplogStatusResponse
	err := c.doRequest(http.MethodGet, "/api/oplog/v1/status", nil, &result, queryParams)
	return &result, err
}
