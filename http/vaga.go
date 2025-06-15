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
}

// DTO

type vagaDTO struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`

	IsActive bool `json:"IsActive"`

	Area         string   `json:"area"`
	Type         string   `json:"type"`
	Location     string   `json:"location"`
	Requirements []string `json:"requirements"`

	Posted_date time.Time `json:"posted_date"`
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

type listVagaResponse struct {
	Vagas []*vagaDTO `json:"vagas"`
	Total int        `json:"total"`
}

// HTTP Handlers

// @summary Get Vaga By ID
// @description Get Vaga By ID
// @router /vaga/{id} [get]
// @tags Vagas
// @produce json
// @param id path int true "Vaga ID"
// @success 200 {object} http.vagaDTO "Vaga achada"
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
		res[k].fromDomain(v)
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
// @param candidato body http.vagaDTO true "Vaga a ser criada"
// @success 201 {object} http.vagaDTO "Vaga criada"
// @success 404 {object} http.ErrorResponse "Mensagem de erro"
func (s *Server) handleVagaCreate(w http.ResponseWriter, r *http.Request) {
	var vagadto *vagaDTO
	if err := json.NewDecoder(r.Body).Decode(&vagadto); err != nil {
		Error(w, r, talenthub.Errorf(talenthub.EINVALID, "invalid json body"))
		return
	}

	vaga := vagadto.toDomain()

	err := s.VagaService.CreateVaga(r.Context(), vaga)
	if err != nil {
		Error(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(vagadto)
}

// @summary Update candidato
// @description Update candidato
// @router /vaga/{id} [put]
// @tags Vagas
// @accept json
// @produce json
// @param id path int true "Vaga ID"
// @param candidato body talenthub.VagaUpdate true "Dados de vagas para atualizar"
// @success 202 {object} http.vagaDTO "Vaga atualizada"
// @success 404 {object} http.ErrorResponse "Mensagem de erro"
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
