package calendar_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
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
		Period:     mcal.AllDay,
	}
	planData := cs.PlanData{
		ID:         plan.ID,
		CalendarID: calendarID,
		UserID:     ownerID,
		Name:       "My plan",
		Color:      "red",
		Shares:     append(shares, calendarID),
		IsAllDay:   true,
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
		res    map[string]interface{}
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
			},
			code: http.StatusUnauthorized,
			res:  map[string]interface{}{"message": "unauthorization"},
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
			},
			code: http.StatusUnauthorized,
			res:  map[string]interface{}{"message": "unauthorization"},
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
			},
			code: http.StatusOK,
			res: map[string]interface{}{
				"id":          "",
				"calendar_id": cal.ID,
				"name":        "all day plan",
				"memo":        "sample text",
				"color":       "red",
				"private":     true,
				"shares":      []interface{}{cal.ID},
				"is_all_day":  true,
				"begin":       float64(0),
				"end":         float64(0),
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

			var actual map[string]interface{}
			if err := json.Unmarshal(rec.Body.Bytes(), &actual); err != nil {
				t.Errorf("invalid response body: %v", rec.Body.String())
			}
			expected := tc.res

			if d := cmp.Diff(expected, actual, testutils.IgnoreKey("id")); d != "" {
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
		res    map[string]interface{}
	}{
		{
			name:   "no registered user",
			cookie: &cookie2,
			body:   map[string]interface{}{"name": "test", "color": "red"},
			code:   http.StatusForbidden,
			res:    map[string]interface{}{"message": "forbidden"},
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
				"begin":       1577836800,
				"end":         1577869200,
			},
			code: http.StatusOK,
			res: map[string]interface{}{
				"id":          "",
				"calendar_id": cal.ID,
				"name":        "plan",
				"memo":        "sample text",
				"color":       "red",
				"private":     true,
				"shares":      []interface{}{cal.ID},
				"is_all_day":  false,
				"begin":       float64(1577836800),
				"end":         float64(1577869200),
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

			var actual map[string]interface{}
			if err := json.Unmarshal(rec.Body.Bytes(), &actual); err != nil {
				t.Errorf("invalid response body: %v", rec.Body.String())
			}
			expected := tc.res

			if d := cmp.Diff(expected, actual, testutils.IgnoreKey("id")); d != "" {
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
		res    map[string]interface{}
	}{
		{
			name:   "invalid calendar",
			cookie: &cookie,
			body: map[string]interface{}{
				"calendar_id": uuid.New().String(),
				"name":        "plan",
				"shares":      []interface{}{uuid.New().String()},
				"is_all_day":  true,
			},
			code: http.StatusBadRequest,
			res:  map[string]interface{}{"message": "bad contents"},
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
			},
			code: http.StatusBadRequest,
			res:  map[string]interface{}{"message": "bad contents"},
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
			},
			code: http.StatusBadRequest,
			res:  map[string]interface{}{"message": "bad contents"},
		},
		{
			name:   "invalid content",
			cookie: &cookie,
			body: map[string]interface{}{
				"calendar_id": cal.ID,
				"shares":      []interface{}{cal.ID},
				"is_all_day":  true,
			},
			code: http.StatusBadRequest,
			res:  map[string]interface{}{"message": "bad contents"},
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
			},
			code: http.StatusBadRequest,
			res:  map[string]interface{}{"message": "bad contents"},
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
				"begin":       1577836800,
				"end":         1577869200,
			},
			code: http.StatusOK,
			res: map[string]interface{}{
				"id":          "",
				"calendar_id": cal.ID,
				"name":        "plan",
				"memo":        "sample text",
				"color":       "red",
				"private":     true,
				"shares":      []interface{}{cal.ID, sharedCal.ID},
				"is_all_day":  false,
				"begin":       float64(1577836800),
				"end":         float64(1577869200),
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
			},
			code: http.StatusOK,
			res: map[string]interface{}{
				"id":          "",
				"calendar_id": cal.ID,
				"name":        "all day plan",
				"memo":        "sample text",
				"color":       "red",
				"private":     true,
				"shares":      []interface{}{cal.ID},
				"is_all_day":  true,
				"begin":       float64(0),
				"end":         float64(0),
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
			},
			code: http.StatusOK,
			res: map[string]interface{}{
				"id":          "",
				"calendar_id": sharedCal.ID,
				"name":        "plan",
				"memo":        "sample text",
				"color":       "red",
				"private":     false,
				"shares":      []interface{}{sharedCal.ID},
				"is_all_day":  true,
				"begin":       float64(0),
				"end":         float64(0),
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

			var actual map[string]interface{}
			if err := json.Unmarshal(rec.Body.Bytes(), &actual); err != nil {
				t.Errorf("invalid response body: %v", rec.Body.String())
			}
			expected := tc.res

			if d := cmp.Diff(expected, actual, testutils.IgnoreKey("id")); d != "" {
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
	sharedPlan := makePlan(calRepo, otherID, sharedCal.ID)
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
		code   int
	}{
		{
			name:   "invalid plan id",
			cookie: &cookie,
			planID: uuid.New().String(),
			code:   http.StatusNotFound,
		},
		{
			name:   "do not permit to access plan",
			cookie: &cookie,
			planID: otherPlan.ID,
			code:   http.StatusForbidden,
		},
		{
			name:   "unshedule plan",
			cookie: &cookie,
			planID: plan.ID,
			code:   http.StatusNoContent,
		},
		{
			name:   "unshedule shared plan",
			cookie: &cookie,
			planID: sharedPlan.ID,
			code:   http.StatusNoContent,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/plans/"+tc.planID, nil)
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
