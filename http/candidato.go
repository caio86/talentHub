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

// DTO

type experience struct {
	ID      int    `json:"experience_id"`
	Company string `json:"company"`
	Role    string `json:"role"`
	Years   int    `json:"years"`
}

type education struct {
	ID          int    `json:"education_id"`
	Institution string `json:"institution"`
	Course      string `json:"course"`
	Level       string `json:"level"`
}

type candidatoDTO struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone,omitempty"`
	Address  string `json:"address,omitempty"`
	Linkedin string `json:"linkedin,omitempty"`

	Experiences []*experience `json:"experience"`
	Education   []*education  `json:"education"`
	Skills      []string      `json:"skills"`
	Interests   []string      `json:"interests"`

	ResumeLink string `json:"resume_link,omitempty"`
}

// DTO Helpers

func (d *candidatoDTO) toDomain() *talenthub.Candidato {
	canEdu := make([]*talenthub.Education, len(d.Education))
	for k, v := range d.Education {
		canEdu[k] = &talenthub.Education{
			ID:          v.ID,
			CandidateID: d.ID,
			Institution: v.Institution,
			Course:      v.Course,
			Level:       v.Level,
		}
	}

	canExp := make([]*talenthub.Experience, len(d.Experiences))
	for k, v := range d.Experiences {
		canExp[k] = &talenthub.Experience{
			ID:          v.ID,
			CandidateID: d.ID,
			Company:     v.Company,
			Role:        v.Role,
			Years:       v.Years,
		}
	}

	return &talenthub.Candidato{
		ID:       d.ID,
		Name:     d.Name,
		Email:    d.Email,
		Password: d.Password,
		Phone:    d.Phone,
		Address:  d.Address,
		Linkedin: d.Linkedin,

		Experiences: canExp,
		Education:   canEdu,
		Skills:      d.Skills,
		Interests:   d.Interests,

		ResumeLink: d.ResumeLink,
	}
}

func (d *candidatoDTO) fromDomain(domain *talenthub.Candidato) {
	canEdu := make([]*education, len(domain.Education))
	for k, v := range domain.Education {
		canEdu[k] = &education{
			ID:          v.ID,
			Institution: v.Institution,
			Course:      v.Course,
			Level:       v.Level,
		}
	}

	canExp := make([]*experience, len(domain.Experiences))
	for k, v := range domain.Experiences {
		canExp[k] = &experience{
			ID:      v.ID,
			Company: v.Company,
			Role:    v.Role,
			Years:   v.Years,
		}
	}

	d.ID = domain.ID
	d.Name = domain.Name
	d.Email = domain.Email
	d.Password = domain.Password
	d.Phone = domain.Phone
	d.Address = domain.Address
	d.Linkedin = domain.Linkedin

	d.Experiences = canExp
	d.Education = canEdu
	d.Skills = domain.Skills
	d.Interests = domain.Interests

	d.ResumeLink = domain.ResumeLink
}

type listCandidatoResponse struct {
	Candidatos []*candidatoDTO `json:"candidatos"`
	Total      int             `json:"total"`
}

// HTTP Handlers

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

	var res candidatoDTO
	res.fromDomain(candidate)

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
		var dto candidatoDTO
		dto.fromDomain(v)
		res[k] = &dto
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

	newCandidato, err := s.CandidatoService.CreateCandidato(r.Context(), candidato.toDomain())
	if err != nil {
		Error(w, r, err)
		return
	}

	var res candidatoDTO
	res.fromDomain(newCandidato)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
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

	var res candidatoDTO
	res.fromDomain(updated)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(res)
}
