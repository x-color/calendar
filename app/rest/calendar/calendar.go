package calendar

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/x-color/calendar/app/rest/middlewares"
	as "github.com/x-color/calendar/auth/service"
	"github.com/x-color/calendar/calendar/model"
	"github.com/x-color/calendar/calendar/service"
	cs "github.com/x-color/calendar/calendar/service"
	cctx "github.com/x-color/calendar/model/ctx"
	cerror "github.com/x-color/calendar/model/error"
)

type CalendarContent struct {
	ID     string        `json:"id"`
	Name   string        `json:"name"`
	Color  string        `json:"color"`
	Shares []string      `json:"shares"`
	Plans  []PlanContent `json:"plans"`
}

func calModelToContent(cal model.Calendar) CalendarContent {
	plans := make([]PlanContent, len(cal.Plans))
	for i, p := range cal.Plans {
		plans[i] = planModelToContent(p)
	}

	c := CalendarContent{
		ID:     cal.ID,
		Name:   cal.Name,
		Color:  string(cal.Color),
		Shares: cal.Shares,
		Plans:  plans,
	}
	return c
}

func planModelToContent(plan model.Plan) PlanContent {
	p := PlanContent{
		ID:         plan.ID,
		CalendarID: plan.CalendarID,
		Name:       plan.Name,
		Memo:       plan.Memo,
		Color:      string(plan.Color),
		Private:    plan.Private,
		Shares:     plan.Shares,
		IsAllDay:   plan.Period.IsAllDay,
		Begin:      plan.Period.Begin.Unix(),
		End:        plan.Period.End.Unix(),
	}
	return p
}

type calEndpoint struct {
	service service.Service
}

func (e *calEndpoint) GetCalendarsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(cctx.UserIDKey).(string)
	cl, err := e.service.GetCalendars(r.Context(), userID)
	if errors.Is(err, cerror.ErrInvalidContent) {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cals := make([]CalendarContent, len(cl))
	for i, c := range cl {
		cals[i] = calModelToContent(c)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cals)
}

func (e *calEndpoint) MakeCalendarHandler(w http.ResponseWriter, r *http.Request) {
	req := CalendarContent{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := r.Context().Value(cctx.UserIDKey).(string)
	cal, err := e.service.MakeCalendar(r.Context(), userID, req.Name, req.Color)
	if errors.Is(err, cerror.ErrInvalidContent) {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(CalendarContent{
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
		return
	} else if errors.Is(err, cerror.ErrAuthorization) {
		w.WriteHeader(http.StatusForbidden)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (e *calEndpoint) ChangeCalendarHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	req := CalendarContent{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	color, err := model.ConvertToColor(req.Color)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
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
		return
	} else if errors.Is(err, cerror.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if errors.Is(err, cerror.ErrAuthorization) {
		w.WriteHeader(http.StatusForbidden)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func NewCalendarRouter(r *mux.Router, calService cs.Service, authService as.Service) {
	e := calEndpoint{calService}
	r.Use(middlewares.ReqIDMiddleware)
	r.Use(middlewares.ResponseHeaderMiddleware)
	r.Use(middlewares.AuthorizationMiddleware(authService))
	r.Use(userCheckerMiddleware(calService))
	r.HandleFunc("", e.GetCalendarsHandler).Methods(http.MethodGet)
	r.HandleFunc("", e.MakeCalendarHandler).Methods(http.MethodPost)
	r.HandleFunc("/{id}", e.RemoveCalendarHandler).Methods(http.MethodDelete)
	r.HandleFunc("/{id}", e.ChangeCalendarHandler).Methods(http.MethodPatch)
}
