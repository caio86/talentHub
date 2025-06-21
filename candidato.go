package talenthub

import (
	"context"
	"net/mail"
)

type Candidato struct {
	ID       int
	Name     string
	Email    string
	Password string
	Phone    string
	Address  string
	Linkedin string

	Experiences []*Experience
	Education   []*Education
	Skills      []string
	Interests   []string

	ResumeLink string
}

func (c *Candidato) validate() error {
	if c.Name == "" {
		return Errorf(EINVALID, "name required")
	}
	if _, err := mail.ParseAddress(c.Email); err != nil {
		return Errorf(EINVALID, "email invalid")
	}
	if c.Password == "" {
		return Errorf(EINVALID, "password required")
	}

	return nil
}

type CandidatoService interface {
	FindCandidatoByID(ctx context.Context, id int) (*Candidato, error)
	FindCandidatos(ctx context.Context, filter CandidatoFilter) ([]*Candidato, int, error)
	CreateCandidato(ctx context.Context, candidato *Candidato) (*Candidato, error)
	RegisterCandidato(ctx context.Context, candidatoID, vagaID int) error
	UnregisterCandidato(ctx context.Context, candidatoID, vagaID int) error
	UpdateCandidato(ctx context.Context, id int, upd CandidatoUpdate) (*Candidato, error)
}

type CandidatoFilter struct {
	Offset int32
	Limit  int32
}

type CandidatoUpdate struct {
	Name     *string `json:"name"`
	Phone    *string `json:"phone"`
	Address  *string `json:"address"`
	Linkedin *string `json:"linkedin"`

	ResumeLink *string `json:"resume_pdf_path"`
}
