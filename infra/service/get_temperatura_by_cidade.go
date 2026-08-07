package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"context"

	"github.com/AndreAfonsoLana/go-degree/infra/service/dto"
	"github.com/AndreAfonsoLana/go-degree/utils"
)

type TemperaturaClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}
type TemperaturaResultado struct {
	Temp_c float64
	Temp_F float64
	Temp_K float64
	Err    error
}

func NewGetTemperaturaService(baseURL string, apiKey string) *TemperaturaClient {
	return &TemperaturaClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 90 * time.Second,
		},
	}
}
func (w *TemperaturaClient) GetTemperaturaByCidade(cidade string) <-chan TemperaturaResultado {

	resultChan := make(chan TemperaturaResultado, 1)

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		defer close(resultChan)
		cidadeFormata := utils.RemoverAcentos(cidade)

		fmt.Printf("Key: %s\n", w.apiKey)

		endpoint := fmt.Sprintf("%s/v1/current.json?key=%s&q=%s", w.baseURL, w.apiKey, cidadeFormata)
		fmt.Printf("Endpoint: %s\n", endpoint)

		request, error := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
		if error != nil {
			resultChan <- TemperaturaResultado{Err: error}
			fmt.Printf("Erro ao criar requisição: %v\n", error)
			return
		}
		resposta, error := w.httpClient.Do(request)

		if error != nil {
			resultChan <- TemperaturaResultado{Err: error}
			return
		}

		//bodyBytes, _ := io.ReadAll(resposta.Body)
		//fmt.Printf("<-> STATUS CODE HTTP: %d\n", resposta.StatusCode)
		//fmt.Printf("<-> RAW JSON DA API: %s\n", string(bodyBytes))

		var payLoad dto.TemperaturaResponseDTO
		if error := json.NewDecoder(resposta.Body).Decode(&payLoad); error != nil {
			resultChan <- TemperaturaResultado{Err: error}
			return
		}

		//fmt.Printf("<-> Payload: %+v\n", payLoad)

		resultChan <- TemperaturaResultado{
			Temp_c: payLoad.Current.Temp_c,
			Temp_F: utils.CelsiusParaFahrenheit(payLoad.Current.Temp_c),
			Temp_K: utils.CelsiusParaKelvin(payLoad.Current.Temp_c),
			Err:    nil,
		}

		defer resposta.Body.Close()
	}()
	return resultChan
}
