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

func (c *Client) GetTelemetry(satelliteID string) ([]Telemetry, error) {
	resp, err := c.Get("/telemetry/", []urlParam{{"sat_id", satelliteID}, {"format", "json"}})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var telemetryResponse TelemetryResponse
	if err := json.NewDecoder(resp.Body).Decode(&telemetryResponse); err != nil {
		return nil, err
	}
	return telemetryResponse.Results, nil
}
