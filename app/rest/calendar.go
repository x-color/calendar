package rest

import (
	"encoding/json"
	"net/http"

	"github.com/x-color/calendar/service/calendar"
)

type CalEndpoint struct {
	service calendar.Service
}

func (e *CalEndpoint) getCalendarsHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(msgContent{"get calendars"})
}

func (e *CalEndpoint) makeCalendarHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(msgContent{"make calendar"})
}

func (e *CalEndpoint) removeCalendarHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(msgContent{"remove calendar"})
}

func (e *CalEndpoint) changeCalendarHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(msgContent{"change calendar"})
}
