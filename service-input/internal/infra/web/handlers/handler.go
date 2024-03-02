package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/kameikay/service-input/internal/service"
	"github.com/kameikay/service-input/internal/usecase"
	"github.com/kameikay/service-input/pkg/exceptions"
	"github.com/kameikay/service-input/pkg/utils"
)

type Handler struct {
	weatherApiService service.GetTemperatureServiceInterface
}

type InputDTO struct {
	Cep string `json:"cep"`
}

func NewHandler(weatherApiService service.GetTemperatureServiceInterface) *Handler {
	return &Handler{
		weatherApiService: weatherApiService,
	}
}

func (h *Handler) GetTemperatures(w http.ResponseWriter, r *http.Request) {
	var input InputDTO
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		utils.JsonResponse(w, utils.ResponseDTO{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Success:    false,
		})
		return
	}

	if !h.validateCEP(input.Cep) {
		utils.JsonResponse(w, utils.ResponseDTO{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    exceptions.ErrInvalidCEP.Error(),
			Success:    false,
		})
		return
	}

	getTemperaturesUseCase := usecase.NewGetTemperatureUseCase(h.weatherApiService)
	data, err := getTemperaturesUseCase.Execute(r.Context(), input.Cep)
	if err != nil {
		if err.Error() == exceptions.ErrInvalidCEP.Error() {
			utils.JsonResponse(w, utils.ResponseDTO{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    err.Error(),
				Success:    false,
			})
			return
		}

		if err.Error() == exceptions.ErrCannotFindZipcode.Error() {
			utils.JsonResponse(w, utils.ResponseDTO{
				StatusCode: http.StatusNotFound,
				Message:    err.Error(),
				Success:    false,
			})
			return
		}

		utils.JsonResponse(w, utils.ResponseDTO{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Success:    false,
		})
		return
	}

	utils.JsonResponse(w, utils.ResponseDTO{
		StatusCode: http.StatusOK,
		Message:    http.StatusText(http.StatusOK),
		Success:    true,
		Data:       data,
	})
}

func (h *Handler) validateCEP(cep string) bool {
	regex := regexp.MustCompile(`^\d{8}$`)

	if len(cep) != 8 {
		return false
	}

	if !regex.MatchString(cep) {
		return false
	}

	return true
}
