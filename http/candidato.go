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
	r.HandleFunc("PATCH /candidato/{id}", s.handleCandidatoPatch)
}

// DTO

type listCandidatoResponse struct {
	Candidatos []*candidatoDTO `json:"candidatos"`
	Total      int             `json:"total"`
}

type experience struct {
	Company string `json:"company"`
	Role    string `json:"role"`
	Years   int    `json:"years"`
}

type education struct {
	Institution string `json:"institution"`
	Course      string `json:"course"`
	Level       string `json:"level"`
}

type candidatoDTO struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Phone    string  `json:"phone,omitempty"`
	Address  string  `json:"address,omitempty"`
	Linkedin *string `json:"linkedin"`

	Experiences []*experience `json:"experiences"`
	Education   []*education  `json:"education"`
	Skills      []string      `json:"skills"`
	Interests   []string      `json:"interests"`

	ResumeLink *string `json:"resume_pdf_path"`
	IsReserve  bool    `json:"is_reserve"`
}

type createCandidatoDTO struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone,omitempty"`
	Address  string `json:"address,omitempty"`
	Linkedin string `json:"linkedin,omitempty"`

	Experiences []*struct {
		Company string `json:"company"`
		Role    string `json:"role"`
		Years   int    `json:"years"`
	} `json:"experience"`
	Education []*struct {
		Institution string `json:"institution"`
		Course      string `json:"course"`
		Level       string `json:"level"`
	} `json:"education"`
	Skills    []string `json:"skills"`
	Interests []string `json:"interests"`

	ResumeLink string `json:"resume_pdf_path,omitempty"`
}

// DTO Helpers

func (d *candidatoDTO) toDomain() *talenthub.Candidato {
	canEdu := make([]*talenthub.Education, len(d.Education))
	for k, v := range d.Education {
		canEdu[k] = &talenthub.Education{
			ID:          0, // será definido pelo banco
			CandidateID: 0, // será definido pelo banco
			Institution: v.Institution,
			Course:      v.Course,
			Level:       v.Level,
		}
	}

	canExp := make([]*talenthub.Experience, len(d.Experiences))
	for k, v := range d.Experiences {
		canExp[k] = &talenthub.Experience{
			ID:          0, // será definido pelo banco
			CandidateID: 0, // será definido pelo banco
			Company:     v.Company,
			Role:        v.Role,
			Years:       v.Years,
		}
	}

	id, _ := strconv.Atoi(d.ID)

	linkedin := ""
	if d.Linkedin != nil {
		linkedin = *d.Linkedin
	}

	resumeLink := ""
	if d.ResumeLink != nil {
		resumeLink = *d.ResumeLink
	}

	return &talenthub.Candidato{
		ID:       id,
		Name:     d.Name,
		Email:    d.Email,
		Password: d.Password,
		Phone:    d.Phone,
		Address:  d.Address,
		Linkedin: linkedin,

		Experiences: canExp,
		Education:   canEdu,
		Skills:      d.Skills,
		Interests:   d.Interests,

		ResumeLink: resumeLink,
	}
}

func (d *candidatoDTO) fromDomain(domain *talenthub.Candidato) {
	canEdu := make([]*education, len(domain.Education))
	for k, v := range domain.Education {
		canEdu[k] = &education{
			Institution: v.Institution,
			Course:      v.Course,
			Level:       v.Level,
		}
	}

	canExp := make([]*experience, len(domain.Experiences))
	for k, v := range domain.Experiences {
		canExp[k] = &experience{
			Company: v.Company,
			Role:    v.Role,
			Years:   v.Years,
		}
	}

	d.ID = strconv.Itoa(domain.ID)
	d.Name = domain.Name
	d.Email = domain.Email
	d.Password = domain.Password
	d.Phone = domain.Phone
	d.Address = domain.Address

	if domain.Linkedin != "" {
		d.Linkedin = &domain.Linkedin
	} else {
		d.Linkedin = nil
	}

	d.Experiences = canExp
	d.Education = canEdu
	d.Skills = domain.Skills
	d.Interests = domain.Interests

	if domain.ResumeLink != "" {
		d.ResumeLink = &domain.ResumeLink
	} else {
		d.ResumeLink = nil
	}

	// Por enquanto, vamos definir is_reserve como false por padrão
	// Isso pode ser ajustado conforme a lógica de negócio
	d.IsReserve = false
}

