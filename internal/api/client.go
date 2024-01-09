package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/mick-io/duplo_go_cloud/internal/models"
)

type ForecastOptions struct {
	Latitude  string
	Longitude string
}

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

type WeatherAPIClient interface {
	GetForecast(opts ForecastOptions, result *models.Forecast) error
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
	}
}

func (c *Client) GetForecast(opts ForecastOptions, result *models.Forecast) error {
	reqURL, err := url.Parse(c.BaseURL + "/forecast")
	if err != nil {
		return err
	}

	if opts.Latitude == "" || opts.Longitude == "" {
		return errors.New("latitude and longitude are required")
	}

	params := url.Values{}
	params.Add("latitude", opts.Latitude)
	params.Add("longitude", opts.Longitude)
	params.Add("hourly", "temperature_2m")
	params.Add("temperature_unit", "fahrenheit")
	params.Add("wind_speed_unit", "mph")
	params.Add("timezone", "auto")
	reqURL.RawQuery = params.Encode()

	resp, err := c.HTTPClient.Get(reqURL.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(result)
}
