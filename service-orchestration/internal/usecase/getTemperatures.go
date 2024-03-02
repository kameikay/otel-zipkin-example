package usecase

import (
	"context"

	"github.com/kameikay/service-orchestration/internal/service"
)

type GetTemperaturesUseCase struct {
	viaCepService     service.ViaCepServiceInterface
	weatherApiService service.WeatherApiServiceInterface
}

type Response struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func NewGetTemperatureUseCase(
	viaCepService service.ViaCepServiceInterface,
	weatherApiService service.WeatherApiServiceInterface,
) *GetTemperaturesUseCase {
	return &GetTemperaturesUseCase{
		viaCepService:     viaCepService,
		weatherApiService: weatherApiService,
	}
}

func (u *GetTemperaturesUseCase) Execute(ctx context.Context, cep string) (Response, error) {
	cepData, err := u.viaCepService.GetCEPData(ctx, cep)
	if err != nil {
		return Response{}, err
	}

	weatherData, err := u.weatherApiService.GetWeatherData(ctx, cepData.Localidade)
	if err != nil {
		return Response{}, err
	}

	tempF := weatherData.Current.TempC*1.8 + 32
	tempK := weatherData.Current.TempC + 273

	return Response{
		City:  cepData.Localidade,
		TempC: weatherData.Current.TempC,
		TempF: tempF,
		TempK: tempK,
	}, nil

}
