package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	talenthub "github.com/caio86/talentHub"
)

func (s *Server) loadVagaRoutes(r *http.ServeMux) {
	r.HandleFunc("GET /vaga/{id}", s.handleVagaGet)
	r.HandleFunc("GET /vaga", s.handleVagaList)
	r.HandleFunc("POST /vaga", s.handleVagaCreate)
	r.HandleFunc("PUT /vaga/{id}", s.handleVagaUpdate)
	r.HandleFunc("POST /vaga/open/{id}", s.handleVagaOpen)
	r.HandleFunc("POST /vaga/close/{id}", s.handleVagaClose)
}

// DTO

type vagaDTO struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`

	IsActive bool `json:"is_active"`

	Area         string   `json:"area,omitempty"`
	Type         string   `json:"type,omitempty"`
	Location     string   `json:"location,omitempty"`
	Requirements []string `json:"requirements"`

	Posted_date time.Time `json:"posted_date"`
}

type createVagaDTO struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`

	Area         string   `json:"area,omitempty"`
	Type         string   `json:"type,omitempty"`
	Location     string   `json:"location,omitempty"`
	Requirements []string `json:"requirements"`
}

type listVagaResponse struct {
	Vagas []*vagaDTO `json:"vagas"`
	Total int        `json:"total"`
}

// DTO Helpers

func (d *vagaDTO) toDomain() *talenthub.Vaga {
	return &talenthub.Vaga{
		ID:          d.ID,
		Title:       d.Title,
		Description: d.Description,

		IsActive: d.IsActive,

		Area:         d.Area,
		Type:         d.Type,
		Location:     d.Location,
		Requirements: d.Requirements,

		Posted_date: d.Posted_date,
	}
}

func (d *vagaDTO) fromDomain(domain *talenthub.Vaga) {
	d.ID = domain.ID
	d.Title = domain.Title
	d.Description = domain.Description

	d.IsActive = domain.IsActive

	d.Area = domain.Area
	d.Type = domain.Type
	d.Location = domain.Location
	d.Requirements = domain.Requirements

	d.Posted_date = domain.Posted_date
}

func (d *createVagaDTO) toDomain() *talenthub.Vaga {
	return &talenthub.Vaga{
		ID:          0,
		Title:       d.Title,
		Description: d.Description,

		IsActive: false,

		Area:         d.Area,
		Type:         d.Type,
		Location:     d.Location,
		Requirements: d.Requirements,
	}
}

// HTTP Handlers

// @summary Get Vaga By ID
// @description Get Vaga By ID
// @router /vaga/{id} [get]
// @tags Vagas
// @produce json
// @param id path int true "Vaga ID"
// @success 200 {object} http.vagaDTO "Vaga achada"
// @success 400 {object} http.ErrorResponse "Bad request"
// @success 404 {object} http.ErrorResponse "Mensagem de error"
func (s *Server) handleVagaGet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		Error(w, r, talenthub.Errorf(talenthub.EINVALID, "invalid id"))
		return
	}

	vaga, err := s.VagaService.FindVagaByID(r.Context(), id)
	if err != nil {
		Error(w, r, err)
		return
	}

	var res vagaDTO
	res.fromDomain(vaga)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		return
	}
}

// @summary Lista Vagas
// @description Lista Vagas
// @router /vaga [get]
// @tags Vagas
// @produce json
// @param limit query int false "Pagination limit"
// @param offset query int false "Pagination offset"
// @success 200 {object} http.listVagaResponse "Lista de vagas"
// @success 400 {object} http.ErrorResponse "Bad request"
// @success 404 {object} http.ErrorResponse "Mensagem de erro"
func (s *Server) handleVagaList(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	var filter talenthub.VagaFilter

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

	vagas, total, err := s.VagaService.FindVagas(r.Context(), filter)
	if err != nil {
		Error(w, r, err)
		return
	}

	res := make([]*vagaDTO, len(vagas))
	for k, v := range vagas {
		var dto vagaDTO
		dto.fromDomain(v)
		res[k] = &dto
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(listVagaResponse{
		Vagas: res,
		Total: total,
	})
}

// @summary Create vaga
// @description Create vaga
// @router /vaga [post]
// @tags Vagas
// @accept json
// @produce json
// @param candidato body http.createVagaDTO true "Vaga a ser criada"
// @success 201 {object} http.vagaDTO "Vaga criada"
// @success 400 {object} http.ErrorResponse "Bad request"
// @success 404 {object} http.ErrorResponse "Mensagem de erro"
// @success 500 {object} http.ErrorResponse "Internal Error"
func (s *Server) handleVagaCreate(w http.ResponseWriter, r *http.Request) {
	var vagadto createVagaDTO

	if err := json.NewDecoder(r.Body).Decode(&vagadto); err != nil {
		Error(w, r, talenthub.Errorf(talenthub.EINVALID, "invalid json body, %s", err))
		return
	}

	domain := vagadto.toDomain()
	if err := domain.Validate(); err != nil {
		Error(w, r, err)
		return
	}

	newVaga, err := s.VagaService.CreateVaga(r.Context(), domain)
	if err != nil {
		Error(w, r, err)
		return
	}

	var res vagaDTO
	res.fromDomain(newVaga)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

// @summary Update vaga
// @description Update vaga
// @router /vaga/{id} [put]
// @tags Vagas
// @accept json
// @produce json
// @param id path int true "Vaga ID"
// @param candidato body talenthub.VagaUpdate true "Dados de vagas para atualizar"
// @success 202 {object} http.vagaDTO "Vaga atualizada"
// @success 400 {object} http.ErrorResponse "Bad request"
// @success 404 {object} http.ErrorResponse "Mensagem de erro"
// @success 500 {object} http.ErrorResponse "Internal Error"
func (s *Server) handleVagaUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		Error(w, r, talenthub.Errorf(talenthub.EINVALID, "invalid id"))
		return
	}

	var upd talenthub.VagaUpdate
	if err := json.NewDecoder(r.Body).Decode(&upd); err != nil {
		Error(w, r, talenthub.Errorf(talenthub.EINVALID, "invalid json body"))
		return
	}

	updated, err := s.VagaService.UpdateVaga(r.Context(), id, upd)
	if err != nil {
		Error(w, r, err)
		return
	}

	var res vagaDTO
	res.fromDomain(updated)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(res)
}

// @summary Open vaga
// @description Open vaga
// @router /vaga/open/{id} [post]
// @tags Vagas
// @param id path int true "Vaga ID"
// @success 204 "Vaga aberta"
// @success 400 {object} http.ErrorResponse "Bad request"
// @success 404 {object} http.ErrorResponse "Mensagem de erro"
func (s *Server) handleVagaOpen(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		Error(w, r, talenthub.Errorf(talenthub.EINVALID, "invalid id"))
		return
	}

	err = s.VagaService.OpenVaga(r.Context(), id)
	if err != nil {
		Error(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @summary Close vaga
// @description Close vaga
// @router /vaga/close/{id} [post]
// @tags Vagas
// @param id path int true "Vaga ID"
// @success 204 "Vaga closed"
// @success 400 {object} http.ErrorResponse "Bad request"
// @success 404 {object} http.ErrorResponse "Mensagem de erro"
func (s *Server) handleVagaClose(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		Error(w, r, talenthub.Errorf(talenthub.EINVALID, "invalid id"))
		return
	}

	err = s.VagaService.CloseVaga(r.Context(), id)
	if err != nil {
		Error(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
