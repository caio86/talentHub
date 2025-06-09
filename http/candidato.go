package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	talenthub "github.com/caio86/talentHub"
)

func (s *Server) loadCandidatoRoutes(r *http.ServeMux) {
	r.HandleFunc("GET /candidato/{id}", s.handleCandidatoGet)
	r.HandleFunc("GET /candidato", s.handleCandidatoList)
	r.HandleFunc("POST /candidato", s.handleCandidatoCreate)
	r.HandleFunc("PUT /candidato/{id}", s.handleCandidatoUpdate)
}

type candidatoDTO struct {
	ID    int    `json:"-"`
	Name  string `json:"name"`
	Email string `json:"email"`
	CPF   string `json:"cpf"`
	Phone string `json:"phone"`
}

type listCandidatoResponse struct {
	Candidatos []*candidatoDTO `json:"candidatos"`
	Total      int             `json:"total"`
}

// @summary Get Candidato By ID
// @description Get Candidato By ID
// @router /candidato/{id} [get]
// @tags Candidatos
// @produce json
// @param id path int true "Candidato ID"
// @success 200 {object} http.candidatoDTO "Candidato achado"
// @success 404 {object} http.ErrorResponse "Mensagem de error"
func (s *Server) handleCandidatoGet(w http.ResponseWriter, r *http.Request) {
	ids := r.PathValue("id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		Error(w, r, talenthub.Errorf(talenthub.EINVALID, "invalid id"))
		return
	}

	candidate, err := s.CandidatoService.FindCandidatoByID(r.Context(), id)
	if err != nil {
		Error(w, r, err)
		return
	}

	res := &candidatoDTO{
		ID:    candidate.ID,
		Name:  candidate.Name,
		Email: candidate.Email,
		CPF:   candidate.CPF,
		Phone: candidate.Phone,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		return
	}
}

// @summary Lista Candidatos
// @description Lista Candidatos
// @router /candidato [get]
// @tags Candidatos
// @produce json
// @param limit query int false "Pagination limit"
// @param offset query int false "Pagination offset"
// @success 200 {object} http.listCandidatoResponse "Lista de candidatos"
// @success 404 {object} http.ErrorResponse "Mensagem de erro"
func (s *Server) handleCandidatoList(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	var filter talenthub.CandidatoFilter

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

	candidates, total, err := s.CandidatoService.FindCandidatos(r.Context(), filter)
	if err != nil {
		Error(w, r, err)
		return
	}

	res := make([]*candidatoDTO, len(candidates))
	for k, v := range candidates {
		res[k] = &candidatoDTO{
			ID:    v.ID,
			Name:  v.Name,
			Email: v.Email,
			CPF:   v.CPF,
			Phone: v.Phone,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(listCandidatoResponse{
		Candidatos: res,
		Total:      total,
	})
}

// @summary Create candidato
// @description Create candidato
// @router /candidato [post]
// @tags Candidatos
// @accept json
// @produce json
// @param candidato body http.candidatoDTO true "Candidato a ser criado"
// @success 201 {object} http.candidatoDTO "Candidato criado"
// @success 404 {object} http.ErrorResponse "Mensagem de erro"
func (s *Server) handleCandidatoCreate(w http.ResponseWriter, r *http.Request) {
	var candidato candidatoDTO

	if err := json.NewDecoder(r.Body).Decode(&candidato); err != nil {
		Error(w, r, talenthub.Errorf(talenthub.EINVALID, "invalid json body"))
		return
	}

	err := s.CandidatoService.CreateCandidato(r.Context(), &talenthub.Candidato{
		ID:    candidato.ID,
		Name:  candidato.Name,
		Email: candidato.Email,
		CPF:   candidato.CPF,
		Phone: candidato.Phone,
	})
	if err != nil {
		Error(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(candidato)
}

// @summary Update candidato
// @description Update candidato
// @router /candidato/{id} [put]
// @tags Candidatos
// @accept json
// @produce json
// @param id path int true "Candidato ID"
// @param candidato body talenthub.CandidatoUpdate true "Dados de candidatos para atualizar"
// @success 202 {object} http.candidatoDTO "Candidato atualizado"
// @success 404 {object} http.ErrorResponse "Mensagem de erro"
func (s *Server) handleCandidatoUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		Error(w, r, talenthub.Errorf(talenthub.EINVALID, "invalid id"))
		return
	}

	var upd talenthub.CandidatoUpdate
	if err := json.NewDecoder(r.Body).Decode(&upd); err != nil {
		Error(w, r, talenthub.Errorf(talenthub.EINVALID, "invalid json body"))
		return
	}

	updated, err := s.CandidatoService.UpdateCandidato(r.Context(), id, upd)
	if err != nil {
		Error(w, r, err)
		return
	}

	res := &candidatoDTO{
		ID:    updated.ID,
		Name:  updated.Name,
		Email: updated.Email,
		CPF:   updated.CPF,
		Phone: updated.Phone,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(res)
}
