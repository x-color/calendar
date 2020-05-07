package calendar_test

import (
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
		res    map[string]interface{}
	}{
		{
			name:   "no cookie",
			cookie: nil,
			code:   http.StatusUnauthorized,
			res:    map[string]interface{}{"message": "unauthorization"},
		},
		{
			name:   "invalid cookie",
			cookie: &http.Cookie{Name: "session_id", Value: uuid.New().String()},
			code:   http.StatusUnauthorized,
			res:    map[string]interface{}{"message": "unauthorization"},
		},
		{
			name:   "valid cookie",
			cookie: &cookie,
			code:   http.StatusOK,
			res:    map[string]interface{}{"message": "register"},
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
		res    map[string]interface{}
	}{
		{
			name:   "register user",
			cookie: &cookie,
			code:   http.StatusOK,
			res:    map[string]interface{}{"message": "register"},
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
