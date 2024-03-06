package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type DataResponse struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

type GetTemperatureServiceResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    DataResponse `json:"data,omitempty"`
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

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	res, err := s.client.Do(req)
	if err != nil {
		return GetTemperatureServiceResponse{}, err
	}

	defer res.Body.Close()

	var response GetTemperatureServiceResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return GetTemperatureServiceResponse{}, err
	}

	if !response.Success {
		return GetTemperatureServiceResponse{}, errors.New(response.Message)
	}

	return response, nil
}
