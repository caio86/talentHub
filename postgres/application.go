package postgres

import (
	"context"

	talenthub "github.com/caio86/talentHub"
	"github.com/caio86/talentHub/postgres/repository"
)

var _ talenthub.ApplicationService = (*ApplicationService)(nil)

type ApplicationService struct {
	db   *DB
	repo *repository.Queries
}

func NewApplicationService(db *DB) *ApplicationService {
	return &ApplicationService{
		db:   db,
		repo: repository.New(db.conn),
	}
}

func (s *ApplicationService) FindApplicationByID(ctx context.Context, id int) (*talenthub.Application, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	repoTx := s.repo.WithTx(tx)

	result, err := repoTx.GetFullApplicationById(ctx, int32(id))
	if err != nil {
		return nil, talenthub.Errorf(talenthub.ENOTFOUND, "application not found")
	}

	res := &talenthub.Application{
		ID:              int(result.ID),
		CandidateID:     int(result.CandidateID),
		VacancyID:       int(result.VacancyID),
		Score:           int(result.Score),
		ApplicationDate: result.ApplicationDate,
	}

	if result.Status != nil {
		res.Status = *result.Status
	}

	return res, nil
}

func (s *ApplicationService) FindApplications(ctx context.Context, filter talenthub.ApplicationFilter) ([]*talenthub.Application, int, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback(ctx)

	var arg repository.ListFullApplicationsParams

	arg.Limit = filter.Limit
	arg.Offset = filter.Offset * arg.Limit

	if arg.Limit <= 0 {
		result, total, err := s.findAllApplications(ctx, tx)
		if err != nil {
			return nil, 0, err
		}

		return result, total, nil

	} else {
		result, total, err := s.findApplications(ctx, tx, arg)
		if err != nil {
			return nil, 0, err
		}

		return result, total, nil

	}
}

func (s *ApplicationService) findAllApplications(ctx context.Context, tx *Tx) ([]*talenthub.Application, int, error) {
	repoTx := s.repo.WithTx(tx)

	total, err := repoTx.CountApplications(ctx)
	if err != nil {
		return nil, 0, err
	}

	result, err := repoTx.ListAllFullApplications(ctx)
	if err != nil {
		return nil, 0, talenthub.Errorf(talenthub.EINTERNAL, "internal server error: %s", err)
	}

	res := make([]*talenthub.Application, len(result))
	for i, v := range result {
		res[i] = &talenthub.Application{
			ID:              int(v.ID),
			CandidateID:     int(v.CandidateID),
			VacancyID:       int(v.VacancyID),
			Score:           int(v.Score),
			ApplicationDate: v.ApplicationDate,
		}

		if v.Status != nil {
			res[i].Status = *v.Status
		}
	}

	return res, int(total), nil
}

func (s *ApplicationService) findApplications(ctx context.Context, tx *Tx, arg repository.ListFullApplicationsParams) ([]*talenthub.Application, int, error) {
	repoTx := s.repo.WithTx(tx)

	total, err := repoTx.CountApplications(ctx)
	if err != nil {
		return nil, 0, err
	}

	result, err := repoTx.ListFullApplications(ctx, arg)
	if err != nil {
		return nil, 0, talenthub.Errorf(talenthub.EINTERNAL, "internal server error: %s", err)
	}

	res := make([]*talenthub.Application, len(result))
	for i, v := range result {
		res[i] = &talenthub.Application{
			ID:              int(v.ID),
			CandidateID:     int(v.CandidateID),
			VacancyID:       int(v.VacancyID),
			Score:           int(v.Score),
			ApplicationDate: v.ApplicationDate,
		}

		if v.Status != nil {
			res[i].Status = *v.Status
		}
	}

	return res, int(total), nil
}

func (s *ApplicationService) SearchApplicationsByCandidateID(ctx context.Context, candidateID int) ([]*talenthub.Application, int, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback(ctx)

	repoTx := s.repo.WithTx(tx)

	total, err := repoTx.CountApplications(ctx)
	if err != nil {
		return nil, 0, err
	}

	result, err := repoTx.SearchApplicationsByCandidateId(ctx, int32(candidateID))
	if err != nil {
		return nil, 0, talenthub.Errorf(talenthub.EINTERNAL, "internal server error: %s", err)
	}

	res := make([]*talenthub.Application, len(result))
	for i, v := range result {
		res[i] = &talenthub.Application{
			ID:              int(v.ID),
			CandidateID:     int(v.CandidateID),
			VacancyID:       int(v.VacancyID),
			Score:           int(v.Score),
			ApplicationDate: v.ApplicationDate,
		}

		if v.StatusID != nil {
			status, err := repoTx.GetApplicationStatusById(ctx, *v.StatusID)
			if err != nil {
				return nil, 0, err
			}
			res[i].Status = status.Status
		}
	}

	return res, int(total), nil
}

