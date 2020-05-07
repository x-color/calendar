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
	cs "github.com/x-color/calendar/calendar/service"
)

func TestNewCalendarRouter_Authoraization(t *testing.T) {
	authRepo := testutils.NewAuthRepo()
	userID, sessionID := testutils.MakeSession(authRepo)
	calRepo := testutils.NewCalRepo()
	calRepo.CalUser().Create(context.Background(), cs.UserData{userID})

	l := testutils.NewLogger()
	authService := as.NewService(authRepo, l)
	calendarService := cs.NewService(calRepo, l)
	r := mux.NewRouter()
	NewCalendarRouter(r.PathPrefix("/calendars").Subrouter(), calendarService, authService)

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
			body:   map[string]interface{}{"name": "My plans", "color": "red"},
			code:   http.StatusUnauthorized,
			res:    map[string]interface{}{"message": "unauthorization"},
		},
		{
			name:   "invalid cookie",
			cookie: &http.Cookie{Name: "session_id", Value: uuid.New().String()},
			body:   map[string]interface{}{"name": "My plans", "color": "red"},
			code:   http.StatusUnauthorized,
			res:    map[string]interface{}{"message": "unauthorization"},
		},
		{
			name:   "valid cookie",
			cookie: &cookie,
			body:   map[string]interface{}{"name": "My plans", "color": "red"},
			code:   http.StatusOK,
			res:    map[string]interface{}{"id": "", "name": "My plans", "color": "red", "shares": []interface{}{userID}, "plans": nil},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.body)
			req := httptest.NewRequest(http.MethodPost, "/calendars", bytes.NewBuffer(body))
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

func TestNewCalendarRouter_UserRegistrationChecker(t *testing.T) {
	authRepo := testutils.NewAuthRepo()
	userID, sessionID := testutils.MakeSession(authRepo)
	_, sessionID2 := testutils.MakeSession(authRepo)
	calRepo := testutils.NewCalRepo()
	calRepo.CalUser().Create(context.Background(), cs.UserData{userID})

	l := testutils.NewLogger()
	authService := as.NewService(authRepo, l)
	calendarService := cs.NewService(calRepo, l)
	r := mux.NewRouter()
	NewCalendarRouter(r.PathPrefix("/calendars").Subrouter(), calendarService, authService)

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
			body:   map[string]interface{}{"name": "test", "color": "red"},
			code:   http.StatusOK,
			res:    map[string]interface{}{"id": "", "name": "test", "color": "red", "shares": []interface{}{userID}, "plans": nil},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.body)
			req := httptest.NewRequest(http.MethodPost, "/calendars", bytes.NewBuffer(body))
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

func TestNewCalendarRouter_MakeCalendar(t *testing.T) {
	authRepo := testutils.NewAuthRepo()
	userID, sessionID := testutils.MakeSession(authRepo)
	calRepo := testutils.NewCalRepo()
	calRepo.CalUser().Create(context.Background(), cs.UserData{userID})

	l := testutils.NewLogger()
	authService := as.NewService(authRepo, l)
	calendarService := cs.NewService(calRepo, l)
	r := mux.NewRouter()
	NewCalendarRouter(r.PathPrefix("/calendars").Subrouter(), calendarService, authService)

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
			name:   "invalid contents",
			cookie: &cookie,
			body:   map[string]interface{}{"name": "", "color": ""},
			code:   http.StatusBadRequest,
			res:    map[string]interface{}{"message": "bad contents"},
		},
		{
			name:   "make calendar",
			cookie: &cookie,
			body:   map[string]interface{}{"name": "My plans", "color": "red"},
			code:   http.StatusOK,
			res:    map[string]interface{}{"id": "", "name": "My plans", "color": "red", "shares": []interface{}{userID}, "plans": nil},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.body)
			req := httptest.NewRequest(http.MethodPost, "/calendars", bytes.NewBuffer(body))
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

