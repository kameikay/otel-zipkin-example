package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/spf13/viper"
)

type WeatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

type WeatherApiServiceInterface interface {
	GetWeatherData(ctx context.Context, location string) (*WeatherAPIResponse, error)
}

type WeatherApiService struct {
	client *http.Client
}

func NewWeatherApiService() *WeatherApiService {
	return &WeatherApiService{
		client: &http.Client{},
	}
}

func (s *WeatherApiService) GetWeatherData(ctx context.Context, location string) (*WeatherAPIResponse, error) {
	WEATHER_API_KEY := viper.GetString("WEATHER_API_KEY")
	urlString := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", WEATHER_API_KEY, url.QueryEscape(location))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlString, nil)
	if err != nil {
		return nil, err
	}

	res, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, errors.New("cannot find weather data")
	}

	var weatherAPIResponse WeatherAPIResponse
	err = json.NewDecoder(res.Body).Decode(&weatherAPIResponse)
	if err != nil {
		return nil, err
	}

	return &weatherAPIResponse, nil
}
