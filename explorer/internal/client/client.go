package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/MiyukiMori11/weatherstat/explorer/internal/config"
	"go.uber.org/zap"
)

type client struct {
	config *config.Client
	logger *zap.Logger

	http *http.Client
}

type Client interface {
	GetTemperature(Latitude, Longitude float64) (float64, error)
}

// TempResponse is a structure of response to weatherApi request
type TempResponse struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
}

type CoordinatesInfo struct {
	Latitude    float64 `json:"lat"`
	Longitude   float64 `json:"lon"`
	CountryCode string  `json:"country"`
}

// GeoResponse is a structure of response to geoApi request
type GeoResponse struct {
	CountriesList []CoordinatesInfo
}

// New is a client constructor
func New(cfg *config.Client, logger *zap.Logger) Client {
	return &client{
		config: cfg,
		logger: logger,
		http: &http.Client{
			Timeout: cfg.Timeout(),
		},
	}
}

// GetTemperature returns temperature in K
func (c *client) GetTemperature(latitude, longitude float64) (float64, error) {
	var tempInfo TempResponse

	resp, err := c.http.Get(fmt.Sprintf("%s?lat=%f&lon=%f&appid=%s", c.config.WeatherAPI, latitude, longitude, c.config.APIKey()))
	if err != nil {
		return 0, fmt.Errorf("can't receive response: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("can't read data from body: %w", err)
	}

	if err := json.Unmarshal(data, &tempInfo); err != nil {
		return 0, fmt.Errorf("can't unmarshal body: %w", err)
	}

	return tempInfo.Main.Temp, nil

}

// GetCoordinates returns coordinates of city
func (c *client) GetCoordinates(city, countryCode string) (float64, float64, error) {
	// http: //api.openweathermap.org/geo/1.0/direct?q=London&limit=5&appid={API key}
	var coordinatesInfo CoordinatesInfo
	resp, err := c.http.Get(fmt.Sprintf("%s?q=%s,%s&limit=1&appid=%s", c.config.GeoAPI, city, countryCode, c.config.APIKey()))
	if err != nil {
		return 0, 0, fmt.Errorf("can't receive response: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, fmt.Errorf("can't read data from body: %w", err)
	}

	if err := json.Unmarshal(data, &coordinatesInfo); err != nil {
		return 0, 0, fmt.Errorf("can't unmarshal body: %w", err)
	}

	return coordinatesInfo.Latitude, coordinatesInfo.Longitude, err

}
