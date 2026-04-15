package shilp

import (
	"net/http"
)

// GetSettings retrieves the settings
func (c *Client) GetSettings() (*GetSettingsResponse, error) {
	var result GetSettingsResponse
	err := c.doRequest(http.MethodGet, "/api/settings/v1/", nil, &result, nil)
	return &result, err
}

// UpdateSettings updates the settings
func (c *Client) UpdateSettings(settings SettingsUpdateRequest) (*GenericResponse, error) {
	var result GenericResponse
	err := c.doRequest(http.MethodPut, "/api/settings/v1/", settings, &result, nil)
	return &result, err
}

func (c *Client) ListProviders() (*SettingsAvailableProvidersResponse, error) {
	var result SettingsAvailableProvidersResponse
	err := c.doRequest(http.MethodGet, "/api/settings/v1/providers", nil, &result, nil)
	return &result, err
}
