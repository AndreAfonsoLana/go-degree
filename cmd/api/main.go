package main

import (
	"fmt"
	"net/http"

	"os"

	httpinfra "github.com/AndreAfonsoLana/go-degree/infra/http"
	"github.com/AndreAfonsoLana/go-degree/infra/service"
	"github.com/AndreAfonsoLana/go-degree/internal/usecase"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Starting main...")
	errou := godotenv.Load("../../.env")

	if errou != nil {
		fmt.Printf("Erro ao carregar .env: %v\n", errou)
	}

	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Erro ao carregar .env: %v\n", err)
	}

	URL_TEMPERATURA := os.Getenv("URL_WEATHER")
	TEMPERATURA_API_KEY := os.Getenv("WEATHER_API_KEY")

	tempService := service.NewGetTemperaturaService(URL_TEMPERATURA, TEMPERATURA_API_KEY)

	consultaClimaUseCase := usecase.NewConsultaClimaUseCase(
		service.NewGetCidadeService(os.Getenv("URL_CEP")),
		tempService,
	)

	temperaturaHandler := httpinfra.NewTemperaturaHandler(consultaClimaUseCase)

	mux := http.NewServeMux()
	mux.HandleFunc("/temperatura", temperaturaHandler.HandleTemperatura)

	fmt.Println("Server is running on port 8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
	}

}
