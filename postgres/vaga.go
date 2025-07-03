package postgres

import (
	"context"

	talenthub "github.com/caio86/talentHub"
	"github.com/caio86/talentHub/postgres/repository"
	"github.com/jackc/pgx/v5/pgtype"
)

var _ talenthub.VagaService = (*VagaService)(nil)

type VagaService struct {
	db   *DB
	repo *repository.Queries
}

func NewVagaService(db *DB) *VagaService {
	return &VagaService{db: db, repo: repository.New(db.conn)}
}

func (s *VagaService) FindVagaByID(ctx context.Context, id int) (*talenthub.Vaga, error) {
	result, err := s.repo.GetFullVacancyByID(ctx, int32(id))
	if err != nil {
		return nil, talenthub.Errorf(talenthub.ENOTFOUND, "vaga not found")
	}

	res := talenthub.Vaga{
		ID:          int(result.ID),
		Title:       result.Title,
		IsActive:    result.IsActive,
		Posted_date: result.PostedDate.Time,
		Benefits:    make([]string, 0), // Inicializar como array vazio
		Company:     "",                // Valor padrão
	}

	if result.Description != nil {
		res.Description = *result.Description
	}
	if result.Area != nil {
		res.Area = *result.Area
	}
	if result.Type != nil {
		res.Type = *result.Type
	}
	if result.Location != nil {
		res.Location = *result.Location
	}

	requirements, _ := s.repo.GetRequirementsByVacancyID(ctx, int32(id))
	if requirements == nil {
		res.Requirements = make([]string, 0)
	} else {
		res.Requirements = requirements
	}

	return &res, nil
}

func (s *VagaService) FindVagas(ctx context.Context, filter talenthub.VagaFilter) ([]*talenthub.Vaga, int, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback(ctx)

	var arg repository.ListVacanciesParams

	arg.Limit = filter.Limit
	arg.Offset = filter.Offset * arg.Limit

	if arg.Limit <= 0 {
		result, total, err := s.findAllVagas(ctx, tx)
		if err != nil {
			return nil, 0, talenthub.Errorf(talenthub.EINTERNAL, "internal server error: %s", err)
		}

		return result, total, nil

	} else {
		result, total, err := s.findVagas(ctx, tx, arg)
		if err != nil {
			return nil, 0, talenthub.Errorf(talenthub.EINTERNAL, "internal server error: %s", err)
		}

		return result, total, nil

	}
}

func (s *VagaService) findAllVagas(ctx context.Context, tx *Tx) ([]*talenthub.Vaga, int, error) {
	repoTx := s.repo.WithTx(tx)

	total, err := repoTx.CountVacancies(ctx)
	if err != nil {
		return nil, 0, err
	}

	result, err := repoTx.ListAllVacancies(ctx)
	if err != nil {
		return nil, 0, talenthub.Errorf(talenthub.EINTERNAL, "internal server error: %s", err)
	}

	res := make([]*talenthub.Vaga, len(result))
	for i, v := range result {
		res[i] = &talenthub.Vaga{
			ID:          int(v.ID),
			Title:       v.Title,
			IsActive:    v.IsActive,
			Posted_date: v.PostedDate.Time,
			Benefits:    make([]string, 0), // Inicializar como array vazio
			Company:     "",                // Valor padrão
		}

		if v.Description != nil {
			res[i].Description = *v.Description
		}
		if v.Area != nil {
			res[i].Area = *v.Area
		}
		if v.Type != nil {
			res[i].Type = *v.Type
		}
		if v.Location != nil {
			res[i].Location = *v.Location
		}

		requirements, _ := repoTx.GetRequirementsByVacancyID(ctx, int32(v.ID))
		if requirements == nil {
			res[i].Requirements = make([]string, 0)
		} else {
			res[i].Requirements = requirements
		}
	}

	return res, int(total), nil
}

