package calendar

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/x-color/calendar/app/rest/middlewares"
	as "github.com/x-color/calendar/auth/service"
	"github.com/x-color/calendar/calendar/model"
	"github.com/x-color/calendar/calendar/service"
	cs "github.com/x-color/calendar/calendar/service"
	cctx "github.com/x-color/calendar/model/ctx"
	cerror "github.com/x-color/calendar/model/error"
)

type PlanContent struct {
	ID         string   `json:"id"`
	CalendarID string   `json:"calendar_id"`
	UserID     string   `json:"user_id"`
	Name       string   `json:"name"`
	Memo       string   `json:"memo"`
	Color      string   `json:"color"`
	Private    bool     `json:"private"`
	Shares     []string `json:"shares"`
	IsAllDay   bool     `json:"is_all_day"`
	Begin      int64    `json:"begin"`
	End        int64    `json:"end"`
}

type planEndpoint struct {
	service service.Service
}

func (e *planEndpoint) ScheduleHandler(w http.ResponseWriter, r *http.Request) {
	req := PlanContent{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	color, err := model.ConvertToColor(req.Color)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	planPram := model.Plan{
		CalendarID: req.CalendarID,
		Name:       req.Name,
		Memo:       req.Memo,
		Color:      color,
		Private:    req.Private,
		Shares:     req.Shares,
		Period: model.Period{
			IsAllDay: req.IsAllDay,
			Begin:    time.Unix(req.Begin, 0),
			End:      time.Unix(req.End, 0),
		},
	}

	userID := r.Context().Value(cctx.UserIDKey).(string)
	plan, err := e.service.Schedule(r.Context(), userID, planPram)
	if errors.Is(err, cerror.ErrInvalidContent) || errors.Is(err, cerror.ErrAuthorization) {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(PlanContent{
		ID:         plan.ID,
		UserID:     plan.UserID,
		CalendarID: plan.CalendarID,
		Name:       plan.Name,
		Memo:       plan.Memo,
		Color:      string(plan.Color),
		Private:    plan.Private,
		Shares:     plan.Shares,
		IsAllDay:   plan.Period.IsAllDay,
		Begin:      plan.Period.Begin.Unix(),
		End:        plan.Period.End.Unix(),
	})
}

func (e *planEndpoint) UnsheduleHandler(w http.ResponseWriter, r *http.Request) {
	req := PlanContent{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	userID := r.Context().Value(cctx.UserIDKey).(string)
	err := e.service.Unschedule(r.Context(), userID, req.CalendarID, vars["id"])
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

func (e *planEndpoint) ResheduleHandler(w http.ResponseWriter, r *http.Request) {
	req := PlanContent{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	color, err := model.ConvertToColor(req.Color)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)

	planPram := model.Plan{
		ID:         vars["id"],
		CalendarID: req.CalendarID,
		Name:       req.Name,
		Memo:       req.Memo,
		Color:      color,
		Private:    req.Private,
		Shares:     req.Shares,
		Period: model.Period{
			IsAllDay: req.IsAllDay,
			Begin:    time.Unix(req.Begin, 0),
			End:      time.Unix(req.End, 0),
		},
	}

	userID := r.Context().Value(cctx.UserIDKey).(string)
	_, err = e.service.Reschedule(r.Context(), userID, planPram)
	if errors.Is(err, cerror.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if errors.Is(err, cerror.ErrInvalidContent) {
		w.WriteHeader(http.StatusBadRequest)
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

func NewPlanRouter(r *mux.Router, calService cs.Service, authService as.Service) {
	e := planEndpoint{calService}
	r.Use(middlewares.ReqIDMiddleware)
	r.Use(middlewares.ResponseHeaderMiddleware)
	r.Use(middlewares.AuthorizationMiddleware(authService))
	r.Use(userCheckerMiddleware(calService))
	r.HandleFunc("", e.ScheduleHandler).Methods(http.MethodPost)
	r.HandleFunc("/{id}", e.UnsheduleHandler).Methods(http.MethodDelete)
	r.HandleFunc("/{id}", e.ResheduleHandler).Methods(http.MethodPatch)
}
