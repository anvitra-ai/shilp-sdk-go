package shilp

// HealthCheck performs a health check on the API
func (c *Client) HealthCheck() (*HealthResponse, error) {
	var result HealthResponse
	err := c.doRequest("GET", "/health", nil, &result, nil)
	return &result, err
}