func (s *VagaService) findVagas(ctx context.Context, tx *Tx, arg repository.ListVacanciesParams) ([]*talenthub.Vaga, int, error) {
	repoTx := s.repo.WithTx(tx)

	total, err := repoTx.CountVacancies(ctx)
	if err != nil {
		return nil, 0, err
	}

	result, err := repoTx.ListVacancies(ctx, arg)
	if err != nil {
		return nil, 0, talenthub.Errorf(talenthub.EINTERNAL, "internal server error: %s", err)
	}

	res := make([]*talenthub.Vaga, len(result))
	for i, v := range result {
		res[i] = &talenthub.Vaga{
			ID:          int(v.ID),
			Title:       v.Title,
			IsActive:    v.IsActive,
			Posted_date: v.PostedDate.Time,
			Benefits:    make([]string, 0), // Inicializar como array vazio
			Company:     "",                // Valor padrão
		}

		if v.Description != nil {
			res[i].Description = *v.Description
		}
		if v.Area != nil {
			res[i].Area = *v.Area
		}
		if v.Type != nil {
			res[i].Type = *v.Type
		}
		if v.Location != nil {
			res[i].Location = *v.Location
		}

		requirements, _ := repoTx.GetRequirementsByVacancyID(ctx, int32(v.ID))
		if requirements == nil {
			res[i].Requirements = make([]string, 0)
		} else {
			res[i].Requirements = requirements
		}
	}

	return res, int(total), nil
}

func (s *VagaService) CreateVaga(ctx context.Context, vaga *talenthub.Vaga) (*talenthub.Vaga, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	repoTx := s.repo.WithTx(tx)

	arg := repository.CreateVacancyParams{
		Title:       vaga.Title,
		Description: &vaga.Description,
		IsActive:    vaga.IsActive,
		Location:    &vaga.Location,
		PostedDate:  pgtype.Date{Time: tx.now, Valid: true},
	}

	if vaga.Type != "" {
		typeID, err := repoTx.GetTypeByName(ctx, vaga.Type)
		if err != nil {
			return nil, talenthub.Errorf(talenthub.EINTERNAL, "employment type does not exists")
		}
		arg.TypeID = &typeID.ID
	} else {
		arg.TypeID = nil
	}

	if vaga.Area != "" {
		areaID, err := repoTx.GetAreaByName(ctx, vaga.Area)
		if err != nil {
			return nil, talenthub.Errorf(talenthub.EINTERNAL, "employment area does not exists")
		}
		arg.AreaID = &areaID.ID
	} else {
		arg.AreaID = nil
	}

	newVaga, err := repoTx.CreateVacancy(ctx, arg)
	if err != nil {
		return nil, talenthub.Errorf(talenthub.EINTERNAL, "internal error: %s", err)
	}

	for _, v := range vaga.Requirements {
		req, err := repoTx.GetRequirementByName(ctx, v)
		if err != nil {
			req, err = repoTx.AddRequirement(ctx, v)
			if err != nil {
				return nil, talenthub.Errorf(talenthub.EINTERNAL, "internal error %s", err)
			}
		}
		err = repoTx.AddVacancyRequirement(ctx, repository.AddVacancyRequirementParams{
			VacancyID:     newVaga.ID,
			RequirementID: req.ID,
		})
		if err != nil {
			return nil, talenthub.Errorf(talenthub.EINTERNAL, "internal error %s", err)
		}
	}

	res := &talenthub.Vaga{
		ID:          int(newVaga.ID),
		Title:       vaga.Title,
		Description: vaga.Description,
		IsActive:    newVaga.IsActive,
		Area:        vaga.Area,
		Type:        vaga.Type,
		Location:    vaga.Location,
		Posted_date: newVaga.PostedDate.Time,
		Benefits:    vaga.Benefits,
		Salary:      vaga.Salary,
		Company:     vaga.Company,
	}

	res.Requirements = vaga.Requirements

	tx.Commit(ctx)

	return res, nil
}

