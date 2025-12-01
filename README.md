# Shilp Go SDK

This is the official Go SDK for the Shilp Vector Database API.

## Installation

```bash
go get github.com/anvitra-ai/shilp-sdk-go
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

	// Create a new collection
	_, err = client.AddCollection(shilp.AddCollectionRequest{
		Name: "my-collection",
	})
	if err != nil {
		log.Printf("Failed to add collection: %v", err)
	}

	// Insert a record
	_, err = client.InsertRecord(shilp.InsertRecordRequest{
		Collection: "my-collection",
		ID:         "record-1",
		Record: map[string]interface{}{
			"title": "Hello World",
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

	// Search
	results, err := client.SearchData("my-collection", "Hello", []string{"title"}, 10)
	if err != nil {
		log.Printf("Search failed: %v", err)
	}
	fmt.Printf("Search results: %+v\n", results.Data)
}
```

## Features

- Collection Management (List, Add, Drop, Rename, Load, Unload, Flush)
- Data Ingestion & Search
- Record Management (Insert, Delete, Expiry Cleanup)
- Storage Listing
- Health Check
