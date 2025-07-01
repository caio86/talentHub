package postgres

import (
	"context"
	"log"
	"sync"
	"time"

	talenthub "github.com/caio86/talentHub"
)

// Ensure ApplicationService implements the interface
var _ talenthub.ApplicationService = (*ApplicationService)(nil)

type ApplicationService struct {
	db *DB
	
	// In-memory storage
	mu           sync.RWMutex
	applications []*talenthub.Application
	nextID       int
}

func NewApplicationService(db *DB) *ApplicationService {
	// Initialize with some mock data
	initialApplications := []*talenthub.Application{
		{
			ID:              1,
			CandidateID:     1,
			VacancyID:       1,
			Score:           85,
			Status:          "pending",
			ApplicationDate: time.Now().AddDate(0, 0, -2),
		},
		{
			ID:              2,
			CandidateID:     2,
			VacancyID:       1,
			Score:           92,
			Status:          "approved",
			ApplicationDate: time.Now().AddDate(0, 0, -1),
		},
	}
	
	service := &ApplicationService{
		db:           db,
		applications: initialApplications,
		nextID:       3,
	}
	
	log.Printf("ApplicationService initialized with %d applications", len(initialApplications))
	return service
}

func (s *ApplicationService) FindApplications(ctx context.Context, filter talenthub.ApplicationFilter) ([]*talenthub.Application, int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	log.Printf("FindApplications called. Current applications count: %d", len(s.applications))
	
	// Return copy of applications to avoid race conditions
	result := make([]*talenthub.Application, len(s.applications))
	for i, app := range s.applications {
		appCopy := *app
		result[i] = &appCopy
		log.Printf("Application %d: ID=%d, CandidateID=%d, VacancyID=%d, Status=%s", 
			i, app.ID, app.CandidateID, app.VacancyID, app.Status)
	}
	
	log.Printf("Returning %d applications", len(result))
	return result, len(result), nil
}

func (s *ApplicationService) SearchApplicationsByCandidateID(ctx context.Context, candidateID int) ([]*talenthub.Application, int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	var result []*talenthub.Application
	for _, app := range s.applications {
		if app.CandidateID == candidateID {
			appCopy := *app
			result = append(result, &appCopy)
		}
	}
	
	return result, len(result), nil
}

func (s *ApplicationService) SearchApplicationsByVacancyID(ctx context.Context, vacancyID int) ([]*talenthub.Application, int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	var result []*talenthub.Application
	for _, app := range s.applications {
		if app.VacancyID == vacancyID {
			appCopy := *app
			result = append(result, &appCopy)
		}
	}
	
	return result, len(result), nil
}

func (s *ApplicationService) RegisterApplication(ctx context.Context, application *talenthub.Application) (*talenthub.Application, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	log.Printf("RegisterApplication called with: CandidateID=%d, VacancyID=%d, Status=%s", 
		application.CandidateID, application.VacancyID, application.Status)
	
	// Assign new ID and set application date if not set
	application.ID = s.nextID
	s.nextID++
	
	if application.ApplicationDate.IsZero() {
		application.ApplicationDate = time.Now()
	}
	
	// Create a copy to avoid reference issues
	appCopy := &talenthub.Application{
		ID:              application.ID,
		CandidateID:     application.CandidateID,
		VacancyID:       application.VacancyID,
		Score:           application.Score,
		Status:          application.Status,
		ApplicationDate: application.ApplicationDate,
	}
	
	// Add to in-memory storage
	s.applications = append(s.applications, appCopy)
	
	log.Printf("Application registered with ID=%d. Total applications: %d", appCopy.ID, len(s.applications))
	
	// Return the original application with updated ID
	return application, nil
}

func (s *ApplicationService) UnregisterApplication(ctx context.Context, id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Find and remove application
	for i, app := range s.applications {
		if app.ID == id {
			s.applications = append(s.applications[:i], s.applications[i+1:]...)
			return nil
		}
	}
	
	return talenthub.Errorf(talenthub.ENOTFOUND, "application not found")
}

func (s *ApplicationService) FindApplicationByID(ctx context.Context, id int) (*talenthub.Application, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	for _, app := range s.applications {
		if app.ID == id {
			appCopy := *app
			return &appCopy, nil
		}
	}
	
	return nil, talenthub.Errorf(talenthub.ENOTFOUND, "application not found")
}

func (s *ApplicationService) UpdateApplication(ctx context.Context, id int, update talenthub.ApplicationUpdate) (*talenthub.Application, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	for _, app := range s.applications {
		if app.ID == id {
			if update.Status != nil {
				app.Status = *update.Status
			}
			if update.Score != nil {
				app.Score = *update.Score
			}
			
			appCopy := *app
			return &appCopy, nil
		}
	}
	
	return nil, talenthub.Errorf(talenthub.ENOTFOUND, "application not found")
}