func TestNewCalendarRouter_RemoveCalendar(t *testing.T) {
	authRepo := testutils.NewAuthRepo()
	userID, sessionID := testutils.MakeSession(authRepo)
	otherID := uuid.New().String()
	calRepo := testutils.NewCalRepo()
	calendarID := uuid.New().String()
	otherCalID := uuid.New().String()
	calRepo.CalUser().Create(context.Background(), cs.UserData{userID})
	calRepo.Calendar().Create(context.Background(), cs.CalendarData{
		ID:     calendarID,
		Name:   "My plans",
		UserID: userID,
		Color:  "red",
		Shares: []string{userID},
	})
	calRepo.Calendar().Create(context.Background(), cs.CalendarData{
		ID:     otherCalID,
		Name:   "Work plans",
		UserID: otherID,
		Color:  "yellow",
		Shares: []string{otherID, userID},
	})

	l := testutils.NewLogger()
	authService := as.NewService(authRepo, l)
	calendarService := cs.NewService(calRepo, l)
	r := mux.NewRouter()
	NewCalendarRouter(r.PathPrefix("/calendars").Subrouter(), calendarService, authService)

	cookie := http.Cookie{
		Name:  "session_id",
		Value: sessionID,
	}

	testcases := []struct {
		name   string
		cookie *http.Cookie
		calID  string
		code   int
		res    map[string]interface{}
	}{
		{
			name:   "invalid calendar id",
			cookie: &cookie,
			calID:  uuid.New().String(),
			code:   http.StatusNotFound,
			res:    map[string]interface{}{"message": "not found"},
		},
		{
			name:   "remove my calendar",
			cookie: &cookie,
			calID:  calendarID,
			code:   http.StatusOK,
			res:    map[string]interface{}{"message": "remove calendar"},
		},
		{
			name:   "second remove my calendar",
			cookie: &cookie,
			calID:  calendarID,
			code:   http.StatusNotFound,
			res:    map[string]interface{}{"message": "not found"},
		},
		{
			name:   "remove other's calendar",
			cookie: &cookie,
			calID:  otherCalID,
			code:   http.StatusOK,
			res:    map[string]interface{}{"message": "remove calendar"},
		},
		{
			name:   "remove not shared calendar",
			cookie: &cookie,
			calID:  otherCalID,
			code:   http.StatusForbidden,
			res:    map[string]interface{}{"message": "unauthorization"},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/calendars/"+tc.calID, nil)
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

func TestNewCalendarRouter_ChangeCalendar(t *testing.T) {
	authRepo := testutils.NewAuthRepo()
	userID, sessionID := testutils.MakeSession(authRepo)
	otherID := uuid.New().String()
	calRepo := testutils.NewCalRepo()
	calRepo.CalUser().Create(context.Background(), cs.UserData{userID})
	calRepo.CalUser().Create(context.Background(), cs.UserData{otherID})
	cal := makeCalendar(calRepo, userID)
	cal2 := makeCalendar(calRepo, userID)
	sharedCal := makeCalendar(calRepo, otherID, userID)

	l := testutils.NewLogger()
	authService := as.NewService(authRepo, l)
	calendarService := cs.NewService(calRepo, l)
	r := mux.NewRouter()
	NewCalendarRouter(r.PathPrefix("/calendars").Subrouter(), calendarService, authService)

	cookie := http.Cookie{
		Name:  "session_id",
		Value: sessionID,
	}

	testcases := []struct {
		name   string
		cookie *http.Cookie
		calID  string
		body   map[string]interface{}
		code   int
		res    map[string]interface{}
	}{
		{
			name:   "invalid calendar id",
			cookie: &cookie,
			calID:  uuid.New().String(),
			body:   map[string]interface{}{"name": "Renamed", "color": "yellow", "shares": []interface{}{userID, otherID}},
			code:   http.StatusNotFound,
			res:    map[string]interface{}{"message": "not found"},
		},
		{
			name:   "invalid content",
			cookie: &cookie,
			calID:  cal.ID,
			body:   map[string]interface{}{"color": "yellow", "shares": []interface{}{userID, otherID}},
			code:   http.StatusBadRequest,
			res:    map[string]interface{}{"message": "bad contents"},
		},
		{
			name:   "invalid user id in shares",
			cookie: &cookie,
			calID:  cal.ID,
			body:   map[string]interface{}{"name": "Renamed", "color": "yellow", "shares": []interface{}{userID, uuid.New().String()}},
			code:   http.StatusBadRequest,
			res:    map[string]interface{}{"message": "bad contents"},
		},
		{
			name:   "not owner",
			cookie: &cookie,
			calID:  sharedCal.ID,
			body:   map[string]interface{}{"name": "Renamed", "color": "yellow", "shares": []interface{}{userID, otherID}},
			code:   http.StatusForbidden,
			res:    map[string]interface{}{"message": "unauthorization"},
		},
		{
			name:   "change calendar",
			cookie: &cookie,
			calID:  cal2.ID,
			body:   map[string]interface{}{"name": "Renamed", "color": "yellow", "shares": []interface{}{userID, otherID}},
			code:   http.StatusOK,
			res:    map[string]interface{}{"message": "change calendar"},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.body)
			req := httptest.NewRequest(http.MethodPatch, "/calendars/"+tc.calID, bytes.NewBuffer(body))
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
