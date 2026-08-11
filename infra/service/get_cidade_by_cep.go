package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/AndreAfonsoLana/go-degree/infra/service/dto"
)

type CEPClient struct {
	baseURL    string
	httpClient *http.Client
}

type CidadeResult struct {
	Cidade string
	Err    error
}

func NewGetCidadeService(baseURL string) *CEPClient /* retorno declaro acima */ {
	return &CEPClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}
func (v *CEPClient) GetCidadeByCEP(cep string) <-chan CidadeResult {
	resultChan := make(chan CidadeResult, 1)

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel() // Garante que o contexto será cancelado após o uso
		defer close(resultChan)

		base := v.baseURL
		if base == "" {
			base = "https://viacep.com.br" // Fallback de segurança
		}

		url := fmt.Sprintf("%s/ws/%s/json/", base, cep)

		request, error := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if error != nil {
			resultChan <- CidadeResult{Err: error}
			return
		}

		response, error := v.httpClient.Do(request)
		if error != nil {
			resultChan <- CidadeResult{Err: error}
			return
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			resultChan <- CidadeResult{Err: fmt.Errorf("erro na API do ViaCEP: status code %d", response.StatusCode)}
			return
		}

		var payload dto.CepResponseDTO
		if error := json.NewDecoder(response.Body).Decode(&payload); error != nil {
			resultChan <- CidadeResult{Err: error}
			return
		}

		if payload.Localidade == "" {
			resultChan <- CidadeResult{Err: fmt.Errorf("can not find zipcode")}
			return
		}

		resultChan <- CidadeResult{Cidade: payload.Localidade, Err: nil}
	}()
	return resultChan
}
