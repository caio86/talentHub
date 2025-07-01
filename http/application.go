package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	talenthub "github.com/caio86/talentHub"
)

func (s *Server) loadApplicationRoutes(r *http.ServeMux) {
	r.HandleFunc("GET /application/{id}", s.handleApplicationGetByID)
	r.HandleFunc("GET /application", s.handleApplicationList)
	r.HandleFunc("POST /application", s.handleApplicationRegister)
	r.HandleFunc("DELETE /application/{id}", s.handleApplicationUnregister)
	r.HandleFunc("PUT /application/{id}", s.handleApplicationUpdate)
	r.HandleFunc("GET /application/candidato/{id}", s.handleApplicationSearchByCandidateID)
	r.HandleFunc("GET /application/vaga/{id}", s.handleApplicationSearchByVacancyID)
}

// DTO

type listApplicationResponse struct {
	Applications []*applicationDTO `json:"applications"`
	Total        int               `json:"total"`
}

type applicationDTO struct {
	ID              int       `json:"id"`
	CandidateID     int       `json:"candidate_id"`
	VacancyID       int       `json:"vacancy_id"`
	Score           int       `json:"score"`
	Status          string    `json:"status"`
	ApplicationDate time.Time `json:"application_date"`
}

type registerApplicationDTO struct {
	CandidateID int    `json:"candidate_id"`
	VacancyID   int    `json:"vacancy_id"`
	Score       int    `json:"score"`
	Status      string `json:"status"`
}

// DTO Helpers

func (d *applicationDTO) fromDomain(domain *talenthub.Application) {
	d.ID = domain.ID
	d.CandidateID = domain.CandidateID
	d.VacancyID = domain.VacancyID
	d.Score = domain.Score
	d.Status = domain.Status
	d.ApplicationDate = domain.ApplicationDate
}

func (d *registerApplicationDTO) toDomain() *talenthub.Application {
	return &talenthub.Application{
		ID:          0,
		CandidateID: d.CandidateID,
		VacancyID:   d.VacancyID,
		Score:       d.Score,
		Status:      d.Status,
	}
}

// HTTP Handlers

// @summary Get Application By ID
// @description Get Application By ID
// @router /application/{id} [get]
// @tags Applications
// @produce json
// @param id path int true "Application ID"
// @success 200 {object} http.applicationDTO "Application achada"
// @success 400 {object} http.ErrorResponse "Bad request"
// @success 404 {object} http.ErrorResponse "Mensagem de error"
// @success 500 {object} http.ErrorResponse "Internal Error"
func (s *Server) handleApplicationGetByID(w http.ResponseWriter, r *http.Request) {
	ids := r.PathValue("id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		Error(w, r, talenthub.Errorf(talenthub.EINVALID, "invalid id"))
		return
	}

	application, err := s.ApplicationService.FindApplicationByID(r.Context(), id)
	if err != nil {
		Error(w, r, err)
		return
	}

	var res applicationDTO
	res.fromDomain(application)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		return
	}
}

