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

type listCandidatoResponse struct {
	Candidatos []*talenthub.Candidato `json:"candidatos"`
	Total      int                    `json:"total"`
}

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

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(candidate); err != nil {
		return
	}
}

func (s *Server) handleCandidatoList(w http.ResponseWriter, r *http.Request) {
	var filter talenthub.CandidatoFilter

	candidates, total, err := s.CandidatoService.FindCandidatos(r.Context(), filter)
	if err != nil {
		Error(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(listCandidatoResponse{
		Candidatos: candidates,
		Total:      total,
	})
}

func (s *Server) handleCandidatoCreate(w http.ResponseWriter, r *http.Request) {
	var candidato talenthub.Candidato

	if err := json.NewDecoder(r.Body).Decode(&candidato); err != nil {
		Error(w, r, talenthub.Errorf(talenthub.EINVALID, "invalid json body"))
		return
	}

	err := s.CandidatoService.CreateCandidato(r.Context(), &candidato)
	if err != nil {
		Error(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(candidato)
}

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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(updated)
}
