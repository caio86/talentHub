package http

import (
	"encoding/json"
	"net/http"

	talenthub "github.com/caio86/talentHub"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func Error(w http.ResponseWriter, r *http.Request, err error) {
	code, message := talenthub.ErrorCode(err), talenthub.ErrorMessage(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(ErrorStatusCode(code))
	json.NewEncoder(w).Encode(&ErrorResponse{Error: message})
}

var codes = map[string]int{
	talenthub.EINVALID:        http.StatusBadRequest,
	talenthub.EINTERNAL:       http.StatusInternalServerError,
	talenthub.ENOTFOUND:       http.StatusNotFound,
	talenthub.ENOTIMPLEMENTED: http.StatusNotImplemented,
}

func ErrorStatusCode(code string) int {
	if v, ok := codes[code]; ok {
		return v
	}
	return http.StatusInternalServerError
}

func FromErrorStatusCode(code int) string {
	for k, v := range codes {
		if v == code {
			return k
		}
	}
	return talenthub.EINTERNAL
}
