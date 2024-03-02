package handlers

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/kameikay/service-orchestration/internal/service"
	"github.com/kameikay/service-orchestration/internal/usecase"
	"github.com/kameikay/service-orchestration/pkg/exceptions"
	"github.com/kameikay/service-orchestration/pkg/utils"
)

type Handler struct {
	viaCepService     service.ViaCepServiceInterface
	weatherApiService service.WeatherApiServiceInterface
}

func NewHandler(
	viaCepService service.ViaCepServiceInterface,
	weatherApiService service.WeatherApiServiceInterface,
) *Handler {
	return &Handler{
		viaCepService:     viaCepService,
		weatherApiService: weatherApiService,
	}
}

func (h *Handler) GetTemperatures(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.JsonResponse(w, utils.ResponseDTO{
			StatusCode: http.StatusMethodNotAllowed,
			Message:    http.StatusText(http.StatusMethodNotAllowed),
			Success:    false,
		})
		return
	}

	cepParam := r.URL.Query().Get("cep")
	cep, err := h.formatCEP(cepParam)
	if err != nil {
		utils.JsonResponse(w, utils.ResponseDTO{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    err.Error(),
			Success:    false,
		})
		return
	}

	getTemperaturesUseCase := usecase.NewGetTemperatureUseCase(h.viaCepService, h.weatherApiService)
	data, err := getTemperaturesUseCase.Execute(r.Context(), cep)
	if err != nil {
		if err == exceptions.ErrCannotFindZipcode {
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

func (h *Handler) formatCEP(cep string) (string, error) {
	cepRegEx := `^\d{5}-\d{3}$`

	if regexp.MustCompile(cepRegEx).MatchString(cep) {
		return cep, nil
	}

	if len(cep) > 9 {
		return "", exceptions.ErrInvalidCEP
	}

	if len(cep) == 8 && !strings.Contains(cep, "-") {
		return cep[:5] + "-" + cep[5:], nil
	}

	return "", exceptions.ErrInvalidCEP
}
