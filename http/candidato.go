package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	talenthub "github.com/caio86/talentHub"
)

type candidatoHandler struct {
	candidatoService talenthub.CandidatoService
}

func newCandidatoHandler(svc talenthub.CandidatoService) *candidatoHandler {
	return &candidatoHandler{
		candidatoService: svc,
	}
}

func (h *candidatoHandler) loadRoutes(r *http.ServeMux) {
	r.HandleFunc("GET /candidato/{id}", h.get)
	r.HandleFunc("GET /candidato", h.list)
	r.HandleFunc("POST /candidato", h.create)
	r.HandleFunc("PUT /candidato/{id}", h.update)
}

type listResponse struct {
	Candidatos []*talenthub.Candidato `json:"candidatos"`
	Total      int                    `json:"total"`
}

func (h *candidatoHandler) get(w http.ResponseWriter, r *http.Request) {
	ids := r.PathValue("id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		Error(w, r, talenthub.Errorf(talenthub.EINVALID, "invalid id"))
		return
	}

	candidate, err := h.candidatoService.FindCandidatoByID(r.Context(), id)
	if err != nil {
		Error(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(candidate); err != nil {
		return
	}
}

func (h *candidatoHandler) list(w http.ResponseWriter, r *http.Request) {
	var filter talenthub.CandidatoFilter

	candidates, total, err := h.candidatoService.FindCandidatos(r.Context(), filter)
	if err != nil {
		Error(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(listResponse{
		Candidatos: candidates,
		Total:      total,
	})
}

func (h *candidatoHandler) create(w http.ResponseWriter, r *http.Request) {
	var candidato talenthub.Candidato

	if err := json.NewDecoder(r.Body).Decode(&candidato); err != nil {
		Error(w, r, talenthub.Errorf(talenthub.EINVALID, "invalid json body"))
		return
	}

	err := h.candidatoService.CreateCandidato(r.Context(), &candidato)
	if err != nil {
		Error(w, r, err)
		return
	}

	json.NewEncoder(w).Encode(candidato)
}

func (h *candidatoHandler) update(w http.ResponseWriter, r *http.Request) {
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

	updated, err := h.candidatoService.UpdateCandidato(r.Context(), id, upd)
	if err != nil {
		Error(w, r, err)
		return
	}

	json.NewEncoder(w).Encode(updated)
}
