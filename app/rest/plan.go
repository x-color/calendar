package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/x-color/calendar/model/calendar"
	cctx "github.com/x-color/calendar/model/ctx"
	cerror "github.com/x-color/calendar/model/error"
	cs "github.com/x-color/calendar/service/calendar"
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

type PlanEndpoint struct {
	service cs.Service
}

func (e *PlanEndpoint) scheduleHandler(w http.ResponseWriter, r *http.Request) {
	req := planContent{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msgContent{"bad contents"})
		return
	}

	color, err := calendar.ConvertToColor(req.Color)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msgContent{"bad contents"})
		return
	}

	planPram := calendar.Plan{
		CalendarID: req.CalendarID,
		Name:       req.Name,
		Memo:       req.Memo,
		Color:      color,
		Private:    req.Private,
		Shares:     req.Shares,
		Period: calendar.Period{
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

func (e *PlanEndpoint) unsheduleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := r.Context().Value(cctx.UserIDKey).(string)
	err := e.service.Unschedule(r.Context(), userID, vars["id"])
	if errors.Is(err, cerror.ErrInvalidContent) {
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

func (e *PlanEndpoint) resheduleHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(msgContent{"reshedule plan"})
}