func (d *createCandidatoDTO) toDomain() *talenthub.Candidato {
	canEdu := make([]*talenthub.Education, len(d.Education))
	for k, v := range d.Education {
		canEdu[k] = &talenthub.Education{
			ID:          0,
			CandidateID: 0,
			Institution: v.Institution,
			Course:      v.Course,
			Level:       v.Level,
		}
	}

	canExp := make([]*talenthub.Experience, len(d.Experiences))
	for k, v := range d.Experiences {
		canExp[k] = &talenthub.Experience{
			ID:          0,
			CandidateID: 0,
			Company:     v.Company,
			Role:        v.Role,
			Years:       v.Years,
		}
	}

	return &talenthub.Candidato{
		ID:       0,
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

// HTTP Handlers

// @summary Get Candidato By ID
// @description Get Candidato By ID
// @router /candidato/{id} [get]
// @tags Candidatos
// @produce json
// @param id path int true "Candidato ID"
// @success 200 {object} http.candidatoDTO "Candidato achado"
// @success 400 {object} http.ErrorResponse "Bad request"
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
// @param email query string false "Email to search"
// @success 200 {object} http.listCandidatoResponse "Lista de candidatos"
// @success 400 {object} http.ErrorResponse "Bad request"
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

	var candidates []*talenthub.Candidato
	var err error
	if email := queryParams.Get("email"); email != "" {
		candidato, err := s.CandidatoService.FindCandidatoByEmail(r.Context(), email)
		if err != nil {
			Error(w, r, err)
			return
		}

		candidates = []*talenthub.Candidato{candidato}
	} else {
		candidates, _, err = s.CandidatoService.FindCandidatos(r.Context(), filter)
		if err != nil {
			Error(w, r, err)
			return
		}
	}

	res := make([]*candidatoDTO, len(candidates))
	for k, v := range candidates {
		var dto candidatoDTO
		dto.fromDomain(v)
		res[k] = &dto
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// @summary Create candidato
// @description Create candidato
// @router /candidato [post]
// @tags Candidatos
// @accept json
// @produce json
// @param candidato body http.createCandidatoDTO true "Candidato a ser criado"
// @success 201 {object} http.candidatoDTO "Candidato criado"
// @success 400 {object} http.ErrorResponse "Bad request"
// @success 404 {object} http.ErrorResponse "Mensagem de erro"
// @success 409 {object} http.ErrorResponse "email already exists"
// @success 500 {object} http.ErrorResponse "Internal Error"
func (s *Server) handleCandidatoCreate(w http.ResponseWriter, r *http.Request) {
	var candidato createCandidatoDTO

	if err := json.NewDecoder(r.Body).Decode(&candidato); err != nil {
		Error(w, r, talenthub.Errorf(talenthub.EINVALID, "invalid json body"))
		return
	}

	domain := candidato.toDomain()
	if err := domain.Validate(); err != nil {
		Error(w, r, err)
		return
	}

	newCandidato, err := s.CandidatoService.CreateCandidato(r.Context(), domain)
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
// @success 400 {object} http.ErrorResponse "Bad request"
// @success 404 {object} http.ErrorResponse "Mensagem de erro"
// @success 500 {object} http.ErrorResponse "Internal Error"
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

// @summary Patch candidato
// @description Patch candidato (partial update)
// @router /candidato/{id} [patch]
// @tags Candidatos
// @accept json
// @produce json
// @param id path int true "Candidato ID"
// @param candidato body talenthub.CandidatoUpdate true "Dados de candidatos para atualizar parcialmente"
// @success 202 {object} http.candidatoDTO "Candidato atualizado"
// @success 400 {object} http.ErrorResponse "Bad request"
// @success 404 {object} http.ErrorResponse "Mensagem de erro"
// @success 500 {object} http.ErrorResponse "Internal Error"
func (s *Server) handleCandidatoPatch(w http.ResponseWriter, r *http.Request) {
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
