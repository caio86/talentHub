package postgres

import (
	"context"

	talenthub "github.com/caio86/talentHub"
	"github.com/caio86/talentHub/postgres/repository"
)

var _ talenthub.CandidatoService = (*CandidatoService)(nil)

type CandidatoService struct {
	repo *repository.Queries
}

func NewCandidatoService(db *DB) *CandidatoService {
	return &CandidatoService{repository.New(db.conn)}
}

func (s *CandidatoService) FindCandidatoByID(ctx context.Context, id int) (*talenthub.Candidato, error) {
	result, err := s.repo.GetCandidato(ctx, int64(id))
	if err != nil {
		return nil, talenthub.Errorf(talenthub.ENOTFOUND, "candidato not found")
	}

	res := talenthub.Candidato{
		ID:    int(result.ID),
		Name:  result.Name,
		Email: result.Email,
		CPF:   result.Cpf,
		Phone: result.Phone,
	}

	return &res, nil
}

func (s *CandidatoService) FindCandidatos(ctx context.Context, filter talenthub.CandidatoFilter) ([]*talenthub.Candidato, int, error) {
	arg := repository.ListCandidatosParams{
		Limit:  int32(filter.Limit),
		Offset: int32(filter.Offset),
	}

	candidatos, err := s.repo.ListCandidatos(ctx, arg)
	if err != nil {
		return nil, 0, err
	}

	res := make([]*talenthub.Candidato, len(candidatos))
	for i, v := range candidatos {
		res[i] = &talenthub.Candidato{
			ID:    int(v.ID),
			Name:  v.Name,
			Email: v.Email,
			CPF:   v.Cpf,
			Phone: v.Phone,
		}
	}

	return res, len(candidatos), nil
}

func (s *CandidatoService) CreateCandidato(ctx context.Context, candidato *talenthub.Candidato) error {
	arg := repository.CreateCandidatoParams{
		Name:  candidato.Name,
		Email: candidato.Email,
		Cpf:   candidato.CPF,
		Phone: candidato.Phone,
	}
	_, err := s.repo.CreateCandidato(ctx, arg)
	if err != nil {
		return talenthub.Errorf(talenthub.EINTERNAL, "internal error: %v", err)
	}

	return nil
}

func (s *CandidatoService) RegisterCandidato(ctx context.Context, candidatoID, vagaID int) error {
	return talenthub.Errorf(talenthub.ENOTIMPLEMENTED, "not implemented")
}

func (s *CandidatoService) UnregisterCandidato(ctx context.Context, candidatoID, vagaID int) error {
	return talenthub.Errorf(talenthub.ENOTIMPLEMENTED, "not implemented")
}

func (s *CandidatoService) UpdateCandidato(ctx context.Context, id int, upd talenthub.CandidatoUpdate) (*talenthub.Candidato, error) {
	arg := repository.UpdateCandidatoParams{
		ID:    int64(id),
		Name:  upd.Name,
		Email: upd.Email,
		Cpf:   upd.CPF,
		Phone: upd.Phone,
	}

	updated, err := s.repo.UpdateCandidato(ctx, arg)
	if err != nil {
		return nil, talenthub.Errorf(talenthub.EINTERNAL, "internal error: %v", err)
	}

	res := &talenthub.Candidato{
		ID:    int(updated.ID),
		Name:  updated.Name,
		Email: updated.Email,
		CPF:   updated.Cpf,
		Phone: updated.Phone,
	}

	return res, nil
}
