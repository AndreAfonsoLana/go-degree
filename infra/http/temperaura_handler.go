package http

import (
	"encoding/json"
	"fmt"
	netHttp "net/http"

	"github.com/AndreAfonsoLana/go-degree/internal/usecase"
)

type TemperaturaHandler struct {
	UseCase *usecase.ConsultaClimaUseCase
}

func NewTemperaturaHandler(uc *usecase.ConsultaClimaUseCase) *TemperaturaHandler {
	return &TemperaturaHandler{
		UseCase: uc,
	}
}

func (h *TemperaturaHandler) HandleTemperatura(w netHttp.ResponseWriter, r *netHttp.Request) {
	w.Header().Set("Content-Type", "application/json")
	cep := r.URL.Query().Get("cep")

	if len(cep) != 8 {
		netHttp.Error(w, "invalid zipcode", netHttp.StatusUnprocessableEntity)
		//w.Write([]byte("invalid zipcode3"))
		return
	}

	output, err := h.UseCase.ConsultarClimaPorCEP(cep)

	fmt.Printf("Output: %+v\n", output)
	if err != nil {
		netHttp.Error(w, err.Error(), netHttp.StatusNotFound)
		//w.Write([]byte("invalid zipcode"))
		return
	}

	w.WriteHeader(netHttp.StatusOK)
	json.NewEncoder(w).Encode(output)
}
