package postgres

import (
	"context"

	talenthub "github.com/caio86/talentHub"
)

var _ talenthub.VagaService = (*VagaService)(nil)

type VagaService struct{}

func NewVagaService() *VagaService {
	return &VagaService{}
}

func (s *VagaService) FindVagaByID(ctx context.Context, id int) (*talenthub.Vaga, error) {
	return nil, talenthub.Errorf(talenthub.ENOTIMPLEMENTED, "not implemented")
}

func (s *VagaService) FindVagas(ctx context.Context, filter talenthub.VagaFilter) ([]*talenthub.Vaga, int, error) {
	return nil, 0, talenthub.Errorf(talenthub.ENOTIMPLEMENTED, "not implemented")
}

func (s *VagaService) CreateVaga(ctx context.Context, vaga *talenthub.Vaga) error {
	return talenthub.Errorf(talenthub.ENOTIMPLEMENTED, "not implemented")
}

func (s *VagaService) UpdateVaga(ctx context.Context, id int, upd talenthub.VagaUpdate) (*talenthub.Vaga, error) {
	return nil, talenthub.Errorf(talenthub.ENOTIMPLEMENTED, "not implemented")
}

func (s *VagaService) OpenVaga(ctx context.Context, id int) error {
	return talenthub.Errorf(talenthub.ENOTIMPLEMENTED, "not implemented")
}

func (s *VagaService) CloseVaga(ctx context.Context, id int) error {
	return talenthub.Errorf(talenthub.ENOTIMPLEMENTED, "not implemented")
}
