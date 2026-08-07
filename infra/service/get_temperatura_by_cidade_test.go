package service

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func NewTemperaturaClientForTest(baseURL string, client *http.Client) *TemperaturaClient {
	return &TemperaturaClient{
		baseURL:    baseURL,
		apiKey:     "CHAVE_MOCK",
		httpClient: client,
	}
}

func TestGetTemperaturaByCidade_Success(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		w.Write([]byte(`{"current": {"temp_c": 28.5, "temp_f": 83.3}}`))
	}))
	defer mockServer.Close() // Desliga o servidor no fim do teste

	// 2. Instanciando o cliente usando a construtora de testes
	cliente := NewTemperaturaClientForTest(mockServer.URL, mockServer.Client())

	// 3. Disparando a função assíncrona
	resultChan := cliente.GetTemperaturaByCidade("Campinas")

	// 4. Lendo e validando o Canal
	select {
	case resposta := <-resultChan:
		if resposta.Err != nil {
			t.Fatalf("não esperava erro, mas recebeu: %v", resposta.Err)
		}

		fmt.Printf("Temp C: %.1f | Temp F: %.1f | Temp K: %.2f\n", resposta.Temp_c, resposta.Temp_F, resposta.Temp_K)

		if resposta.Temp_c != 28.5 {
			t.Errorf("esperava TempC 28.5, mas recebeu %.1f", resposta.Temp_c)
		}

		if resposta.Temp_F != 83.3 {
			t.Errorf("esperava TempF 83.3, mas recebeu %.1f", resposta.Temp_F)
		}

		if resposta.Temp_K != 301.65 {
			t.Errorf("esperava TempK 301.65, mas recebeu %.2f", resposta.Temp_K)
		}

	case <-time.After(2 * time.Second):
		t.Fatal("timeout: a goroutine demorou muito para responder")
	}
}
