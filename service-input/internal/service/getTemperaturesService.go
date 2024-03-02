package service

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/spf13/viper"
)

type GetTemperatureServiceResponse struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

type GetTemperatureServiceInterface interface {
	GetTemperatureService(ctx context.Context, cep string) (GetTemperatureServiceResponse, error)
}

type GetTemperatureService struct {
	client *http.Client
}

func NewGetTemperatureService() *GetTemperatureService {
	return &GetTemperatureService{
		client: &http.Client{},
	}
}

func (s *GetTemperatureService) GetTemperatureService(ctx context.Context, cep string) (GetTemperatureServiceResponse, error) {
	WEATHER_SERVICE_URL := viper.GetString("WEATHER_SERVICE_URL")
	URL := WEATHER_SERVICE_URL + "?cep=" + cep

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, URL, nil)
	if err != nil {
		return GetTemperatureServiceResponse{}, err
	}

	res, err := s.client.Do(req)
	if err != nil {
		return GetTemperatureServiceResponse{}, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return GetTemperatureServiceResponse{}, err
	}

	var response GetTemperatureServiceResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return GetTemperatureServiceResponse{}, err
	}

	return response, nil
}
