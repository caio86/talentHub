package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	talenthub "github.com/caio86/talentHub"
)

func (s *Server) loadVagaRoutes(r *http.ServeMux) {
	r.HandleFunc("GET /vaga/{id}", s.handleVagaGet)
	r.HandleFunc("GET /vaga", s.handleVagaList)
	r.HandleFunc("POST /vaga", s.handleVagaCreate)
	r.HandleFunc("PUT /vaga/{id}", s.handleVagaUpdate)
}

type listVagaResponse struct {
	Vagas []*talenthub.Vaga `json:"vagas"`
	Total int               `json:"total"`
}

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

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(vaga); err != nil {
		return
	}
}

func (s *Server) handleVagaList(w http.ResponseWriter, r *http.Request) {
	var filter talenthub.VagaFilter

	vagas, total, err := s.VagaService.FindVagas(r.Context(), filter)
	if err != nil {
		Error(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(listVagaResponse{
		Vagas: vagas,
		Total: total,
	})
}

func (s *Server) handleVagaCreate(w http.ResponseWriter, r *http.Request) {
	var vaga talenthub.Vaga
	if err := json.NewDecoder(r.Body).Decode(&vaga); err != nil {
		Error(w, r, talenthub.Errorf(talenthub.EINVALID, "invalid json body"))
		return
	}

	err := s.VagaService.CreateVaga(r.Context(), &vaga)
	if err != nil {
		Error(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(vaga)
}

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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(updated)
}
