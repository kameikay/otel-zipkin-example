package service

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/kameikay/service-orchestration/pkg/exceptions"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
)

type ViaCEPResponse struct {
	Erro        string `json:"erro"`
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type ViaCepServiceInterface interface {
	GetCEPData(ctx context.Context, cep string) (*ViaCEPResponse, error)
}

type ViaCepService struct {
	client *http.Client
}

func NewViaCepService() *ViaCepService {
	return &ViaCepService{client: &http.Client{}}
}

func (s *ViaCepService) GetCEPData(ctx context.Context, cep string) (*ViaCEPResponse, error) {
	tracer := otel.Tracer(viper.GetString("SERVICE_NAME"))
	ctx, span := tracer.Start(ctx, "ViaCEPService.GetCEPData")
	defer span.End()

	url := "http://viacep.com.br/ws/" + cep + "/json"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		var viaCEPResponse ViaCEPResponse
		err = json.NewDecoder(res.Body).Decode(&viaCEPResponse)
		if err != nil {
			return nil, err
		}

		if viaCEPResponse.Erro == "true" {
			return nil, exceptions.ErrCannotFindZipcode
		}

		return &viaCEPResponse, nil
	}

	return nil, nil
}