// @summary Lista Applications
// @description Lista Applications
// @router /application [get]
// @tags Applications
// @produce json
// @param limit query int false "Pagination limit"
// @param offset query int false "Pagination offset"
// @success 200 {object} http.listApplicationResponse "Lista de applications"
// @success 400 {object} http.ErrorResponse "Bad request"
// @success 404 {object} http.ErrorResponse "Mensagem de erro"
// @success 500 {object} http.ErrorResponse "Internal Error"
func (s *Server) handleApplicationList(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	var filter talenthub.ApplicationFilter

	if res := queryParams.Get("offset"); res == "" {
		filter.Offset = 0
	} else {
		offset, err := strconv.ParseUint(res, 10, 32)
		filter.Offset = int32(offset)
		if err != nil {
			Error(w, r, talenthub.Errorf(talenthub.EINVALID, "invalid page param"))
			return
		}
	}

	if res := queryParams.Get("limit"); res == "" {
		filter.Limit = 0
	} else {
		limit, err := strconv.ParseUint(res, 10, 32)
		filter.Limit = int32(limit)
		if err != nil {
			Error(w, r, talenthub.Errorf(talenthub.EINVALID, "invalid limit param"))
			return
		}
	}

	applications, total, err := s.ApplicationService.FindApplications(r.Context(), filter)
	if err != nil {
		Error(w, r, err)
		return
	}

	res := make([]*applicationDTO, len(applications))
	for k, v := range applications {
		var dto applicationDTO
		dto.fromDomain(v)
		res[k] = &dto
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(listApplicationResponse{
		Applications: res,
		Total:        total,
	})
}

// @summary Register application
// @description Register application
// @router /application [post]
// @tags Applications
// @accept json
// @produce json
// @param candidato body http.registerApplicationDTO true "Application a ser registrada"
// @success 201 {object} http.applicationDTO "Application criada"
// @success 400 {object} http.ErrorResponse "Bad request"
// @success 404 {object} http.ErrorResponse "Mensagem de erro"
// @success 500 {object} http.ErrorResponse "Internal Error"
func (s *Server) handleApplicationRegister(w http.ResponseWriter, r *http.Request) {
	var dto registerApplicationDTO

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		Error(w, r, talenthub.Errorf(talenthub.EINVALID, "invalid json body, %s", err))
		return
	}

	domain := dto.toDomain()
	if err := domain.Validate(); err != nil {
		Error(w, r, err)
		return
	}

	newApplication, err := s.ApplicationService.RegisterApplication(r.Context(), domain)
	if err != nil {
		Error(w, r, err)
		return
	}

	var res applicationDTO
	res.fromDomain(newApplication)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

// @summary Unregister application
// @description Unregister vaga
// @router /application/{id} [delete]
// @tags Applications
// @param id path int true "Application ID"
// @success 204 "application unregistered"
// @success 400 {object} http.ErrorResponse "Bad request"
// @success 404 {object} http.ErrorResponse "Not found"
// @success 500 {object} http.ErrorResponse "Internal Error"
func (s *Server) handleApplicationUnregister(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		Error(w, r, talenthub.Errorf(talenthub.EINVALID, "invalid id"))
		return
	}

	if err := s.ApplicationService.UnregisterApplication(r.Context(), id); err != nil {
		Error(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @summary Update application
// @description Update application
// @router /application/{id} [put]
// @tags Applications
// @accept json
// @produce json
// @param id path int true "Application ID"
// @param candidato body talenthub.ApplicationUpdate true "Dados de applications para atualizar"
// @success 202 {object} http.applicationDTO "Application atualizada"
// @success 400 {object} http.ErrorResponse "Bad request"
// @success 404 {object} http.ErrorResponse "Mensagem de erro"
// @success 500 {object} http.ErrorResponse "Internal Error"
func (s *Server) handleApplicationUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		Error(w, r, talenthub.Errorf(talenthub.EINVALID, "invalid id"))
		return
	}

	var upd talenthub.ApplicationUpdate
	if err := json.NewDecoder(r.Body).Decode(&upd); err != nil {
		Error(w, r, talenthub.Errorf(talenthub.EINVALID, "invalid json body"))
		return
	}

	updated, err := s.ApplicationService.UpdateApplication(r.Context(), id, upd)
	if err != nil {
		Error(w, r, err)
		return
	}

	var res applicationDTO
	res.fromDomain(updated)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(res)
}

// @summary Pesquisa applications por id da vaga
// @description Pesquisa applications por id da vaga
// @router /application/candidato/{id} [get]
// @tags Applications
// @produce json
// @success 200 {object} http.listApplicationResponse "Lista de applications"
// @success 400 {object} http.ErrorResponse "Bad request"
// @success 404 {object} http.ErrorResponse "Mensagem de erro"
// @success 500 {object} http.ErrorResponse "Internal Error"
func (s *Server) handleApplicationSearchByCandidateID(w http.ResponseWriter, r *http.Request) {
	ids := r.PathValue("id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		Error(w, r, talenthub.Errorf(talenthub.EINVALID, "invalid id"))
		return
	}

	applications, total, err := s.ApplicationService.SearchApplicationsByCandidateID(r.Context(), id)
	if err != nil {
		Error(w, r, err)
		return
	}

	res := make([]*applicationDTO, len(applications))
	for k, v := range applications {
		var dto applicationDTO
		dto.fromDomain(v)
		res[k] = &dto
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(listApplicationResponse{
		Applications: res,
		Total:        total,
	})
}

// @summary Pesquisa applications por id da vaga
// @description Pesquisa applications por id da vaga
// @router /application/vaga/{id} [get]
// @tags Applications
// @produce json
// @success 200 {object} http.listApplicationResponse "Lista de applications"
// @success 400 {object} http.ErrorResponse "Bad request"
// @success 404 {object} http.ErrorResponse "Mensagem de erro"
// @success 500 {object} http.ErrorResponse "Internal Error"
func (s *Server) handleApplicationSearchByVacancyID(w http.ResponseWriter, r *http.Request) {
	ids := r.PathValue("id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		Error(w, r, talenthub.Errorf(talenthub.EINVALID, "invalid id"))
		return
	}

	applications, total, err := s.ApplicationService.SearchApplicationsByVacancyID(r.Context(), id)
	if err != nil {
		Error(w, r, err)
		return
	}

	res := make([]*applicationDTO, len(applications))
	for k, v := range applications {
		var dto applicationDTO
		dto.fromDomain(v)
		res[k] = &dto
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(listApplicationResponse{
		Applications: res,
		Total:        total,
	})
}
