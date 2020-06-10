package calendar_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	. "github.com/x-color/calendar/app/rest/calendar"
	"github.com/x-color/calendar/app/rest/testutils"
	as "github.com/x-color/calendar/auth/service"
	mcal "github.com/x-color/calendar/calendar/model"
	cs "github.com/x-color/calendar/calendar/service"
)

func makeCalendar(calRepo cs.Repogitory, ownerID string, shares ...string) mcal.Calendar {
	cal := mcal.Calendar{
		ID:     uuid.New().String(),
		Name:   "My plans",
		UserID: ownerID,
		Color:  mcal.RED,
		Plans:  []mcal.Plan{},
		Shares: append(shares, ownerID),
	}
	calData := cs.CalendarData{
		ID:     cal.ID,
		Name:   "My plans",
		UserID: ownerID,
		Color:  "red",
		Shares: append(shares, ownerID),
	}
	calRepo.Calendar().Create(context.Background(), calData)
	return cal
}

func makePlan(calRepo cs.Repogitory, ownerID, calendarID string, shares ...string) mcal.Plan {
	plan := mcal.Plan{
		ID:         uuid.New().String(),
		CalendarID: calendarID,
		UserID:     ownerID,
		Name:       "My plan",
		Color:      "red",
		Shares:     append(shares, calendarID),
		Period: mcal.Period{
			IsAllDay: true,
			Begin:    time.Now(),
			End:      time.Now(),
		},
	}
	planData := cs.PlanData{
		ID:         plan.ID,
		CalendarID: calendarID,
		UserID:     ownerID,
		Name:       "My plan",
		Color:      "red",
		Shares:     append(shares, calendarID),
		IsAllDay:   true,
		Begin:      plan.Period.Begin.Unix(),
		End:        plan.Period.End.Unix(),
	}
	calRepo.Plan().Create(context.Background(), planData)
	return plan
}

func makePrivatePlan(calRepo cs.Repogitory, ownerID, calendarID string, shares ...string) mcal.Plan {
	plan := mcal.Plan{
		ID:         uuid.New().String(),
		CalendarID: calendarID,
		UserID:     ownerID,
		Name:       "My plan",
		Color:      "red",
		Private:    true,
		Shares:     append(shares, calendarID),
		Period: mcal.Period{
			IsAllDay: true,
			Begin:    time.Now(),
			End:      time.Now(),
		},
	}
	planData := cs.PlanData{
		ID:         plan.ID,
		CalendarID: calendarID,
		UserID:     ownerID,
		Name:       "My plan",
		Color:      "red",
		Private:    true,
		Shares:     append(shares, calendarID),
		IsAllDay:   true,
		Begin:      plan.Period.Begin.Unix(),
		End:        plan.Period.End.Unix(),
	}
	calRepo.Plan().Create(context.Background(), planData)
	return plan
}

