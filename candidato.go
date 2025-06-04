package talenthub

import (
	"context"
	"net/mail"
	"net/url"
)

type Candidato struct {
	ID    int    `json:"-"`
	Name  string `json:"name"`
	Email string `json:"email"`
	CPF   string `json:"cpf"`
	Phone string `json:"phone"`
	// Experiencas
	// Formação

	Vaga string `json:"-"`

	LinkCurriculo *url.URL `json:"-"`
}

func (c *Candidato) validate() error {
	if c.Name == "" {
		return Errorf(EINVALID, "name required")
	}
	if _, err := mail.ParseAddress(c.Email); err != nil {
		return Errorf(EINVALID, "email invalid")
	}
	if c.CPF == "" {
		return Errorf(EINVALID, "cpf required")
	}
	if c.Phone == "" {
		return Errorf(EINVALID, "phone required")
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
