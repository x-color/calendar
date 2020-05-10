package calendar_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	. "github.com/x-color/calendar/app/rest/calendar"
	"github.com/x-color/calendar/app/rest/testutils"
	as "github.com/x-color/calendar/auth/service"
	cs "github.com/x-color/calendar/calendar/service"
)

func TestNewUserRouter_Authoraization(t *testing.T) {
	authRepo := testutils.NewAuthRepo()
	_, sessionID := testutils.MakeSession(authRepo)
	calRepo := testutils.NewCalRepo()

	l := testutils.NewLogger()
	authService := as.NewService(authRepo, l)
	calendarService := cs.NewService(calRepo, l)
	r := mux.NewRouter()
	NewUserRouter(r.PathPrefix("/register").Subrouter(), calendarService, authService)

	cookie := http.Cookie{
		Name:  "session_id",
		Value: sessionID,
	}

	testcases := []struct {
		name   string
		cookie *http.Cookie
		code   int
	}{
		{
			name:   "no cookie",
			cookie: nil,
			code:   http.StatusUnauthorized,
		},
		{
			name:   "invalid cookie",
			cookie: &http.Cookie{Name: "session_id", Value: uuid.New().String()},
			code:   http.StatusUnauthorized,
		},
		{
			name:   "valid cookie",
			cookie: &cookie,
			code:   http.StatusNoContent,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/register", nil)
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

func TestNewUserRouter_RegisterUser(t *testing.T) {
	authRepo := testutils.NewAuthRepo()
	_, sessionID := testutils.MakeSession(authRepo)
	calRepo := testutils.NewCalRepo()

	l := testutils.NewLogger()
	authService := as.NewService(authRepo, l)
	calendarService := cs.NewService(calRepo, l)
	r := mux.NewRouter()
	NewUserRouter(r.PathPrefix("/register").Subrouter(), calendarService, authService)

	cookie := http.Cookie{
		Name:  "session_id",
		Value: sessionID,
	}

	testcases := []struct {
		name   string
		cookie *http.Cookie
		code   int
	}{
		{
			name:   "register user",
			cookie: &cookie,
			code:   http.StatusNoContent,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/register", nil)
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
