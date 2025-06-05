package talenthub

import (
	"context"
	"time"
)

type Vaga struct {
	ID          int    `json:"-"`
	Name        string `json:"name"`
	Description string `json:"description"`

	Open bool `json:"open"`
	// Localidade
	// Tipo
	// AreaAtuação

	CreatedAt time.Time `json:"-"`
	ExpiresAt time.Time `json:"-"`
}

func (v *Vaga) Validate() error {
	if v.Name == "" {
		return Errorf(EINVALID, "name required")
	}
	if v.Description == "" {
		return Errorf(EINVALID, "description required")
	}

	return nil
}

type VagaService interface {
	FindVagaByID(ctx context.Context, id int) (*Vaga, error)
	FindVagas(ctx context.Context, filter VagaFilter) ([]*Vaga, int, error)
	CreateVaga(ctx context.Context, vaga *Vaga) error
	UpdateVaga(ctx context.Context, id int, upd VagaUpdate) (*Vaga, error)
	OpenVaga(ctx context.Context, id int) error
	CloseVaga(ctx context.Context, id int) error
}

type VagaFilter struct {
	Open *bool

	Offset int
	Limit  int
}

type VagaUpdate struct {
	Name        string
	Description string
}
