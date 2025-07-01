package talenthub

import (
	"context"
	"time"
)

type Application struct {
	ID              int
	CandidateID     int
	VacancyID       int
	Score           int
	Status          string
	ApplicationDate time.Time
}

func (a *Application) Validate() error {
	if a.CandidateID <= 0 {
		return Errorf(EINVALID, "candidate id invalid")
	}
	if a.VacancyID <= 0 {
		return Errorf(EINVALID, "vacancy id invalid")
	}
	if a.Score < 0 {
		return Errorf(EINVALID, "score invalid")
	}
	if a.Status == "" {
		return Errorf(EINVALID, "status required")
	}

	return nil
}

type ApplicationService interface {
	FindApplicationByID(ctx context.Context, id int) (*Application, error)
	FindApplications(ctx context.Context, filter ApplicationFilter) ([]*Application, int, error)
	SearchApplicationsByCandidateID(ctx context.Context, candidateID int) ([]*Application, int, error)
	SearchApplicationsByVacancyID(ctx context.Context, vacancyID int) ([]*Application, int, error)
	RegisterApplication(ctx context.Context, application *Application) (*Application, error)
	UnregisterApplication(ctx context.Context, id int) error
	UpdateApplication(ctx context.Context, id int, upd ApplicationUpdate) (*Application, error)
}

type ApplicationFilter struct {
	Offset int32
	Limit  int32
}

type ApplicationUpdate struct {
	Score  *int    `json:"score"`
	Status *string `json:"status"`
}
