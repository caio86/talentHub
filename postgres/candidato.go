package postgres

import (
	"context"

	talenthub "github.com/caio86/talentHub"
)

var _ talenthub.CandidatoService = (*CandidatoService)(nil)

type CandidatoService struct{}

func NewCandidatoService() *CandidatoService {
	return &CandidatoService{}
}

func (s *CandidatoService) FindCandidatoByID(ctx context.Context, id int) (*talenthub.Candidato, error) {
	return nil, talenthub.Errorf(talenthub.ENOTIMPLEMENTED, "not implemented")
}

func (s *CandidatoService) FindCandidatos(ctx context.Context, filter talenthub.CandidatoFilter) ([]*talenthub.Candidato, int, error) {
	return nil, 0, talenthub.Errorf(talenthub.ENOTIMPLEMENTED, "not implemented")
}

func (s *CandidatoService) CreateCandidato(ctx context.Context, candidato *talenthub.Candidato) error {
	return talenthub.Errorf(talenthub.ENOTIMPLEMENTED, "not implemented")
}

func (s *CandidatoService) RegisterCandidato(ctx context.Context, candidatoID, vagaID int) error {
	return talenthub.Errorf(talenthub.ENOTIMPLEMENTED, "not implemented")
}

func (s *CandidatoService) UnregisterCandidato(ctx context.Context, candidatoID, vagaID int) error {
	return talenthub.Errorf(talenthub.ENOTIMPLEMENTED, "not implemented")
}

func (s *CandidatoService) UpdateCandidato(ctx context.Context, id int, upd talenthub.CandidatoUpdate) (*talenthub.Candidato, error) {
	return nil, talenthub.Errorf(talenthub.ENOTIMPLEMENTED, "not implemented")
}
