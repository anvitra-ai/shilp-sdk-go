package shilp

import (
	"fmt"
	"strconv"
)

// GetCollectionDistance gets the distance of a node in a collection for debug purposes
func (c *Client) GetCollectionDistance(collectionName, field string, nodeID int, text string) (*DebugDistanceResponse, error) {
	var result DebugDistanceResponse
	path := fmt.Sprintf("/api/collections/v1/debug/%s/%s/distance/%d", collectionName, field, nodeID)
	queryParams := map[string]string{
		"text": text,
	}
	err := c.doRequest("GET", path, nil, &result, queryParams)
	return &result, err
}

// GetCollectionNodeInfo gets node info of a collection for debug purposes
func (c *Client) GetCollectionNodeInfo(collectionName, field string, nodeID int) (*DebugNodeInfoResponse, error) {
	var result DebugNodeInfoResponse
	path := fmt.Sprintf("/api/collections/v1/debug/%s/%s/nodes/%d", collectionName, field, nodeID)
	err := c.doRequest("GET", path, nil, &result, nil)
	return &result, err
}

// GetCollectionNodeNeighborsAtLevel gets node neighbors at a level of a collection for debug purposes
func (c *Client) GetCollectionNodeNeighborsAtLevel(collectionName, field string, nodeID, level int, limit, offset *int) (*DebugNodeInfoResponse, error) {
	var result DebugNodeInfoResponse
	path := fmt.Sprintf("/api/collections/v1/debug/%s/%s/nodes/%d/neighbors/%d", collectionName, field, nodeID, level)

	queryParams := map[string]string{}
	if limit != nil {
		queryParams["limit"] = strconv.Itoa(*limit)
	}
	if offset != nil {
		queryParams["offset"] = strconv.Itoa(*offset)
	}

	err := c.doRequest("GET", path, nil, &result, queryParams)
	return &result, err
}

// GetCollectionLevels gets levels of a collection for debug purposes
func (c *Client) GetCollectionLevels(collectionName string) (*DebugLevelsResponse, error) {
	var result DebugLevelsResponse
	path := fmt.Sprintf("/api/collections/v1/debug/%s/levels", collectionName)
	err := c.doRequest("GET", path, nil, &result, nil)
	return &result, err
}

// GetCollectionNodesAtLevel gets nodes at a level of a collection for debug purposes
func (c *Client) GetCollectionNodesAtLevel(collectionName string, level int) (*DebugNodesAtLevelResponse, error) {
	var result DebugNodesAtLevelResponse
	path := fmt.Sprintf("/api/collections/v1/debug/%s/levels/%d", collectionName, level)
	err := c.doRequest("GET", path, nil, &result, nil)
	return &result, err
}

// GetCollectionNodeByReferenceNodeID gets node by reference node ID of a collection for debug purposes
func (c *Client) GetCollectionNodeByReferenceNodeID(collectionName string, nodeID int) (*DebugReferenceNodeResponse, error) {
	var result DebugReferenceNodeResponse
	path := fmt.Sprintf("/api/collections/v1/debug/%s/nodes/reference_node/%d", collectionName, nodeID)
	err := c.doRequest("GET", path, nil, &result, nil)
	return &result, err
}
