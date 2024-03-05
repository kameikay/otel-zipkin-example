package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kameikay/service-orchestration/internal/service"
	mock "github.com/kameikay/service-orchestration/internal/service/mocks"
	"github.com/kameikay/service-orchestration/internal/usecase"
	"github.com/kameikay/service-orchestration/pkg/exceptions"
	"github.com/kameikay/service-orchestration/pkg/utils"
	"github.com/stretchr/testify/suite"
)

type HandlerSuite struct {
	suite.Suite
	ctrl              *gomock.Controller
	viaCepService     *mock.MockViaCepServiceInterface
	weatherApiService *mock.MockWeatherApiServiceInterface
	ctx               context.Context
}

func TestHandlerStart(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}

func (suite *HandlerSuite) HandlerSuiteDown() {
	defer suite.ctrl.Finish()
}

func (suite *HandlerSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.viaCepService = mock.NewMockViaCepServiceInterface(suite.ctrl)
	suite.weatherApiService = mock.NewMockWeatherApiServiceInterface(suite.ctrl)
	suite.ctx = context.Background()
}

func (suite *HandlerSuite) TestNewHandler() {
	handler := NewHandler(suite.viaCepService, suite.weatherApiService)
	suite.NotNil(handler)
}

func (suite *HandlerSuite) TestGetTemperatures() {
	testCases := []struct {
		name             string
		cep              string
		expectations     func(viaCepService *mock.MockViaCepServiceInterface, weatherApiService *mock.MockWeatherApiServiceInterface)
		expectedResponse utils.ResponseDTO
	}{
		{
			name: "should return correct temperatures",
			cep:  "12345-678",
			expectations: func(viaCepService *mock.MockViaCepServiceInterface, weatherApiService *mock.MockWeatherApiServiceInterface) {
				viaCepService.EXPECT().GetCEPData(gomock.Any(), "12345-678").Return(&service.ViaCEPResponse{
					Localidade: "São Paulo",
				}, nil)
				weatherApiService.EXPECT().GetWeatherData(gomock.Any(), "São Paulo").Return(&service.WeatherAPIResponse{
					Current: struct {
						TempC float64 `json:"temp_c"`
					}{
						TempC: 20,
					},
				}, nil)
			},
			expectedResponse: utils.ResponseDTO{
				StatusCode: http.StatusOK,
				Message:    http.StatusText(http.StatusOK),
				Success:    true,
				Data:       usecase.Response{TempC: 20, TempF: 68, TempK: 293},
			},
		},
		{
			name: "should return error when cep is invalid",
			cep:  "123451s",
			expectations: func(viaCepService *mock.MockViaCepServiceInterface, weatherApiService *mock.MockWeatherApiServiceInterface) {
				viaCepService.EXPECT().GetCEPData(gomock.Any(), "12345-678").Times(0)
			},
			expectedResponse: utils.ResponseDTO{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    exceptions.ErrInvalidCEP.Error(),
				Success:    false,
			},
		},
		{
			name: "should return error when there is an error getting data from services",
			cep:  "12345-678",
			expectations: func(viaCepService *mock.MockViaCepServiceInterface, weatherApiService *mock.MockWeatherApiServiceInterface) {
				viaCepService.EXPECT().GetCEPData(gomock.Any(), "12345-678").Return(nil, errors.New("error"))
			},
			expectedResponse: utils.ResponseDTO{
				StatusCode: http.StatusBadRequest,
				Message:    "error",
				Success:    false,
			},
		},
		{
			name: "should return error when cep is not found",
			cep:  "12345678",
			expectations: func(viaCepService *mock.MockViaCepServiceInterface, weatherApiService *mock.MockWeatherApiServiceInterface) {
				viaCepService.EXPECT().GetCEPData(gomock.Any(), "12345-678").Return(nil, exceptions.ErrCannotFindZipcode)
			},
			expectedResponse: utils.ResponseDTO{
				StatusCode: http.StatusNotFound,
				Message:    exceptions.ErrCannotFindZipcode.Error(),
				Success:    false,
			},
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			tc.expectations(suite.viaCepService, suite.weatherApiService)
			request := httptest.NewRequest(http.MethodGet, "http://test/?cep="+tc.cep, nil)
			recorder := httptest.NewRecorder()

			handler := NewHandler(suite.viaCepService, suite.weatherApiService)
			handler.GetTemperatures(recorder, request)

			suite.Equal(tc.expectedResponse, utils.ResponseDTO{
				StatusCode: recorder.Code,
				Message:    tc.expectedResponse.Message,
				Success:    tc.expectedResponse.Success,
				Data:       tc.expectedResponse.Data,
			})
		})
	}
}

func (suite *HandlerSuite) TestFormatCep() {
	ceps := []struct {
		cep           string
		expectedCep   string
		expectedError error
	}{
		{
			cep:           "12345678",
			expectedCep:   "12345-678",
			expectedError: nil,
		},
		{
			cep:           "12345-678",
			expectedCep:   "12345-678",
			expectedError: nil,
		},
		{
			cep:           "12345 678",
			expectedCep:   "",
			expectedError: exceptions.ErrInvalidCEP,
		},
		{
			cep:           "1234567",
			expectedCep:   "",
			expectedError: exceptions.ErrInvalidCEP,
		},
		{
			cep:           "0123456789",
			expectedCep:   "",
			expectedError: exceptions.ErrInvalidCEP,
		},
	}

	for _, tc := range ceps {
		suite.T().Run(tc.cep, func(t *testing.T) {
			handler := NewHandler(suite.viaCepService, suite.weatherApiService)
			cep, err := handler.formatCEP(tc.cep)
			suite.Equal(tc.expectedCep, cep)
			suite.Equal(tc.expectedError, err)
		})
	}
}