func TestNewPlanRouter_Authoraization(t *testing.T) {
	authRepo := testutils.NewAuthRepo()
	userID, sessionID := testutils.MakeSession(authRepo)
	calRepo := testutils.NewCalRepo()
	calRepo.User().Create(context.Background(), cs.UserData{userID})
	cal := makeCalendar(calRepo, userID)

	l := testutils.NewLogger()
	authService := as.NewService(authRepo, l)
	calendarService := cs.NewService(calRepo, l)
	r := mux.NewRouter()
	NewPlanRouter(r.PathPrefix("/plans").Subrouter(), calendarService, authService)

	cookie := http.Cookie{
		Name:  "session_id",
		Value: sessionID,
	}

	testcases := []struct {
		name   string
		cookie *http.Cookie
		body   map[string]interface{}
		code   int
		res    PlanContent
	}{
		{
			name:   "no cookie",
			cookie: nil,
			body: map[string]interface{}{
				"calendar_id": cal.ID,
				"name":        "all day plan",
				"memo":        "sample text",
				"color":       "red",
				"private":     true,
				"shares":      []interface{}{cal.ID},
				"is_all_day":  true,
				"begin":       time.Date(2020, 4, 1, 0, 0, 0, 0, time.Local).Unix(),
				"end":         time.Date(2020, 4, 1, 0, 0, 0, 0, time.Local).Unix(),
			},
			code: http.StatusUnauthorized,
		},
		{
			name:   "invalid cookie",
			cookie: &http.Cookie{Name: "session_id", Value: uuid.New().String()},
			body: map[string]interface{}{
				"calendar_id": cal.ID,
				"name":        "all day plan",
				"memo":        "sample text",
				"color":       "red",
				"private":     true,
				"shares":      []interface{}{cal.ID},
				"is_all_day":  true,
				"begin":       time.Date(2020, 4, 1, 0, 0, 0, 0, time.Local).Unix(),
				"end":         time.Date(2020, 4, 1, 0, 0, 0, 0, time.Local).Unix(),
			},
			code: http.StatusUnauthorized,
		},
		{
			name:   "valid cookie",
			cookie: &cookie,
			body: map[string]interface{}{
				"calendar_id": cal.ID,
				"name":        "all day plan",
				"memo":        "sample text",
				"color":       "red",
				"private":     true,
				"shares":      []interface{}{cal.ID},
				"is_all_day":  true,
				"begin":       time.Date(2020, 4, 1, 0, 0, 0, 0, time.Local).Unix(),
				"end":         time.Date(2020, 4, 1, 0, 0, 0, 0, time.Local).Unix(),
			},
			code: http.StatusOK,
			res: PlanContent{
				ID:         "",
				UserID:     userID,
				CalendarID: cal.ID,
				Name:       "all day plan",
				Memo:       "sample text",
				Color:      "red",
				Private:    true,
				Shares:     []string{cal.ID},
				IsAllDay:   true,
				Begin:      time.Date(2020, 4, 1, 0, 0, 0, 0, time.Local).Unix(),
				End:        time.Date(2020, 4, 1, 0, 0, 0, 0, time.Local).Unix(),
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.body)
			req := httptest.NewRequest(http.MethodPost, "/plans", bytes.NewBuffer(body))
			if tc.cookie != nil {
				req.AddCookie(tc.cookie)
			}
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			if rec.Code != tc.code {
				t.Errorf("status code: want %v but %v", tc.code, rec.Code)
			}

			var actual PlanContent
			if len(rec.Body.Bytes()) > 0 {
				if err := json.Unmarshal(rec.Body.Bytes(), &actual); err != nil {
					t.Errorf("invalid response body: %v", rec.Body.String())
				}
			}
			expected := tc.res

			if d := cmp.Diff(expected, actual, cmpopts.IgnoreFields(PlanContent{}, "ID")); d != "" {
				t.Errorf("invalid response body: \n%v", d)
			}
		})
	}
}

func TestNewPlanRouter_UserRegistrationChecker(t *testing.T) {
	authRepo := testutils.NewAuthRepo()
	userID, sessionID := testutils.MakeSession(authRepo)
	_, sessionID2 := testutils.MakeSession(authRepo)
	calRepo := testutils.NewCalRepo()
	calRepo.User().Create(context.Background(), cs.UserData{userID})
	cal := makeCalendar(calRepo, userID)

	l := testutils.NewLogger()
	authService := as.NewService(authRepo, l)
	calendarService := cs.NewService(calRepo, l)
	r := mux.NewRouter()
	NewPlanRouter(r.PathPrefix("/plans").Subrouter(), calendarService, authService)

	cookie := http.Cookie{
		Name:  "session_id",
		Value: sessionID,
	}

	cookie2 := http.Cookie{
		Name:  "session_id",
		Value: sessionID2,
	}

	testcases := []struct {
		name   string
		cookie *http.Cookie
		body   map[string]interface{}
		code   int
		res    PlanContent
	}{
		{
			name:   "no registered user",
			cookie: &cookie2,
			body:   map[string]interface{}{"name": "test", "color": "red"},
			code:   http.StatusForbidden,
		},
		{
			name:   "registered user",
			cookie: &cookie,
			body: map[string]interface{}{
				"calendar_id": cal.ID,
				"name":        "plan",
				"memo":        "sample text",
				"color":       "red",
				"private":     true,
				"shares":      []interface{}{cal.ID},
				"begin":       time.Date(2020, 4, 1, 9, 0, 0, 0, time.Local).Unix(),
				"end":         time.Date(2020, 4, 1, 18, 0, 0, 0, time.Local).Unix(),
			},
			code: http.StatusOK,
			res: PlanContent{
				ID:         "",
				UserID:     userID,
				CalendarID: cal.ID,
				Name:       "plan",
				Memo:       "sample text",
				Color:      "red",
				Private:    true,
				Shares:     []string{cal.ID},
				IsAllDay:   false,
				Begin:      time.Date(2020, 4, 1, 9, 0, 0, 0, time.Local).Unix(),
				End:        time.Date(2020, 4, 1, 18, 0, 0, 0, time.Local).Unix(),
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.body)
			req := httptest.NewRequest(http.MethodPost, "/plans", bytes.NewBuffer(body))
			if tc.cookie != nil {
				req.AddCookie(tc.cookie)
			}
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			if rec.Code != tc.code {
				t.Errorf("status code: want %v but %v", tc.code, rec.Code)
			}

			var actual PlanContent
			if len(rec.Body.Bytes()) > 0 {
				if err := json.Unmarshal(rec.Body.Bytes(), &actual); err != nil {
					t.Errorf("invalid response body: %v", rec.Body.String())
				}
			}
			expected := tc.res

			if d := cmp.Diff(expected, actual, cmpopts.IgnoreFields(PlanContent{}, "ID")); d != "" {
				t.Errorf("invalid response body: \n%v", d)
			}
		})
	}
}

