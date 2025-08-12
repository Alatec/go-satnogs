package gosatnogs

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

const (
	baseURL = "https://db.satnogs.org/api"
)

type urlParam struct {
	Key   string
	Value string
}

type Client struct {
	client  *http.Client
	baseURL string
	apiKey  string
}

func NewClient(apiKey string) *Client {
	return &Client{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: baseURL,
		apiKey:  apiKey,
	}
}

func (c *Client) Get(endpoint string, params []urlParam) (*http.Response, error) {
	// Create URL
	u, err := url.Parse(c.baseURL + endpoint)
	if err != nil {
		return nil, err
	}

	// Add query parameters
	q := u.Query()
	for _, param := range params {
		q.Add(param.Key, param.Value) // This automatically URL-encodes the values
	}
	u.RawQuery = q.Encode()

	// Create request
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	// Add authorization header if API key is set
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Token "+c.apiKey)
	}
	return c.client.Do(req)
}

type Telemetry struct {
	SatID         string    `json:"sat_id"`
	NoradCatID    int       `json:"norad_cat_id"`
	Transmitter   string    `json:"transmitter"`
	AppSource     string    `json:"app_source"`
	Decoded       string    `json:"decoded"`
	Frame         string    `json:"frame"`
	Observer      string    `json:"observer"`
	Timestamp     time.Time `json:"timestamp"`
	Version       string    `json:"version"`
	ObservationID int       `json:"observation_id"`
	StationID     int       `json:"station_id"`
}

type TelemetryResponse struct {
	Next    string      `json:"next"`
	Prev    string      `json:"prev"`
	Results []Telemetry `json:"results"`
}

// GetTelemetry retrieves telemetry data for a specific satellite from the SatNOGS database.
// It returns a slice of Telemetry structs containing the decoded data, or an error if the request fails.
//
// Parameters:
//   - satelliteID: The unique identifier of the satellite (typically a NORAD ID as a string)
//
// Returns:
//   - []Telemetry: A slice of Telemetry structs containing the satellite's telemetry data
//   - error: An error object if the request fails or if the response cannot be decoded
//
// Note: This function only returns the first page of results. For complete telemetry data,
// consider using GetTelemetryResponse() which has pagination support.
func (c *Client) GetTelemetry(satelliteID string) ([]Telemetry, error) {
	resp, err := c.GetTelemetryResponse(satelliteID)
	if err != nil {
		return nil, err
	}
	return resp.Results, nil
}

func (c *Client) GetTelemetryResponse(satelliteID string) (*TelemetryResponse, error) {
	resp, err := c.Get("/telemetry/", []urlParam{{"sat_id", satelliteID}, {"format", "json"}})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var telemetryResponse TelemetryResponse
	if err := json.NewDecoder(resp.Body).Decode(&telemetryResponse); err != nil {
		return nil, err
	}
	return &telemetryResponse, nil
}

func (c *Client) GetTelemetryResponseNextPage(t *TelemetryResponse) (*TelemetryResponse, error) {
	if t.Next == "" {
		return nil, nil
	}
	// Create request
	req, err := http.NewRequest("GET", t.Next, nil)
	if err != nil {
		return nil, err
	}

	// Add authorization header if API key is set
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Token "+c.apiKey)
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var telemetryResponse TelemetryResponse
	if err := json.NewDecoder(resp.Body).Decode(&telemetryResponse); err != nil {
		return nil, err
	}
	return &telemetryResponse, nil
}

func (c *Client) GetTelemetryResponsePrevPage(t *TelemetryResponse) (*TelemetryResponse, error) {
	if t.Prev == "" {
		return nil, nil
	}
	// Create request
	req, err := http.NewRequest("GET", t.Prev, nil)
	if err != nil {
		return nil, err
	}

	// Add authorization header if API key is set
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Token "+c.apiKey)
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var telemetryResponse TelemetryResponse
	if err := json.NewDecoder(resp.Body).Decode(&telemetryResponse); err != nil {
		return nil, err
	}
	return &telemetryResponse, nil
}
