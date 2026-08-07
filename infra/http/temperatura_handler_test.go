package http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AndreAfonsoLana/go-degree/infra/service"
	"github.com/AndreAfonsoLana/go-degree/internal/usecase"
)

type mockCEPProvedor struct {
	Result service.CidadeResult
}

func (m *mockCEPProvedor) GetCidadeByCEP(cep string) <-chan service.CidadeResult {
	ch := make(chan service.CidadeResult, 1)
	ch <- m.Result
	close(ch)
	return ch
}

type mockTempProvedor struct {
	Result service.TemperaturaResultado
}

func (m *mockTempProvedor) GetTemperaturaByCidade(cidade string) <-chan service.TemperaturaResultado {
	ch := make(chan service.TemperaturaResultado, 1)
	ch <- m.Result
	close(ch)
	return ch
}
func TestHandleTemperatura_InvalidZipcode(t *testing.T) {
	uc := usecase.NewConsultaClimaUseCase(&mockCEPProvedor{}, &mockTempProvedor{})
	handler := NewTemperaturaHandler(uc)

	req := httptest.NewRequest(http.MethodGet, "/temperatura?cep=123", nil)
	rr := httptest.NewRecorder() // rr atua como o http.ResponseWriter

	handler.HandleTemperatura(rr, req)

	if rr.Code != http.StatusUnprocessableEntity {
		t.Errorf("Esperado status %d, recebido %d", http.StatusUnprocessableEntity, rr.Code)
	}

	if rr.Body.String() != "invalid zipcode" {
		t.Errorf("Esperado body 'invalid zipcode', recebido '%s'", rr.Body.String())
	}
}

func TestHandleTemperatura_NotFound(t *testing.T) {
	// Setup simulando erro no provedor de CEP
	mockCEP := &mockCEPProvedor{
		Result: service.CidadeResult{Err: errors.New("cep não encontrado")},
	}
	uc := usecase.NewConsultaClimaUseCase(mockCEP, &mockTempProvedor{})
	handler := NewTemperaturaHandler(uc)

	req := httptest.NewRequest(http.MethodGet, "/temperatura?cep=99999999", nil)
	rr := httptest.NewRecorder()

	handler.HandleTemperatura(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Esperado status %d, recebido %d", http.StatusNotFound, rr.Code)
	}

	if rr.Body.String() != "can not find zipcode" {
		t.Errorf("Esperado body 'can not find zipcode', recebido '%s'", rr.Body.String())
	}
}

func TestHandleTemperatura_Success(t *testing.T) {
	// Setup simulando sucesso
	mockCEP := &mockCEPProvedor{
		Result: service.CidadeResult{Cidade: "São Paulo", Err: nil},
	}
	mockTemp := &mockTempProvedor{
		Result: service.TemperaturaResultado{Temp_c: 28.5, Temp_F: 83.3, Temp_K: 301.65, Err: nil},
	}
	uc := usecase.NewConsultaClimaUseCase(mockCEP, mockTemp)
	handler := NewTemperaturaHandler(uc)

	req := httptest.NewRequest(http.MethodGet, "/temperatura?cep=01153000", nil)
	rr := httptest.NewRecorder()

	handler.HandleTemperatura(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Esperado status %d, recebido %d", http.StatusOK, rr.Code)
	}
}
