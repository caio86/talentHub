package postgres

import (
	"context"

	talenthub "github.com/caio86/talentHub"
	"github.com/caio86/talentHub/postgres/repository"
	"github.com/jackc/pgx/v5/pgtype"
)

var _ talenthub.VagaService = (*VagaService)(nil)

type VagaService struct {
	repo *repository.Queries
}

func NewVagaService(db *DB) *VagaService {
	return &VagaService{repository.New(db.conn)}
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
	var arg repository.ListVacanciesParams

	arg.Limit = filter.Limit
	arg.Offset = filter.Offset * arg.Limit

	if arg.Limit <= 0 {
		result, total, err := s.findAllVagas(ctx)
		if err != nil {
			return nil, 0, talenthub.Errorf(talenthub.EINTERNAL, "internal server error: %s", err)
		}

		return result, total, nil

	} else {
		result, total, err := s.findVagas(ctx, arg)
		if err != nil {
			return nil, 0, talenthub.Errorf(talenthub.EINTERNAL, "internal server error: %s", err)
		}

		return result, total, nil

	}
}

func (s *VagaService) findAllVagas(ctx context.Context) ([]*talenthub.Vaga, int, error) {
	total, err := s.repo.CountVacancies(ctx)
	if err != nil {
		return nil, 0, err
	}

	result, err := s.repo.ListAllVacancies(ctx)
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

		requirements, _ := s.repo.GetRequirementsByVacancyID(ctx, int32(v.ID))
		if requirements == nil {
			res[i].Requirements = make([]string, 0)
		} else {
			res[i].Requirements = requirements
		}
	}

	return res, int(total), nil
}

func (s *VagaService) findVagas(ctx context.Context, arg repository.ListVacanciesParams) ([]*talenthub.Vaga, int, error) {
	total, err := s.repo.CountVacancies(ctx)
	if err != nil {
		return nil, 0, err
	}

	result, err := s.repo.ListVacancies(ctx, arg)
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

		requirements, _ := s.repo.GetRequirementsByVacancyID(ctx, int32(v.ID))
		if requirements == nil {
			res[i].Requirements = make([]string, 0)
		} else {
			res[i].Requirements = requirements
		}
	}

	return res, int(total), nil
}

func (s *VagaService) CreateVaga(ctx context.Context, vaga *talenthub.Vaga) (*talenthub.Vaga, error) {
	arg := repository.CreateVacancyParams{
		Title:       vaga.Title,
		Description: &vaga.Description,
		IsActive:    vaga.IsActive,
		Location:    &vaga.Location,
		PostedDate:  pgtype.Date{Time: vaga.Posted_date, Valid: true},
	}

	if vaga.Type != "" {
		typeID, err := s.repo.GetTypeByName(ctx, vaga.Type)
		if err != nil {
			return nil, talenthub.Errorf(talenthub.EINTERNAL, "employment type does not exists")
		}
		arg.TypeID = &typeID.ID
	} else {
		arg.TypeID = nil
	}

	if vaga.Area != "" {
		areaID, err := s.repo.GetAreaByName(ctx, vaga.Area)
		if err != nil {
			return nil, talenthub.Errorf(talenthub.EINTERNAL, "employment area does not exists")
		}
		arg.AreaID = &areaID.ID
	} else {
		arg.AreaID = nil
	}

	newVaga, err := s.repo.CreateVacancy(ctx, arg)
	if err != nil {
		return nil, talenthub.Errorf(talenthub.EINTERNAL, "internal error: %s", err)
	}

	for _, v := range vaga.Requirements {
		req, err := s.repo.GetRequirementByName(ctx, v)
		if err != nil {
			req, err = s.repo.AddRequirement(ctx, v)
			if err != nil {
				return nil, talenthub.Errorf(talenthub.EINTERNAL, "internal error %s", err)
			}
		}
		err = s.repo.AddVacancyRequirement(ctx, repository.AddVacancyRequirementParams{
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
	}

	res.Requirements = vaga.Requirements

	return res, nil
}

func (s *VagaService) UpdateVaga(ctx context.Context, id int, upd talenthub.VagaUpdate) (*talenthub.Vaga, error) {
	vaga, err := s.repo.GetVacancyByID(ctx, int32(id))
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
		areaID, err := s.repo.GetAreaByName(ctx, *upd.Area)
		if err != nil {
			return nil, talenthub.Errorf(talenthub.EINTERNAL, "employment area does not exists")
		}
		arg.AreaID = &areaID.ID
	}
	if upd.Type != nil {
		typeID, err := s.repo.GetTypeByName(ctx, *upd.Type)
		if err != nil {
			return nil, talenthub.Errorf(talenthub.EINTERNAL, "employment type does not exists")
		}
		arg.TypeID = &typeID.ID
	}

	updated, err := s.repo.UpdateVacancy(ctx, arg)
	if err != nil {
		return nil, err
	}

	res := &talenthub.Vaga{
		ID:          int(updated.ID),
		Title:       updated.Title,
		IsActive:    updated.IsActive,
		Posted_date: updated.PostedDate.Time,
	}

	if updated.Description != nil {
		res.Description = *updated.Description
	}
	if updated.AreaID != nil {
		areaID, _ := s.repo.GetAreaByID(ctx, *updated.AreaID)
		res.Area = areaID.Name
	}
	if updated.TypeID != nil {
		typeID, _ := s.repo.GetTypeByID(ctx, *updated.TypeID)
		res.Type = typeID.Name
	}
	if updated.Location != nil {
		res.Location = *updated.Location
	}

	requirements, _ := s.repo.GetRequirementsByVacancyID(ctx, int32(id))
	if requirements == nil {
		res.Requirements = make([]string, 0)
	} else {
		res.Requirements = requirements
	}

	return res, nil
}

func (s *VagaService) OpenVaga(ctx context.Context, id int) error {
	return talenthub.Errorf(talenthub.ENOTIMPLEMENTED, "not implemented")
}

func (s *VagaService) CloseVaga(ctx context.Context, id int) error {
	return talenthub.Errorf(talenthub.ENOTIMPLEMENTED, "not implemented")
}