func (s *VagaService) UpdateVaga(ctx context.Context, id int, upd talenthub.VagaUpdate) (*talenthub.Vaga, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	repoTx := s.repo.WithTx(tx)

	vaga, err := repoTx.GetVacancyByID(ctx, int32(id))
	if err != nil {
		return nil, talenthub.Errorf(talenthub.ENOTFOUND, "vaga not found")
	}

	arg := repository.UpdateVacancyParams{
		ID:          vaga.ID,
		Title:       vaga.Title,
		Description: vaga.Description,
		Location:    vaga.Location,
		AreaID:      vaga.AreaID,
		TypeID:      vaga.TypeID,
	}

	// Apply updates
	if upd.Title != nil {
		arg.Title = *upd.Title
	}
	if upd.Description != nil {
		arg.Description = upd.Description
	}
	if upd.Location != nil {
		arg.Location = upd.Location
	}
	if upd.Area != nil {
		areaID, err := repoTx.GetAreaByName(ctx, *upd.Area)
		if err != nil {
			return nil, talenthub.Errorf(talenthub.EINTERNAL, "employment area does not exists")
		}
		arg.AreaID = &areaID.ID
	}
	if upd.Type != nil {
		typeID, err := repoTx.GetTypeByName(ctx, *upd.Type)
		if err != nil {
			return nil, talenthub.Errorf(talenthub.EINTERNAL, "employment type does not exists")
		}
		arg.TypeID = &typeID.ID
	}

	updated, err := repoTx.UpdateVacancy(ctx, arg)
	if err != nil {
		return nil, err
	}

	res := &talenthub.Vaga{
		ID:          int(updated.ID),
		Title:       updated.Title,
		IsActive:    updated.IsActive,
		Posted_date: updated.PostedDate.Time,
		Benefits:    make([]string, 0), // Inicializar como array vazio
		Company:     "",                // Valor padrão
	}

	if updated.Description != nil {
		res.Description = *updated.Description
	}
	if updated.AreaID != nil {
		areaID, _ := repoTx.GetAreaByID(ctx, *updated.AreaID)
		res.Area = areaID.Name
	}
	if updated.TypeID != nil {
		typeID, _ := repoTx.GetTypeByID(ctx, *updated.TypeID)
		res.Type = typeID.Name
	}
	if updated.Location != nil {
		res.Location = *updated.Location
	}

	requirements, _ := repoTx.GetRequirementsByVacancyID(ctx, int32(id))
	if requirements == nil {
		res.Requirements = make([]string, 0)
	} else {
		res.Requirements = requirements
	}

	tx.Commit(ctx)

	return res, nil
}

func (s *VagaService) OpenVaga(ctx context.Context, id int) error {
	_, err := s.repo.GetVacancyByID(ctx, int32(id))
	if err != nil {
		return talenthub.Errorf(talenthub.ENOTFOUND, "vaga not found")
	}

	err = s.repo.OpenVacancy(ctx, int32(id))
	if err != nil {
		return talenthub.Errorf(talenthub.EINTERNAL, "internal error %s", err)
	}

	return nil
}

func (s *VagaService) CloseVaga(ctx context.Context, id int) error {
	_, err := s.repo.GetVacancyByID(ctx, int32(id))
	if err != nil {
		return talenthub.Errorf(talenthub.ENOTFOUND, "vaga not found")
	}

	err = s.repo.CloseVacancy(ctx, int32(id))
	if err != nil {
		return talenthub.Errorf(talenthub.EINTERNAL, "internal error %s", err)
	}

	return nil
}

func (s *VagaService) DeleteVaga(ctx context.Context, id int) error {
	_, err := s.repo.GetVacancyByID(ctx, int32(id))
	if err != nil {
		return talenthub.Errorf(talenthub.ENOTFOUND, "vaga not found")
	}

	err = s.repo.DeleteVacancy(ctx, int32(id))
	if err != nil {
		return talenthub.Errorf(talenthub.EINTERNAL, "internal error %s", err)
	}

	return nil
}
