package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kameikay/service-orchestration/internal/service"
	mock "github.com/kameikay/service-orchestration/internal/service/mocks"
	"github.com/stretchr/testify/suite"
)

type GetTemperaturesUseCaseSuite struct {
	suite.Suite
	ctrl              *gomock.Controller
	viaCepService     *mock.MockViaCepServiceInterface
	weatherApiService *mock.MockWeatherApiServiceInterface
	ctx               context.Context
}

func TestGetTemperaturesUseCaseStart(t *testing.T) {
	suite.Run(t, new(GetTemperaturesUseCaseSuite))
}

func (suite *GetTemperaturesUseCaseSuite) GetTemperaturesUseCaseSuiteDown() {
	defer suite.ctrl.Finish()
}

func (suite *GetTemperaturesUseCaseSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.viaCepService = mock.NewMockViaCepServiceInterface(suite.ctrl)
	suite.weatherApiService = mock.NewMockWeatherApiServiceInterface(suite.ctrl)
	suite.ctx = context.Background()
}

func (suite *GetTemperaturesUseCaseSuite) TestNewGetCEPDataUseCase() {
	useCase := NewGetTemperatureUseCase(suite.viaCepService, suite.weatherApiService)
	suite.NotNil(useCase)
}

func (suite *GetTemperaturesUseCaseSuite) TestExecute() {
	testCases := []struct {
		name         string
		cep          string
		expectations func(viaCepService *mock.MockViaCepServiceInterface, weatherApiService *mock.MockWeatherApiServiceInterface)
		expectedResp Response
		expectedErr  error
	}{
		{
			name: "should return correct temperatures",
			cep:  "12345678",
			expectations: func(viaCepService *mock.MockViaCepServiceInterface, weatherApiService *mock.MockWeatherApiServiceInterface) {
				viaCepService.EXPECT().GetCEPData(suite.ctx, "12345678").Return(&service.ViaCEPResponse{
					Localidade: "São Paulo",
				}, nil)
				weatherApiService.EXPECT().GetWeatherData(suite.ctx, "São Paulo").Return(&service.WeatherAPIResponse{
					Current: struct {
						TempC float64 `json:"temp_c"`
					}{
						TempC: 25,
					},
				}, nil)
			},
			expectedResp: Response{
				TempC: 25,
				TempF: 77,
				TempK: 298,
			},
			expectedErr: nil,
		},
		{
			name: "should return error when Via Cep Service returns error",
			cep:  "12345678",
			expectations: func(viaCepService *mock.MockViaCepServiceInterface, weatherApiService *mock.MockWeatherApiServiceInterface) {
				viaCepService.EXPECT().GetCEPData(suite.ctx, "12345678").Return(nil, errors.New("error"))
				weatherApiService.EXPECT().GetWeatherData(suite.ctx, "São Paulo").Times(0)
			},
			expectedResp: Response{},
			expectedErr:  errors.New("error"),
		},
		{
			name: "should return error when Weather API Service returns error",
			cep:  "12345678",
			expectations: func(viaCepService *mock.MockViaCepServiceInterface, weatherApiService *mock.MockWeatherApiServiceInterface) {
				viaCepService.EXPECT().GetCEPData(suite.ctx, "12345678").Return(&service.ViaCEPResponse{
					Localidade: "São Paulo",
				}, nil)
				weatherApiService.EXPECT().GetWeatherData(suite.ctx, "São Paulo").Return(nil, errors.New("error"))
			},
			expectedResp: Response{},
			expectedErr:  errors.New("error"),
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			tc.expectations(suite.viaCepService, suite.weatherApiService)
			useCase := NewGetTemperatureUseCase(suite.viaCepService, suite.weatherApiService)
			res, err := useCase.Execute(suite.ctx, tc.cep)
			suite.Equal(tc.expectedResp, res)
			suite.Equal(tc.expectedErr, err)
		})
	}

}
