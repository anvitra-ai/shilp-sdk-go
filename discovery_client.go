package shilp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// DiscoveryClient is the client for the Discovery API
type DiscoveryClient struct {
	baseURL    string
	httpClient *http.Client
}

// DiscoveryClientOption is a function that configures the DiscoveryClient
type DiscoveryClientOption func(*DiscoveryClient)

// NewDiscoveryClient creates a new Discovery API client
func NewDiscoveryClient(baseURL string, opts ...DiscoveryClientOption) *DiscoveryClient {
	c := &DiscoveryClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// doRequest performs an HTTP request
func (d *DiscoveryClient) doRequest(method, path string, body interface{}, result interface{}, queryParams map[string]string) error {
	req, err := d.prepareRequest(method, path, body, queryParams)
	if err != nil {
		return err
	}

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error: %s (status: %d)", string(bodyBytes), resp.StatusCode)
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

func (d *DiscoveryClient) prepareRequest(method, path string, body interface{}, queryParams map[string]string) (*http.Request, error) {
	u, err := url.Parse(d.baseURL + path)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	if len(queryParams) > 0 {
		q := u.Query()
		for k, v := range queryParams {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
	}

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, u.String(), reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (d *DiscoveryClient) GetShilpStats(accountId string) (*DiscoveryStats, error) {
	var stats DiscoveryStats
	err := d.doRequest(http.MethodGet, "/control/shilp/stats", nil, &stats, map[string]string{"account_id": accountId})
	if err != nil {
		return nil, fmt.Errorf("failed to get discovery stats: %v", err)
	}
	return &stats, nil
}

func (d *DiscoveryClient) UpdateShilpSyncStatus(accountId string, address string, status SyncStatus) (*GenericResponse, error) {
	req := UpdateSyncStatusRequest{
		Status:    status,
		AccountID: accountId,
		Address:   address,
	}
	res := &GenericResponse{}
	err := d.doRequest(http.MethodPost, "/control/shilp/sync-status", req, res, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to update sync status: %v", err)
	}
	return res, nil
}

func (d *DiscoveryClient) RegisterShilpService(accountId string, address string, id string, replicaType ReplicaType) error {
	// Register with node manager
	registrationData := RegisterToDiscoveryRequest{
		AccountID: accountId,
		Address:   address,
		Id:        id,
	}

	switch replicaType {
	case ReadReplica:
		registrationData.IsRead = true
	case WriteReplica:
		registrationData.IsWrite = true
	}
	if replicaType != SingleNode {
		_, err := d.registerationShilp(registrationData, true)
		if err != nil {
			return fmt.Errorf("failed to register service: %v", err)
		}
		return nil
	}
	registrationData.IsRead = true
	_, err := d.registerationShilp(registrationData, true)
	if err != nil {
		return fmt.Errorf("failed to register read replica: %v", err)
	}
	registrationData.IsRead = false
	registrationData.IsWrite = true
	_, err = d.registerationShilp(registrationData, true)
	if err != nil {
		return fmt.Errorf("failed to register write replica: %v", err)
	}
	return nil
}

func (d *DiscoveryClient) UnregisterShilpService(accountId string, address string, id string, replicaType ReplicaType) error {
	// Register with node manager
	registrationData := RegisterToDiscoveryRequest{
		AccountID: accountId,
		Address:   address,
		Id:        id,
	}

	switch replicaType {
	case ReadReplica:
		registrationData.IsRead = true
	case WriteReplica:
		registrationData.IsWrite = true
	}
	if replicaType != SingleNode {
		_, err := d.registerationShilp(registrationData, false)
		if err != nil {
			return fmt.Errorf("failed to unregister service: %v", err)
		}
		return nil
	}
	registrationData.IsRead = true
	_, err := d.registerationShilp(registrationData, false)
	if err != nil {
		return fmt.Errorf("failed to unregister read replica: %v", err)
	}
	registrationData.IsRead = false
	registrationData.IsWrite = true
	_, err = d.registerationShilp(registrationData, false)
	if err != nil {
		return fmt.Errorf("failed to unregister write replica: %v", err)
	}
	return nil
}

func (d *DiscoveryClient) registerationShilp(payload RegisterToDiscoveryRequest, isRegistration bool) (*GenericResponse, error) {

	endpoint := "register"
	if !isRegistration {
		endpoint = "unregister"
	}

	res := &GenericResponse{}
	err := d.doRequest(http.MethodPost, fmt.Sprintf("/control/shilp/%s", endpoint), payload, res, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to update sync status: %v", err)
	}
	return res, nil
}

func (d *DiscoveryClient) RegisterTeiService(accountId string, address string, id string) error {
	// Register with node manager
	registrationData := RegisterToDiscoveryRequest{
		AccountID: accountId,
		Address:   address,
		Id:        id,
		IsRead:    true,
	}
	err := d.doRequest(http.MethodPost, "/control/tei/register", registrationData, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to register service: %v", err)
	}
	return nil
}

func (d *DiscoveryClient) UnregisterTeiService(accountId string, address string, id string) error {
	registrationData := RegisterToDiscoveryRequest{
		AccountID: accountId,
		Address:   address,
		Id:        id,
		IsRead:    true,
	}

	err := d.doRequest(http.MethodPost, "/control/tei/unregister", registrationData, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to unregister service: %v", err)
	}
	return nil
}