func TestNewPlanRouter_Shedule(t *testing.T) {
	authRepo := testutils.NewAuthRepo()
	userID, sessionID := testutils.MakeSession(authRepo)
	otherID := uuid.New().String()
	calRepo := testutils.NewCalRepo()
	calRepo.User().Create(context.Background(), cs.UserData{userID})
	cal := makeCalendar(calRepo, userID)
	otherCal := makeCalendar(calRepo, otherID)
	sharedCal := makeCalendar(calRepo, otherID, userID)
	l := testutils.NewLogger()
	authService := as.NewService(authRepo, l)
	calendarService := cs.NewService(calRepo, l)
	r := mux.NewRouter()
	NewPlanRouter(r.PathPrefix("/plans").Subrouter(), calendarService, authService)

	cookie := http.Cookie{
		Name:  "session_id",
		Value: sessionID,
	}

	testcases := []struct {
		name   string
		cookie *http.Cookie
		body   map[string]interface{}
		code   int
		res    PlanContent
	}{
		{
			name:   "invalid calendar",
			cookie: &cookie,
			body: map[string]interface{}{
				"calendar_id": uuid.New().String(),
				"name":        "plan",
				"shares":      []interface{}{uuid.New().String()},
				"is_all_day":  true,
				"begin":       time.Date(2020, 4, 1, 0, 0, 0, 0, time.Local).Unix(),
				"end":         time.Date(2020, 4, 1, 0, 0, 0, 0, time.Local).Unix(),
			},
			code: http.StatusBadRequest,
		},
		{
			name:   "do not permit to access calendar",
			cookie: &cookie,
			body: map[string]interface{}{
				"calendar_id": otherCal.ID,
				"name":        "all day plan",
				"memo":        "sample text",
				"color":       "red",
				"private":     true,
				"shares":      []interface{}{otherCal.ID},
				"is_all_day":  true,
				"begin":       time.Date(2020, 4, 1, 0, 0, 0, 0, time.Local).Unix(),
				"end":         time.Date(2020, 4, 1, 0, 0, 0, 0, time.Local).Unix(),
			},
			code: http.StatusBadRequest,
		},
		{
			name:   "not shared calendar",
			cookie: &cookie,
			body: map[string]interface{}{
				"calendar_id": cal.ID,
				"name":        "all day plan",
				"memo":        "sample text",
				"color":       "red",
				"private":     true,
				"shares":      []interface{}{cal.ID, otherCal.ID},
				"is_all_day":  true,
				"begin":       time.Date(2020, 4, 1, 0, 0, 0, 0, time.Local).Unix(),
				"end":         time.Date(2020, 4, 1, 0, 0, 0, 0, time.Local).Unix(),
			},
			code: http.StatusBadRequest,
		},
		{
			name:   "invalid content",
			cookie: &cookie,
			body: map[string]interface{}{
				"calendar_id": cal.ID,
				"shares":      []interface{}{cal.ID},
				"is_all_day":  true,
				"begin":       time.Date(2020, 4, 1, 0, 0, 0, 0, time.Local).Unix(),
				"end":         time.Date(2020, 4, 1, 0, 0, 0, 0, time.Local).Unix(),
			},
			code: http.StatusBadRequest,
		},
		{
			name:   "begin is after end in all day plan",
			cookie: &cookie,
			body: map[string]interface{}{
				"calendar_id": cal.ID,
				"name":        "all day plan",
				"memo":        "sample text",
				"color":       "red",
				"private":     true,
				"shares":      []interface{}{cal.ID, uuid.New().String()},
				"is_all_day":  true,
				"begin":       time.Date(2020, 4, 1, 1, 0, 0, 0, time.Local).Unix(),
				"end":         time.Date(2020, 4, 1, 0, 0, 0, 0, time.Local).Unix(),
			},
			code: http.StatusBadRequest,
		},
		{
			name:   "invalid calendar id in shares",
			cookie: &cookie,
			body: map[string]interface{}{
				"calendar_id": cal.ID,
				"name":        "all day plan",
				"memo":        "sample text",
				"color":       "red",
				"private":     true,
				"shares":      []interface{}{cal.ID, uuid.New().String()},
				"is_all_day":  true,
				"begin":       time.Date(2020, 4, 1, 0, 0, 0, 0, time.Local).Unix(),
				"end":         time.Date(2020, 4, 1, 0, 0, 0, 0, time.Local).Unix(),
			},
			code: http.StatusBadRequest,
		},
		{
			name:   "shedule plan",
			cookie: &cookie,
			body: map[string]interface{}{
				"calendar_id": cal.ID,
				"name":        "plan",
				"memo":        "sample text",
				"color":       "red",
				"private":     true,
				"shares":      []interface{}{cal.ID, sharedCal.ID},
				"begin":       time.Date(2020, 4, 1, 9, 0, 0, 0, time.Local).Unix(),
				"end":         time.Date(2020, 4, 1, 18, 0, 0, 0, time.Local).Unix(),
			},
			code: http.StatusOK,
			res: PlanContent{
				ID:         "",
				UserID:     userID,
				CalendarID: cal.ID,
				Name:       "plan",
				Memo:       "sample text",
				Color:      "red",
				Private:    true,
				Shares:     []string{cal.ID, sharedCal.ID},
				IsAllDay:   false,
				Begin:      time.Date(2020, 4, 1, 9, 0, 0, 0, time.Local).Unix(),
				End:        time.Date(2020, 4, 1, 18, 0, 0, 0, time.Local).Unix(),
			},
		},
		{
			name:   "shedule all day plan",
			cookie: &cookie,
			body: map[string]interface{}{
				"calendar_id": cal.ID,
				"name":        "all day plan",
				"memo":        "sample text",
				"color":       "red",
				"private":     true,
				"shares":      []interface{}{cal.ID},
				"is_all_day":  true,
				"begin":       time.Date(2020, 4, 1, 0, 0, 0, 0, time.Local).Unix(),
				"end":         time.Date(2020, 4, 1, 0, 0, 0, 0, time.Local).Unix(),
			},
			code: http.StatusOK,
			res: PlanContent{
				ID:         "",
				UserID:     userID,
				CalendarID: cal.ID,
				Name:       "all day plan",
				Memo:       "sample text",
				Color:      "red",
				Private:    true,
				Shares:     []string{cal.ID},
				IsAllDay:   true,
				Begin:      time.Date(2020, 4, 1, 0, 0, 0, 0, time.Local).Unix(),
				End:        time.Date(2020, 4, 1, 0, 0, 0, 0, time.Local).Unix(),
			},
		},
		{
			name:   "shedule plan in shared calendar",
			cookie: &cookie,
			body: map[string]interface{}{
				"id":          "",
				"calendar_id": sharedCal.ID,
				"name":        "plan",
				"memo":        "sample text",
				"color":       "red",
				"shares":      []interface{}{sharedCal.ID},
				"is_all_day":  true,
				"begin":       time.Date(2020, 4, 1, 0, 0, 0, 0, time.Local).Unix(),
				"end":         time.Date(2020, 4, 1, 0, 0, 0, 0, time.Local).Unix(),
			},
			code: http.StatusOK,
			res: PlanContent{
				ID:         "",
				UserID:     userID,
				CalendarID: sharedCal.ID,
				Name:       "plan",
				Memo:       "sample text",
				Color:      "red",
				Private:    false,
				Shares:     []string{sharedCal.ID},
				IsAllDay:   true,
				Begin:      time.Date(2020, 4, 1, 0, 0, 0, 0, time.Local).Unix(),
				End:        time.Date(2020, 4, 1, 0, 0, 0, 0, time.Local).Unix(),
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.body)
			req := httptest.NewRequest(http.MethodPost, "/plans", bytes.NewBuffer(body))
			if tc.cookie != nil {
				req.AddCookie(tc.cookie)
			}
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			if rec.Code != tc.code {
				t.Errorf("status code: want %v but %v", tc.code, rec.Code)
			}

			var actual PlanContent
			if len(rec.Body.Bytes()) > 0 {
				if err := json.Unmarshal(rec.Body.Bytes(), &actual); err != nil {
					t.Errorf("invalid response body: %v", rec.Body.String())
				}
			}
			expected := tc.res

			if d := cmp.Diff(expected, actual, cmpopts.IgnoreFields(PlanContent{}, "ID")); d != "" {
				t.Errorf("invalid response body: \n%v", d)
			}
		})
	}
}

