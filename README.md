# Shilp Go SDK

This is the official Go SDK for the Shilp Vector Database API.

## Installation

```bash
go get -u github.com/anvitra-ai/shilp-sdk-go
```

## Usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/anvitra-ai/shilp-sdk-go"
)

func main() {
	// Initialize the client
	client := shilp.NewClient("http://localhost:3000")

	// Check health
	health, err := client.HealthCheck()
	if err != nil {
		log.Fatalf("Health check failed: %v", err)
	}
	fmt.Printf("Health: %v\n", health.Success)

	// List collections
	collections, err := client.ListCollections()
	if err != nil {
		log.Fatalf("Failed to list collections: %v", err)
	}
	fmt.Printf("Collections: %+v\n", collections.Data)

    // Drop collection if exists
	client.DropCollection("my-collection")

	// Create a new collection
	_, err = client.AddCollection(shilp.AddCollectionRequest{
		Name: "my-collection",
	})
	if err != nil {
		log.Fatalf("Failed to create collection: %v", err)
	}

	// Insert a record
	_, err = client.InsertRecord(shilp.InsertRecordRequest{
		Collection: "my-collection",
		ID:         "record-1",
		Record: map[string]interface{}{
			"title":  "Hello World",
			"vector": []float64{0.1, 0.2, 0.3},
		},
	})
	if err != nil {
		log.Printf("Failed to insert record: %v", err)
	}
	// Flush collection incase you are using insert record.
	// Flush can be used post inserting the batch of records.
	_, err = client.FlushCollection("my-collection")
	if err != nil {
		log.Printf("Failed to flush collection: %v", err)
	}

	// Search with max distance filter
	maxDist := 0.5
	advancedResults, err := client.SearchData(shilp.SearchRequest{
		Collection:  "my-collection",
		Query:       "Hello",
		Fields:      []string{"title"},
		Limit:       10,
		MaxDistance: &maxDist,
	})
	if err != nil {
		log.Printf("Advanced search failed: %v", err)
	}
	fmt.Printf("Advanced search results: %+v\n", advancedResults.Data)

	_, err = client.DropCollection("my-collection")
	if err != nil {
		log.Printf("Failed to drop collection: %v", err)
	}
}
```

### Debug Operations

The SDK also provides debug endpoints for inspecting collection internals:

```go
// Re-index a collection
_, err = client.ReIndexCollection("my-collection")

// Get collection levels
levels, err := client.GetCollectionLevels("my-collection")
if err != nil {
	log.Printf("Failed to get levels: %v", err)
}

// Get nodes at a specific level
nodes, err := client.GetCollectionNodesAtLevel("my-collection", 0)
if err != nil {
	log.Printf("Failed to get nodes: %v", err)
}

// Get node information
nodeInfo, err := client.GetCollectionNodeInfo("my-collection", "title", 123)
if err != nil {
	log.Printf("Failed to get node info: %v", err)
}

// Get distance to a node
limit := 10
offset := 0
neighbors, err := client.GetCollectionNodeNeighborsAtLevel("my-collection", "title", 123, 0, &limit, &offset)
if err != nil {
	log.Printf("Failed to get neighbors: %v", err)
}

// Get distance calculation
distance, err := client.GetCollectionDistance("my-collection", "title", 123, "some text")
if err != nil {
	log.Printf("Failed to get distance: %v", err)
}

// Get node by reference ID
refNode, err := client.GetCollectionNodeByReferenceNodeID("my-collection", 456)
if err != nil {
	log.Printf("Failed to get reference node: %v", err)
}
```

## Features

- Collection Management (List, Add, Drop, Rename, Load, Unload, Flush, ReIndex)
- Data Ingestion & Search (with keyword fields support)
- Record Management (Insert, Delete, Expiry Cleanup)
- Debug Collection Operations (Distance, Node Info, Levels, Neighbors)
- Oplog Operations (Replica Registration, Heartbeat, Get Entries, Status)
- Storage Listing
- Health Check

### Oplog Operations

The SDK provides oplog (operation log) endpoints for replica synchronization:

```go
// Register a replica for oplog tracking
registerResp, err := client.RegisterReplica("replica-1")
if err != nil {
	log.Printf("Failed to register replica: %v", err)
}
fmt.Printf("Registered replica, last LSN: %d\n", registerResp.LastLSN)

// Get oplog status for a collection
status, err := client.GetOplogStatus("my-collection")
if err != nil {
	log.Printf("Failed to get oplog status: %v", err)
}
fmt.Printf("Oplog status - Last LSN: %d, Retention LSN: %d, Replicas: %d\n",
	status.LastLSN, status.RetentionLSN, status.ReplicaCount)

// Get oplog entries after a specific LSN
entries, err := client.GetOplogEntries("my-collection", 1000, 100)
if err != nil {
	log.Printf("Failed to get oplog entries: %v", err)
}
fmt.Printf("Retrieved %d oplog entries, last LSN: %d\n", entries.Count, entries.LastLSN)

// Get oplog entries for all collections
allEntries, err := client.GetOplogEntries("", 1000, 100)
if err != nil {
	log.Printf("Failed to get oplog entries: %v", err)
}

// Update replica LSN (heartbeat)
updateResp, err := client.UpdateReplicaLSN("my-collection", "replica-1", 1050)
if err != nil {
	log.Printf("Failed to update replica LSN: %v", err)
}
fmt.Printf("Updated replica LSN: %v\n", updateResp.Success)
```
