package rest

import (
	"encoding/json"
	"net/http"

	cs "github.com/x-color/calendar/service/calendar"
)

type planContent struct {
	ID         string `json:"id"`
	CalendarID string `json:"calendar_id"`
	UserID     string `json:"user_id"`
	Name       string `json:"name"`
	Memo       string `json:"memo"`
	Color      string `json:"color"`
	Private    bool   `json:"private"`
	IsAllDay   bool   `json:"is_all_day"`
	Begin      string `json:"begin"`
	End        string `json:"end"`
}

type PlanEndpoint struct {
	service cs.Service
}

func (e *PlanEndpoint) scheduleHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(msgContent{Msg: "schedule plan"})
}

func (e *PlanEndpoint) unsheduleHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(msgContent{"unshedule plan"})
}

func (e *PlanEndpoint) resheduleHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(msgContent{"reshedule plan"})
}