func TestNewPlanRouter_Unshedule(t *testing.T) {
	authRepo := testutils.NewAuthRepo()
	userID, sessionID := testutils.MakeSession(authRepo)
	otherID := uuid.New().String()
	calRepo := testutils.NewCalRepo()
	calRepo.User().Create(context.Background(), cs.UserData{userID})
	cal := makeCalendar(calRepo, userID, otherID)
	sharedCal := makeCalendar(calRepo, otherID, userID)
	otherCal := makeCalendar(calRepo, otherID)

	plan := makePlan(calRepo, userID, cal.ID)
	sharedPlan := makePlan(calRepo, userID, cal.ID, sharedCal.ID)
	otherPlan := makePlan(calRepo, otherID, otherCal.ID)

	l := testutils.NewLogger()
	authService := as.NewService(authRepo, l)
	calendarService := cs.NewService(calRepo, l)
	r := mux.NewRouter()
	NewPlanRouter(r.PathPrefix("/plans").Subrouter(), calendarService, authService)

	cookie := http.Cookie{
		Name:  "session_id",
		Value: sessionID,
	}

	testcases := []struct {
		name   string
		cookie *http.Cookie
		planID string
		body   map[string]string
		code   int
	}{
		{
			name:   "invalid plan id",
			cookie: &cookie,
			planID: uuid.New().String(),
			body:   map[string]string{"calendar_id": plan.CalendarID},
			code:   http.StatusNotFound,
		},
		{
			name:   "invalid contents",
			cookie: &cookie,
			planID: plan.ID,
			body:   map[string]string{},
			code:   http.StatusBadRequest,
		},
		{
			name:   "invalid calendar id",
			cookie: &cookie,
			planID: plan.ID,
			body:   map[string]string{"calendar_id": uuid.New().String()},
			code:   http.StatusBadRequest,
		},
		{
			name:   "other calendar id",
			cookie: &cookie,
			planID: plan.ID,
			body:   map[string]string{"calendar_id": otherCal.ID},
			code:   http.StatusBadRequest,
		},
		{
			name:   "do not permit to access plan",
			cookie: &cookie,
			planID: otherPlan.ID,
			body:   map[string]string{"calendar_id": otherPlan.CalendarID},
			code:   http.StatusForbidden,
		},
		{
			name:   "unshedule plan",
			cookie: &cookie,
			planID: plan.ID,
			body:   map[string]string{"calendar_id": plan.CalendarID},
			code:   http.StatusNoContent,
		},
		{
			name:   "unshedule shared plan",
			cookie: &cookie,
			planID: sharedPlan.ID,
			body:   map[string]string{"calendar_id": sharedPlan.CalendarID},
			code:   http.StatusNoContent,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.body)
			req := httptest.NewRequest(http.MethodDelete, "/plans/"+tc.planID, bytes.NewBuffer(body))
			if tc.cookie != nil {
				req.AddCookie(tc.cookie)
			}
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			if rec.Code != tc.code {
				t.Errorf("status code: want %v but %v", tc.code, rec.Code)
			}
		})
	}
}

