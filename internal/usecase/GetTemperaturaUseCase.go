package usecase

import (
	"errors"
	"fmt"

	"github.com/AndreAfonsoLana/go-degree/infra/service"
	"github.com/AndreAfonsoLana/go-degree/internal/usecase/dto"
)

type CEPProvedor interface {
	GetCidadeByCEP(cep string) <-chan service.CidadeResult
}
type TemperaturaProvedor interface {
	GetTemperaturaByCidade(cidade string) <-chan service.TemperaturaResultado
}
type ConsultaClimaUseCase struct {
	cepProvedor         CEPProvedor
	temperaturaProvedor TemperaturaProvedor
}

func NewConsultaClimaUseCase(
	cep CEPProvedor,
	temperatura TemperaturaProvedor,
) *ConsultaClimaUseCase {

	return &ConsultaClimaUseCase{
		cepProvedor:         cep,
		temperaturaProvedor: temperatura,
	}
}
func (c *ConsultaClimaUseCase) ConsultarClimaPorCEP(cep string) (dto.GetTemperaturaOutputDTO, error) {
	fmt.Printf("Cep %v\n", cep)
	if len(cep) != 8 {
		return dto.GetTemperaturaOutputDTO{}, errors.New("CEP inválido: Precisa de 8 dígitos")
	}

	cidadeResult := <-c.cepProvedor.GetCidadeByCEP(cep)
	if cidadeResult.Err != nil {
		return dto.GetTemperaturaOutputDTO{}, errors.New("Erro ao obter cidade pelo CEP: " + cidadeResult.Err.Error())
	}

	temperaturas := <-c.temperaturaProvedor.GetTemperaturaByCidade(cidadeResult.Cidade)
	fmt.Printf("<-> Temperaturas: %+v\n", temperaturas)
	return dto.GetTemperaturaOutputDTO{
		Cidade: cidadeResult.Cidade,
		TempC:  temperaturas.Temp_c,
		TempF:  temperaturas.Temp_F,
		TempK:  temperaturas.Temp_K,
	}, nil
}
