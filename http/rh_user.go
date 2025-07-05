package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	talenthub "github.com/caio86/talentHub"
)

func (s *Server) loadRHUserRoutes(r *http.ServeMux) {
	r.HandleFunc("GET /rh_users", s.handleRHUserList)
}

// DTO

type rhUserDTO struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// DTO Helpers

func (d *rhUserDTO) fromDomain(domain *talenthub.RHUser) {
	d.ID = strconv.Itoa(domain.ID)
	d.Name = domain.Name
	d.Email = domain.Email
	d.Password = domain.Password
}

// HTTP Handlers

// @summary Lista RH Users
// @description Lista RH Users
// @router /rh_users [get]
// @tags RH Users
// @produce json
// @success 200 {array} http.rhUserDTO "Lista de usuários RH"
// @success 400 {object} http.ErrorResponse "Bad request"
// @success 404 {object} http.ErrorResponse "Mensagem de erro"
func (s *Server) handleRHUserList(w http.ResponseWriter, r *http.Request) {
	// Por enquanto, retornar dados mock
	// Em uma implementação real, isso viria do RHUserService
	res := []*rhUserDTO{
		{
			ID:       "1",
			Name:     "Gestor RH",
			Email:    "rh@empresa.com",
			Password: "rhpassword",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

