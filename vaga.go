package talenthub

import (
	"context"
	"fmt"
	"time"
)

type Vaga struct {
	ID          int
	Name        string
	Description string

	Open bool
	// Localidade
	// Tipo
	// AreaAtuação

	CreatedAt time.Time
	ExpiredAt time.Time
}

func (v *Vaga) Validate() error {
	if v.Name == "" {
		return fmt.Errorf("name required")
	}
	if v.Description == "" {
		return fmt.Errorf("description required")
	}

	return nil
}

type VagaService interface {
	FindVagaByID(ctx context.Context, id int) (*Vaga, error)
	FindVagas(ctx context.Context, filter VagaFilter) ([]*Vaga, int, error)
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
