package handlers

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kameikay/service-input/internal/service"
	mock "github.com/kameikay/service-input/internal/service/mocks"
	"github.com/kameikay/service-input/internal/usecase"
	"github.com/kameikay/service-input/pkg/exceptions"
	"github.com/kameikay/service-input/pkg/utils"
	"github.com/stretchr/testify/suite"
)

type HandlerSuite struct {
	suite.Suite
	ctrl                  *gomock.Controller
	getTemperatureService *mock.MockGetTemperatureServiceInterface
	ctx                   context.Context
}

func TestHandlerStart(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}

func (suite *HandlerSuite) HandlerSuiteDown() {
	defer suite.ctrl.Finish()
}

func (suite *HandlerSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.getTemperatureService = mock.NewMockGetTemperatureServiceInterface(suite.ctrl)
	suite.ctx = context.Background()
}

func (suite *HandlerSuite) TestNewHandler() {
	handler := NewHandler(suite.getTemperatureService)
	suite.NotNil(handler)
}

func (suite *HandlerSuite) TestGetTemperatures() {
	testCases := []struct {
		name             string
		expectations     func(getTemperatureService *mock.MockGetTemperatureServiceInterface)
		expectedResponse utils.ResponseDTO
		requestJson      string
	}{
		{
			name: "should return correct temperatures",
			expectations: func(getTemperatureService *mock.MockGetTemperatureServiceInterface) {
				getTemperatureService.EXPECT().GetTemperatureService(gomock.Any(), "12345678").Return(service.GetTemperatureServiceResponse{
					Success: true,
					Message: "success",
					Data: service.DataResponse{
						City:  "city",
						TempC: 20,
						TempF: 68,
						TempK: 293,
					},
				}, nil)
			},
			expectedResponse: utils.ResponseDTO{
				StatusCode: http.StatusOK,
				Message:    http.StatusText(http.StatusOK),
				Success:    true,
				Data:       usecase.Response{City: "city", TempC: 20, TempF: 68, TempK: 293},
			},
			requestJson: `{"cep":"12345678"}`,
		},
		{
			name: "should return error when cep length is different from 8",
			expectations: func(getTemperatureService *mock.MockGetTemperatureServiceInterface) {
				getTemperatureService.EXPECT().GetTemperatureService(gomock.Any(), "123451s").Times(0)
			},
			expectedResponse: utils.ResponseDTO{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    exceptions.ErrInvalidCEP.Error(),
				Success:    false,
			},
			requestJson: `{"cep":"123451s"}`,
		},
		{
			name: "should return error when cep is invalid",
			expectations: func(getTemperatureService *mock.MockGetTemperatureServiceInterface) {
				getTemperatureService.EXPECT().GetTemperatureService(gomock.Any(), "123451s").Times(0)
			},
			expectedResponse: utils.ResponseDTO{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    exceptions.ErrInvalidCEP.Error(),
				Success:    false,
			},
			requestJson: `{"cep":"1234567a"}`,
		},
		{
			name: "should return error when there is an error getting data from services",
			expectations: func(getTemperatureService *mock.MockGetTemperatureServiceInterface) {
				getTemperatureService.EXPECT().GetTemperatureService(gomock.Any(), "12345678").Return(service.GetTemperatureServiceResponse{}, errors.New("error"))
			},
			expectedResponse: utils.ResponseDTO{
				StatusCode: http.StatusBadRequest,
				Message:    "error",
				Success:    false,
			},
			requestJson: `{"cep":"12345678"}`,
		},
		{
			name: "should return error when cep is not found",
			expectations: func(getTemperatureService *mock.MockGetTemperatureServiceInterface) {
				getTemperatureService.EXPECT().GetTemperatureService(gomock.Any(), "12345678").Return(service.GetTemperatureServiceResponse{}, exceptions.ErrCannotFindZipcode)
			},
			expectedResponse: utils.ResponseDTO{
				StatusCode: http.StatusNotFound,
				Message:    exceptions.ErrCannotFindZipcode.Error(),
				Success:    false,
			},
			requestJson: `{"cep":"12345678"}`,
		},
		{
			name: "should return error when request is invalid",
			expectations: func(getTemperatureService *mock.MockGetTemperatureServiceInterface) {
				getTemperatureService.EXPECT().GetTemperatureService(gomock.Any(), "").Times(0)
			},
			expectedResponse: utils.ResponseDTO{
				StatusCode: http.StatusBadRequest,
				Message:    "error",
				Success:    false,
			},
			requestJson: `{"cep":123}`,
		},
		{
			name: "should return error when cep is invalid",
			expectations: func(getTemperatureService *mock.MockGetTemperatureServiceInterface) {
				getTemperatureService.EXPECT().GetTemperatureService(gomock.Any(), "12345678").Return(service.GetTemperatureServiceResponse{}, exceptions.ErrInvalidCEP)
			},
			expectedResponse: utils.ResponseDTO{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    exceptions.ErrInvalidCEP.Error(),
				Success:    false,
			},
			requestJson: `{"cep":"12345678"}`,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			tc.expectations(suite.getTemperatureService)
			request := httptest.NewRequest(http.MethodPost, "http://test/", nil)
			request = request.WithContext(suite.ctx)
			request.Header.Set("Content-Type", "application/json")
			request.Body = io.NopCloser(strings.NewReader(tc.requestJson))
			recorder := httptest.NewRecorder()

			handler := NewHandler(suite.getTemperatureService)
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
