package talenthub

import (
	"context"
	"time"
)

type Vaga struct {
	ID          int
	Title       string
	Description string

	IsActive bool

	Area         string
	Type         string
	Location     string
	Requirements []string
	Benefits     []string
	Salary       *string
	Company      string

	Posted_date time.Time
}

func (v *Vaga) Validate() error {
	if v.Title == "" {
		return Errorf(EINVALID, "title required")
	}
	if v.Description == "" {
		return Errorf(EINVALID, "description required")
	}

	return nil
}

type VagaService interface {
	FindVagaByID(ctx context.Context, id int) (*Vaga, error)
	FindVagas(ctx context.Context, filter VagaFilter) ([]*Vaga, int, error)
	CreateVaga(ctx context.Context, vaga *Vaga) (*Vaga, error)
	UpdateVaga(ctx context.Context, id int, upd VagaUpdate) (*Vaga, error)
	DeleteVaga(ctx context.Context, id int) error
	OpenVaga(ctx context.Context, id int) error
	CloseVaga(ctx context.Context, id int) error
}

type VagaFilter struct {
	Open *bool

	Offset int32
	Limit  int32
}

type VagaUpdate struct {
	Title       *string   `json:"title"`
	Description *string   `json:"description"`
	Area        *string   `json:"area"`
	Type        *string   `json:"type"`
	Location    *string   `json:"location"`
	Benefits    *[]string `json:"benefits"`
	Salary      *string   `json:"salary"`
	Company     *string   `json:"company"`
}
