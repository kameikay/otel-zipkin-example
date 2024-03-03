package usecase

import (
	"context"

	"github.com/kameikay/service-input/internal/service"
)

type GetTemperaturesUseCase struct {
	weatherApiService service.GetTemperatureServiceInterface
}

type Response struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func NewGetTemperatureUseCase(weatherApiService service.GetTemperatureServiceInterface) *GetTemperaturesUseCase {
	return &GetTemperaturesUseCase{
		weatherApiService: weatherApiService,
	}
}

func (u *GetTemperaturesUseCase) Execute(ctx context.Context, cep string) (Response, error) {
	weatherData, err := u.weatherApiService.GetTemperatureService(ctx, cep)
	if err != nil {
		return Response{}, err
	}

	return Response{
		City:  weatherData.Data.City,
		TempC: weatherData.Data.TempC,
		TempF: weatherData.Data.TempF,
		TempK: weatherData.Data.TempK,
	}, nil

}
