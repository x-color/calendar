package calendar

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/x-color/calendar/calendar/model"
	"github.com/x-color/calendar/calendar/service"
	cctx "github.com/x-color/calendar/model/ctx"
	cerror "github.com/x-color/calendar/model/error"
)

type calendarContent struct {
	ID     string        `json:"id"`
	Name   string        `json:"name"`
	Color  string        `json:"color"`
	Shares []string      `json:"shares"`
	Plans  []planContent `json:"plans"`
}

type calEndpoint struct {
	service service.Service
}

func (e *calEndpoint) GetCalendarsHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(msgContent{"get calendars"})
}

func (e *calEndpoint) MakeCalendarHandler(w http.ResponseWriter, r *http.Request) {
	req := calendarContent{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msgContent{"bad contents"})
		return
	}

	userID := r.Context().Value(cctx.UserIDKey).(string)
	cal, err := e.service.MakeCalendar(r.Context(), userID, req.Name, req.Color)
	if errors.Is(err, cerror.ErrInvalidContent) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msgContent{"bad contents"})
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msgContent{"internal server error"})
		return
	}

	json.NewEncoder(w).Encode(calendarContent{
		ID:     cal.ID,
		Name:   cal.Name,
		Color:  string(cal.Color),
		Shares: cal.Shares,
	})
}

func (e *calEndpoint) RemoveCalendarHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := r.Context().Value(cctx.UserIDKey).(string)
	err := e.service.RemoveCalendar(r.Context(), userID, vars["id"])
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

func (e *calEndpoint) ChangeCalendarHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	req := calendarContent{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msgContent{"bad contents"})
		return
	}

	color, err := model.ConvertToColor(req.Color)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msgContent{"bad contents"})
		return
	}

	cal := model.Calendar{
		ID:     vars["id"],
		Name:   req.Name,
		Color:  color,
		Shares: req.Shares,
	}

	userID := r.Context().Value(cctx.UserIDKey).(string)
	err = e.service.ChangeCalendar(r.Context(), userID, cal)
	if errors.Is(err, cerror.ErrInvalidContent) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msgContent{"bad contents"})
		return
	} else if errors.Is(err, cerror.ErrNotFound) {
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

	json.NewEncoder(w).Encode(msgContent{"change calendar"})
}