func TestNewPlanRouter_Reschedule(t *testing.T) {
	authRepo := testutils.NewAuthRepo()
	userID, sessionID := testutils.MakeSession(authRepo)
	otherID, otherSessionID := testutils.MakeSession(authRepo)
	calRepo := testutils.NewCalRepo()
	calRepo.User().Create(context.Background(), cs.UserData{userID})
	cal := makeCalendar(calRepo, userID, otherID)
	sharedCal := makeCalendar(calRepo, otherID, userID)
	otherCal := makeCalendar(calRepo, otherID)

	plan := makePlan(calRepo, userID, cal.ID)
	sharedPlan := makePlan(calRepo, userID, cal.ID, sharedCal.ID)
	otherPlan := makePlan(calRepo, otherID, otherCal.ID)

	l := testutils.NewLogger()
	authService := as.NewService(authRepo, l)
	calendarService := cs.NewService(calRepo, l)
	r := mux.NewRouter()
	NewPlanRouter(r.PathPrefix("/plans").Subrouter(), calendarService, authService)

	cookie := http.Cookie{
		Name:  "session_id",
		Value: sessionID,
	}

	testcases := []struct {
		name   string
		cookie *http.Cookie
		planID string
		body   map[string]interface{}
		code   int
	}{
		{
			name:   "invalid plan id",
			cookie: &cookie,
			planID: uuid.New().String(),
			body: map[string]interface{}{
				"calendar_id": cal.ID,
				"name":        "renamed",
				"memo":        "edited text",
				"color":       "yellow",
				"private":     false,
				"shares":      []interface{}{cal.ID},
				"begin":       time.Date(2020, 5, 1, 9, 0, 0, 0, time.Local).Unix(),
				"end":         time.Date(2020, 5, 1, 18, 0, 0, 0, time.Local).Unix(),
			},
			code: http.StatusNotFound,
		},
		{
			name:   "invalid contents",
			cookie: &cookie,
			planID: plan.ID,
			body:   map[string]interface{}{},
			code:   http.StatusBadRequest,
		},
		{
			name:   "invalid calendar id",
			cookie: &cookie,
			planID: plan.ID,
			body: map[string]interface{}{
				"calendar_id": uuid.New().String(),
				"name":        "renamed",
				"memo":        "edited text",
				"color":       "yellow",
				"private":     false,
				"shares":      []interface{}{cal.ID},
				"begin":       time.Date(2020, 5, 1, 9, 0, 0, 0, time.Local).Unix(),
				"end":         time.Date(2020, 5, 1, 18, 0, 0, 0, time.Local).Unix(),
			},
			code: http.StatusBadRequest,
		},
		{
			name:   "other calendar id",
			cookie: &cookie,
			planID: plan.ID,
			body: map[string]interface{}{
				"calendar_id": otherCal.ID,
				"name":        "renamed",
				"memo":        "edited text",
				"color":       "yellow",
				"private":     false,
				"shares":      []interface{}{cal.ID},
				"begin":       time.Date(2020, 5, 1, 9, 0, 0, 0, time.Local).Unix(),
				"end":         time.Date(2020, 5, 1, 18, 0, 0, 0, time.Local).Unix(),
			},
			code: http.StatusBadRequest,
		},
		{
			name:   "do not permit to access other plan",
			cookie: &cookie,
			planID: otherPlan.ID,
			body: map[string]interface{}{
				"calendar_id": otherCal.ID,
				"name":        "renamed",
				"memo":        "edited text",
				"color":       "yellow",
				"private":     false,
				"shares":      []interface{}{cal.ID},
				"begin":       time.Date(2020, 5, 1, 9, 0, 0, 0, time.Local).Unix(),
				"end":         time.Date(2020, 5, 1, 18, 0, 0, 0, time.Local).Unix(),
			},
			code: http.StatusForbidden,
		},
		{
			name:   "do not permit to reshcedule shared plan",
			cookie: &http.Cookie{Name: "session_id", Value: otherSessionID},
			planID: sharedPlan.ID,
			body: map[string]interface{}{
				"calendar_id": sharedCal.ID,
				"name":        "renamed",
				"memo":        "edited text",
				"color":       "yellow",
				"private":     false,
				"shares":      []interface{}{cal.ID, sharedCal.ID},
				"begin":       time.Date(2020, 5, 1, 9, 0, 0, 0, time.Local).Unix(),
				"end":         time.Date(2020, 5, 1, 18, 0, 0, 0, time.Local).Unix(),
			},
			code: http.StatusForbidden,
		},
		{
			name:   "reshedule plan",
			cookie: &cookie,
			planID: plan.ID,
			body: map[string]interface{}{
				"calendar_id": cal.ID,
				"name":        "renamed",
				"memo":        "edited text",
				"color":       "yellow",
				"private":     false,
				"shares":      []interface{}{cal.ID},
				"begin":       time.Date(2020, 5, 1, 9, 0, 0, 0, time.Local).Unix(),
				"end":         time.Date(2020, 5, 1, 18, 0, 0, 0, time.Local).Unix(),
			},
			code: http.StatusNoContent,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.body)
			req := httptest.NewRequest(http.MethodPatch, "/plans/"+tc.planID, bytes.NewBuffer(body))
			if tc.cookie != nil {
				req.AddCookie(tc.cookie)
			}
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			if rec.Code != tc.code {
				t.Errorf("status code: want %v but %v", tc.code, rec.Code)
			}
		})
	}
}
