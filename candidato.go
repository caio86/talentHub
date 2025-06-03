package talenthub

import (
	"context"
	"fmt"
	"net/mail"
	"net/url"
)

type Candidato struct {
	ID    int
	Name  string
	Email string
	CPF   string
	Phone string
	// Experiencas
	// Formação

	Vaga string

	LinkCurriculo *url.URL
}

func (c *Candidato) validate() error {
	if c.Name == "" {
		return fmt.Errorf("name required")
	}
	if _, err := mail.ParseAddress(c.Email); err != nil {
		return fmt.Errorf("email invalid")
	}
	if c.CPF == "" {
		return fmt.Errorf("cpf required")
	}
	if c.Phone == "" {
		return fmt.Errorf("phone required")
	}

	return nil
}

type CandidatoService interface {
	FindCandidatoByID(ctx context.Context, id int) (*Candidato, error)
	FindCandidatos(ctx context.Context, filter CandidatoFilter) ([]*Candidato, int, error)
	RegisterCandidato(ctx context.Context, candidatoID, vagaID int) error
	UnregisterCandidato(ctx context.Context, candidatoID, vagaID int) error
	UpdateCandidato(ctx context.Context, id int, upd CandidatoUpdate) (*Candidato, error)
}

type CandidatoFilter struct {
	Offset int
	Limit  int
}

type CandidatoUpdate struct {
	Name  string
	Email string
	CPF   string
	Phone string

	LinkCurriculo *url.URL
}
