package openweather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	apiKey     string
	httpClient *http.Client
}

type WeatherResponse struct {
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
	} `json:"main"`
	Weather []struct {
		Main string `json:"main"`
	} `json:"weather"`
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) GetWeather(city string) (*WeatherResponse, error) {
	url := fmt.Sprintf(
		"https://api.openweathermap.org/data/2.5/weather?q=%s,in&appid=%s",
		city,
		c.apiKey,
	)

	response, err := c.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var weather WeatherResponse
	if err := json.NewDecoder(response.Body).Decode(&weather); err != nil {
		return nil, err
	}

	return &weather, nil
}
