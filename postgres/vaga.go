package postgres

import (
	"context"

	talenthub "github.com/caio86/talentHub"
	"github.com/caio86/talentHub/postgres/repository"
)

var _ talenthub.VagaService = (*VagaService)(nil)

type VagaService struct {
	repo *repository.Queries
}

func NewVagaService(db *DB) *VagaService {
	return &VagaService{repository.New(db.conn)}
}

func (s *VagaService) FindVagaByID(ctx context.Context, id int) (*talenthub.Vaga, error) {
	result, err := s.repo.GetVaga(ctx, int64(id))
	if err != nil {
		return nil, talenthub.Errorf(talenthub.ENOTFOUND, "vaga not found")
	}

	res := talenthub.Vaga{
		ID:          int(result.ID),
		Name:        result.Name,
		Description: result.Description,
		Open:        result.Open,
		CreatedAt:   result.CreatedAt,
		ExpiresAt:   result.ExpiresAt,
	}

	return &res, nil
}

func (s *VagaService) FindVagas(ctx context.Context, filter talenthub.VagaFilter) ([]*talenthub.Vaga, int, error) {
	var arg repository.ListVagasParams

	arg.Limit = int32(filter.Limit)
	arg.Offset = int32(filter.Offset)

	result, err := s.repo.ListVagas(ctx, arg)
	if err != nil {
		return nil, 0, talenthub.Errorf(talenthub.EINTERNAL, "internal server error: %s", err)
	}

	res := make([]*talenthub.Vaga, len(result))
	for i, v := range result {
		res[i] = &talenthub.Vaga{
			ID:          int(v.ID),
			Name:        v.Name,
			Description: v.Description,
			Open:        v.Open,
			CreatedAt:   v.CreatedAt,
			ExpiresAt:   v.ExpiresAt,
		}
	}

	return res, len(res), nil
}

func (s *VagaService) CreateVaga(ctx context.Context, vaga *talenthub.Vaga) error {
	arg := repository.CreateVagaParams{
		Name:        vaga.Name,
		Description: vaga.Description,
		Open:        vaga.Open,
		CreatedAt:   vaga.CreatedAt,
		ExpiresAt:   vaga.CreatedAt,
	}

	_, err := s.repo.CreateVaga(ctx, arg)
	if err != nil {
		return err
	}

	return nil
}

func (s *VagaService) UpdateVaga(ctx context.Context, id int, upd talenthub.VagaUpdate) (*talenthub.Vaga, error) {
	arg := repository.UpdateVagaParams{
		Name:        upd.Name,
		Description: upd.Description,
		Open:        upd.Open,
	}

	updated, err := s.repo.UpdateVaga(ctx, arg)
	if err != nil {
		return nil, err
	}

	res := &talenthub.Vaga{
		ID:          int(updated.ID),
		Name:        updated.Name,
		Description: updated.Description,
		Open:        updated.Open,
		CreatedAt:   updated.CreatedAt,
		ExpiresAt:   updated.CreatedAt,
	}

	return res, nil
}

func (s *VagaService) OpenVaga(ctx context.Context, id int) error {
	return talenthub.Errorf(talenthub.ENOTIMPLEMENTED, "not implemented")
}

func (s *VagaService) CloseVaga(ctx context.Context, id int) error {
	return talenthub.Errorf(talenthub.ENOTIMPLEMENTED, "not implemented")
}
