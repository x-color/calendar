package rest

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	cerror "github.com/x-color/calendar/model/error"
	"github.com/x-color/calendar/service/calendar"
)

type calendarContent struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Color   string   `json:"color"`
	Private bool     `json:"private"`
	Shares  []string `json:"shares"`
	Plans   []string `json:"plans"`
}

type CalEndpoint struct {
	service calendar.Service
}

func (e *CalEndpoint) getCalendarsHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(msgContent{"get calendars"})
}

func (e *CalEndpoint) makeCalendarHandler(w http.ResponseWriter, r *http.Request) {
	req := calendarContent{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msgContent{"bad contents"})
		return
	}

	cal, err := e.service.MakeCalendar(r.Context(), req.Name, req.Color)
	if errors.Is(err, cerror.ErrInvalidContent) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msgContent{"bad contents"})
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msgContent{"internal server error"})
		return
	}

	plans := []string{}
	for _, p := range cal.Plans {
		plans = append(plans, p.ID)
	}

	json.NewEncoder(w).Encode(calendarContent{
		ID:      cal.ID,
		Name:    cal.Name,
		Color:   string(cal.Color),
		Private: cal.Private,
		Shares:  cal.Shares,
		Plans:   plans,
	})
}

func (e *CalEndpoint) removeCalendarHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	err := e.service.RemoveCalendar(r.Context(), vars["id"])
	if errors.Is(err, cerror.ErrInvalidContent) || errors.Is(err, cerror.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(msgContent{"not found"})
		return
	} else if errors.Is(err, cerror.ErrAuthorization) {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(msgContent{"unauthorization"})
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msgContent{"internal server error"})
		return
	}

	json.NewEncoder(w).Encode(msgContent{"remove calendar"})
}

func (e *CalEndpoint) changeCalendarHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(msgContent{"change calendar"})
}
