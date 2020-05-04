package rest

import (
	"encoding/json"
	"net/http"

	cs "github.com/x-color/calendar/service/calendar"
)

type UserEndpoint struct {
	service cs.Service
}

func (e *UserEndpoint) registerHandler(w http.ResponseWriter, r *http.Request) {
	_, err := e.service.RegisterUser(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msgContent{"internal server error"})
		return
	}

	json.NewEncoder(w).Encode(msgContent{Msg: "register"})
}
