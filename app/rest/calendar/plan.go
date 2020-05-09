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

type planContent struct {
	ID         string   `json:"id"`
	CalendarID string   `json:"calendar_id"`
	Name       string   `json:"name"`
	Memo       string   `json:"memo"`
	Color      string   `json:"color"`
	Private    bool     `json:"private"`
	Shares     []string `json:"shares"`
	IsAllDay   bool     `json:"is_all_day"`
	Begin      int      `json:"begin"`
	End        int      `json:"end"`
}

type planEndpoint struct {
	service service.Service
}

func (e *planEndpoint) ScheduleHandler(w http.ResponseWriter, r *http.Request) {
	req := planContent{}
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

	planPram := model.Plan{
		CalendarID: req.CalendarID,
		Name:       req.Name,
		Memo:       req.Memo,
		Color:      color,
		Private:    req.Private,
		Shares:     req.Shares,
		Period: model.Period{
			IsAllDay: req.IsAllDay,
			Begin:    time.Unix(int64(req.Begin), 0),
			End:      time.Unix(int64(req.End), 0),
		},
	}

	userID := r.Context().Value(cctx.UserIDKey).(string)
	plan, err := e.service.Schedule(r.Context(), userID, planPram)
	if errors.Is(err, cerror.ErrInvalidContent) || errors.Is(err, cerror.ErrAuthorization) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msgContent{"bad contents"})
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msgContent{"internal server error"})
		return
	}

	json.NewEncoder(w).Encode(planContent{
		ID:         plan.ID,
		CalendarID: plan.CalendarID,
		Name:       plan.Name,
		Memo:       plan.Memo,
		Color:      string(plan.Color),
		Private:    plan.Private,
		Shares:     plan.Shares,
		IsAllDay:   plan.Period.IsAllDay,
		Begin:      int(plan.Period.Begin.Unix()),
		End:        int(plan.Period.End.Unix()),
	})
}

func (e *planEndpoint) UnsheduleHandler(w http.ResponseWriter, r *http.Request) {
	req := planContent{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msgContent{"bad contents"})
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
	json.NewEncoder(w).Encode(msgContent{"reshedule plan"})
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
