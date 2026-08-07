package entity

import "errors"

type CEP struct {
	Value string
}

func NewCEP(value string) (*CEP, error) {
	if len(value) != 8 {
		return nil, errors.New("CEP inválido")
	}
	return &CEP{Value: value}, nil
}
