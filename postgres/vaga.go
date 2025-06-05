package postgres

import (
	"context"

	talenthub "github.com/caio86/talentHub"
	"github.com/caio86/talentHub/postgres/repository"
	"github.com/jackc/pgx/v5"
)

var _ talenthub.VagaService = (*VagaService)(nil)

type VagaService struct {
	conn *pgx.Conn
}

func NewVagaService(conn *pgx.Conn) *VagaService {
	return &VagaService{conn: conn}
}

func (s *VagaService) FindVagaByID(ctx context.Context, id int) (*talenthub.Vaga, error) {
	repo := repository.New(s.conn)
	result, err := repo.GetVaga(ctx, int64(id))
	if err != nil {
		return nil, err
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
	repo := repository.New(s.conn)
	var arg repository.ListVagasParams

	arg.Limit = int32(filter.Limit)
	arg.Offset = int32(filter.Offset)

	result, err := repo.ListVagas(ctx, arg)
	if err != nil {
		return nil, 0, err
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
	repo := repository.New(s.conn)
	arg := repository.CreateVagaParams{
		Name:        vaga.Name,
		Description: vaga.Description,
		Open:        vaga.Open,
		CreatedAt:   vaga.CreatedAt,
		ExpiresAt:   vaga.CreatedAt,
	}

	_, err := repo.CreateVaga(ctx, arg)
	if err != nil {
		return err
	}

	return nil
}

func (s *VagaService) UpdateVaga(ctx context.Context, id int, upd talenthub.VagaUpdate) (*talenthub.Vaga, error) {
	repo := repository.New(s.conn)
	arg := repository.UpdateVagaParams{
		Name:        upd.Name,
		Description: upd.Description,
		Open:        upd.Open,
	}

	updated, err := repo.UpdateVaga(ctx, arg)
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
