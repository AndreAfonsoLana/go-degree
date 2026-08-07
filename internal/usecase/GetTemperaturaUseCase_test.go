package usecase

import (
	"fmt"
	"testing"

	"os"

	"github.com/AndreAfonsoLana/go-degree/infra/service"
	"github.com/AndreAfonsoLana/go-degree/internal/usecase/dto"
	"github.com/joho/godotenv"
)

func TestGetTemperatura(t *testing.T) {
	cepTeste := "13183310"

	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatalf("Erro ao carregar .env: %v", err)
	}

	URL_TEMPERATURA := os.Getenv("URL_WEATHER")
	URL_CEP := os.Getenv("URL_CEP")
	TEMPERATURA_API_KEY := os.Getenv("WEATHER_API_KEY")

	fmt.Printf("URL_WEATHER: %s\nURL_CEP: %s\n", URL_TEMPERATURA, URL_CEP)
	fmt.Printf("Key API %s\n", TEMPERATURA_API_KEY)

	cepService := service.NewGetCidadeService(URL_CEP)
	tempService := service.NewGetTemperaturaService(URL_TEMPERATURA, TEMPERATURA_API_KEY)

	usecaseNew := NewConsultaClimaUseCase(
		cepService,
		tempService,
	)

	resultado, erro := usecaseNew.ConsultarClimaPorCEP(cepTeste)

	if erro != nil {
		t.Errorf("Erro ao consultar clima: %v", erro)
	}
	if resultado == (dto.GetTemperaturaOutputDTO{}) {
		t.Errorf("Resultado vazio")
	}

	if resultado != (dto.GetTemperaturaOutputDTO{}) {
		fmt.Printf("<-> Resultado: %+v\n", resultado)
	}
}
