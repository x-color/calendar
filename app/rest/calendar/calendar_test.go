package calendar_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
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
	calRepo.User().Create(context.Background(), cs.UserData{userID})

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
		res    CalendarContent
	}{
		{
			name:   "no cookie",
			cookie: nil,
			body:   map[string]interface{}{"name": "My plans", "color": "red"},
			code:   http.StatusUnauthorized,
		},
		{
			name:   "invalid cookie",
			cookie: &http.Cookie{Name: "session_id", Value: uuid.New().String()},
			body:   map[string]interface{}{"name": "My plans", "color": "red"},
			code:   http.StatusUnauthorized,
		},
		{
			name:   "valid cookie",
			cookie: &cookie,
			body:   map[string]interface{}{"name": "My plans", "color": "red"},
			code:   http.StatusOK,
			res: CalendarContent{
				ID:     "",
				Name:   "My plans",
				Color:  "red",
				Shares: []string{userID},
				Plans:  nil,
			},
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

			var actual CalendarContent
			if len(rec.Body.Bytes()) > 0 {
				if err := json.Unmarshal(rec.Body.Bytes(), &actual); err != nil {
					t.Errorf("invalid response body: %v", rec.Body.String())
				}
			}
			expected := tc.res

			if d := cmp.Diff(expected, actual, cmpopts.IgnoreFields(CalendarContent{}, "ID")); d != "" {
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
	calRepo.User().Create(context.Background(), cs.UserData{userID})

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
		res    CalendarContent
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
			body:   map[string]interface{}{"name": "test", "color": "red"},
			code:   http.StatusOK,
			res: CalendarContent{
				ID:     "",
				Name:   "test",
				Color:  "red",
				Shares: []string{userID},
				Plans:  nil,
			},
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

			var actual CalendarContent
			if len(rec.Body.Bytes()) > 0 {
				if err := json.Unmarshal(rec.Body.Bytes(), &actual); err != nil {
					t.Errorf("invalid response body: %v", rec.Body.String())
				}
			}
			expected := tc.res

			if d := cmp.Diff(expected, actual, cmpopts.IgnoreFields(CalendarContent{}, "ID")); d != "" {
				t.Errorf("invalid response body: \n%v", d)
			}
		})
	}
}

func TestNewCalendarRouter_MakeCalendar(t *testing.T) {
	authRepo := testutils.NewAuthRepo()
	userID, sessionID := testutils.MakeSession(authRepo)
	calRepo := testutils.NewCalRepo()
	calRepo.User().Create(context.Background(), cs.UserData{userID})

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
		res    CalendarContent
	}{
		{
			name:   "invalid contents",
			cookie: &cookie,
			body:   map[string]interface{}{"name": "", "color": ""},
			code:   http.StatusBadRequest,
		},
		{
			name:   "make calendar",
			cookie: &cookie,
			body:   map[string]interface{}{"name": "My plans", "color": "red"},
			code:   http.StatusOK,
			res: CalendarContent{
				ID:     "",
				Name:   "My plans",
				Color:  "red",
				Shares: []string{userID},
				Plans:  nil,
			},
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

			var actual CalendarContent
			if len(rec.Body.Bytes()) > 0 {
				if err := json.Unmarshal(rec.Body.Bytes(), &actual); err != nil {
					t.Errorf("invalid response body: %v", rec.Body.String())
				}
			}
			expected := tc.res

			if d := cmp.Diff(expected, actual, cmpopts.IgnoreFields(CalendarContent{}, "ID")); d != "" {
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
	calRepo.User().Create(context.Background(), cs.UserData{userID})
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
	}{
		{
			name:   "invalid calendar id",
			cookie: &cookie,
			calID:  uuid.New().String(),
			code:   http.StatusNotFound,
		},
		{
			name:   "remove my calendar",
			cookie: &cookie,
			calID:  calendarID,
			code:   http.StatusNoContent,
		},
		{
			name:   "second remove my calendar",
			cookie: &cookie,
			calID:  calendarID,
			code:   http.StatusNotFound,
		},
		{
			name:   "remove other's calendar",
			cookie: &cookie,
			calID:  otherCalID,
			code:   http.StatusNoContent,
		},
		{
			name:   "remove not shared calendar",
			cookie: &cookie,
			calID:  otherCalID,
			code:   http.StatusForbidden,
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
		})
	}
}

func TestNewCalendarRouter_ChangeCalendar(t *testing.T) {
	authRepo := testutils.NewAuthRepo()
	userID, sessionID := testutils.MakeSession(authRepo)
	otherID := uuid.New().String()
	calRepo := testutils.NewCalRepo()
	calRepo.User().Create(context.Background(), cs.UserData{userID})
	calRepo.User().Create(context.Background(), cs.UserData{otherID})
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
	}{
		{
			name:   "invalid calendar id",
			cookie: &cookie,
			calID:  uuid.New().String(),
			body:   map[string]interface{}{"name": "Renamed", "color": "yellow", "shares": []interface{}{userID, otherID}},
			code:   http.StatusNotFound,
		},
		{
			name:   "invalid content",
			cookie: &cookie,
			calID:  cal.ID,
			body:   map[string]interface{}{"color": "yellow", "shares": []interface{}{userID, otherID}},
			code:   http.StatusBadRequest,
		},
		{
			name:   "invalid user id in shares",
			cookie: &cookie,
			calID:  cal.ID,
			body:   map[string]interface{}{"name": "Renamed", "color": "yellow", "shares": []interface{}{userID, uuid.New().String()}},
			code:   http.StatusBadRequest,
		},
		{
			name:   "not owner",
			cookie: &cookie,
			calID:  sharedCal.ID,
			body:   map[string]interface{}{"name": "Renamed", "color": "yellow", "shares": []interface{}{userID, otherID}},
			code:   http.StatusForbidden,
		},
		{
			name:   "change calendar",
			cookie: &cookie,
			calID:  cal2.ID,
			body:   map[string]interface{}{"name": "Renamed", "color": "yellow", "shares": []interface{}{userID, otherID}},
			code:   http.StatusNoContent,
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
		})
	}
}
