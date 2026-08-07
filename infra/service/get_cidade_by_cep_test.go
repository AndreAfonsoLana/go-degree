package service

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func NewCEPClientForTest(baseURL string, client *http.Client) *CEPClient {
	return &CEPClient{
		baseURL:    baseURL,
		httpClient: client,
	}
}

func TestGetCidadeByCEP_Sucess(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"cep": "13183000", "localidade": "Hortolândia", "uf": "SP"}`))
	}))
	defer mockServer.Close()

	cliente := NewCEPClientForTest(mockServer.URL, mockServer.Client())

	resultoChan := cliente.GetCidadeByCEP("13180000")

	select {

	case resposta := <-resultoChan:
		fmt.Println(resposta.Cidade, " // teste")

		if resposta.Err != nil {
			t.Fatalf("esperava err nil, mas recebeu o erro: %v", resposta.Err)
		}

		if resposta.Cidade != "Hortolândia" {
			t.Errorf("esperava cidade 'Hortolândia', recebeu: %s", resposta.Cidade)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout: a goroutine demorou muito para responder no channel")
	}

}
func TestFuncGetCidadeByCEP_NotFound(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json") // Resposta em JSON
		w.WriteHeader(http.StatusOK)                       // Retorna o status HTTP 200 OK
		w.Write([]byte(`{"erro": "true"}`))                // Retorna o corpo (payload)
	}))
	defer mockServer.Close() // Gararante que o servidor falso criado na memória será desligado

	cliente := NewCEPClientForTest(mockServer.URL, mockServer.Client())
	resultadoChan := cliente.GetCidadeByCEP("")

	resposta := <-resultadoChan
	if resposta.Err == nil {
		t.Fatal("esperava erro de CEP não encontrado, mas recebeu nil")
	}
	if resposta.Err.Error() != "não foi possível encontrar esse CEP" {
		t.Errorf("mensagem de erro incorreta: %v", resposta.Err)
	}

}
func TestFuncGetCidadeByCEP_HTTPError(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError) // Retorna o status HTTP 200 OK
	}))
	defer mockServer.Close() // Gararante que o servidor falso criado na memória será desligado

	client := NewCEPClientForTest(mockServer.URL, mockServer.Client())
	resultChan := client.GetCidadeByCEP("13183000")

	resposta := <-resultChan
	if resposta.Err == nil {
		t.Fatal("esperava erro de HTTP 500, recebeu nil")
	}
}