func (s *ApplicationService) SearchApplicationsByVacancyID(ctx context.Context, vacancyID int) ([]*talenthub.Application, int, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback(ctx)

	repoTx := s.repo.WithTx(tx)

	total, err := repoTx.CountApplications(ctx)
	if err != nil {
		return nil, 0, err
	}

	result, err := repoTx.SearchApplicationsByVacancyId(ctx, int32(vacancyID))
	if err != nil {
		return nil, 0, talenthub.Errorf(talenthub.EINTERNAL, "internal server error: %s", err)
	}

	res := make([]*talenthub.Application, len(result))
	for i, v := range result {
		res[i] = &talenthub.Application{
			ID:              int(v.ID),
			CandidateID:     int(v.CandidateID),
			VacancyID:       int(v.VacancyID),
			Score:           int(v.Score),
			ApplicationDate: v.ApplicationDate,
		}

		if v.StatusID != nil {
			status, err := repoTx.GetApplicationStatusById(ctx, *v.StatusID)
			if err != nil {
				return nil, 0, err
			}
			res[i].Status = status.Status
		}
	}

	return res, int(total), nil
}

func (s *ApplicationService) RegisterApplication(ctx context.Context, application *talenthub.Application) (*talenthub.Application, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	repoTx := s.repo.WithTx(tx)

	recebidoStatus, err := repoTx.GetApplicationStatusByName(ctx, "Recebido")
	if err != nil {
		return nil, err
	}

	arg := repository.RegisterApplicationParams{
		CandidateID: int32(application.CandidateID),
		VacancyID:   int32(application.VacancyID),
		StatusID:    &recebidoStatus.ID,
	}

	newApp, err := repoTx.RegisterApplication(ctx, arg)
	if err != nil {
		return nil, err
	}

	res := &talenthub.Application{
		ID:          int(newApp.ID),
		CandidateID: int(newApp.CandidateID),
		VacancyID:   int(newApp.VacancyID),
		Score:       int(newApp.Score),
		Status:      recebidoStatus.Status,
	}

	tx.Commit(ctx)

	return res, nil
}

func (s *ApplicationService) UnregisterApplication(ctx context.Context, id int) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	repoTx := s.repo.WithTx(tx)

	_, err = repoTx.GetApplicationById(ctx, int32(id))
	if err != nil {
		return talenthub.Errorf(talenthub.ENOTFOUND, "application not found")
	}

	err = repoTx.UnregisterApplication(ctx, int32(id))
	if err != nil {
		return err
	}

	tx.Commit(ctx)

	return nil
}

func (s *ApplicationService) UpdateApplication(ctx context.Context, id int, upd talenthub.ApplicationUpdate) (*talenthub.Application, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	repoTx := s.repo.WithTx(tx)

	application, err := repoTx.GetApplicationById(ctx, int32(id))
	if err != nil {
		return nil, talenthub.Errorf(talenthub.ENOTFOUND, "application not found")
	}

	arg := repository.UpdateApplicationParams{
		ID:       application.ID,
		Score:    application.Score,
		StatusID: application.StatusID,
	}

	// Apply updates
	if upd.Score != nil {
		arg.Score = int32(*upd.Score)
	}

	if upd.Status != nil {
		status, err := repoTx.GetApplicationStatusByName(ctx, *upd.Status)
		if err != nil {
			return nil, err
		}

		arg.StatusID = &status.ID
	}

	updated, err := repoTx.UpdateApplication(ctx, arg)
	if err != nil {
		return nil, err
	}

	res := &talenthub.Application{
		ID:              int(updated.ID),
		CandidateID:     int(updated.CandidateID),
		VacancyID:       int(updated.VacancyID),
		Score:           int(updated.Score),
		ApplicationDate: updated.ApplicationDate,
	}

	if updated.StatusID != nil {
		status, err := repoTx.GetApplicationStatusById(ctx, *updated.StatusID)
		if err != nil {
			return nil, err
		}

		res.Status = status.Status
	}

	tx.Commit(ctx)

	return res, nil
}
