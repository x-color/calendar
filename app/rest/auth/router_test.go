package auth_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	. "github.com/x-color/calendar/app/rest/auth"
	"github.com/x-color/calendar/auth/repogitory/inmem"
	as "github.com/x-color/calendar/auth/service"
	cs "github.com/x-color/calendar/calendar/service"
	"github.com/x-color/calendar/logging"
	"github.com/x-color/calendar/testutils"
	"golang.org/x/crypto/bcrypt"
)

func newAuthRepo() as.Repogitory {
	r := inmem.NewRepogitory()
	return &r
}

func newLogger() logging.Logger {
	l := logging.NewLogger(ioutil.Discard)
	return &l
}

func dummyCalService() cs.Service {
	return cs.Service{}
}

func makeSession(authRepo as.Repogitory) (string, string) {
	userID := uuid.New().String()
	sessionID := uuid.New().String()
	authRepo.Session().Create(context.Background(), as.SessionData{
		ID:      sessionID,
		UserID:  userID,
		Expires: time.Now().Add(time.Hour).Unix(),
	})
	return userID, sessionID
}

func TestNewRouter_Signup(t *testing.T) {
	repo := newAuthRepo()
	l := newLogger()
	authService := as.NewService(repo, l)
	r := mux.NewRouter()
	NewRouter(r.PathPrefix("/auth").Subrouter(), authService)

	testcases := []struct {
		name string
		body map[string]string
		code int
		res  map[string]string
	}{
		{
			name: "invalid password",
			body: map[string]string{"name": "Alice", "password": "Password"},
			code: http.StatusBadRequest,
			res:  map[string]string{"message": "bad contents"},
		},
		{
			name: "invalid name",
			body: map[string]string{"name": "", "password": "P@ssw0rd"},
			code: http.StatusBadRequest,
			res:  map[string]string{"message": "bad contents"},
		},
		{
			name: "signup new user",
			body: map[string]string{"name": "Alice", "password": "P@ssw0rd"},
			code: http.StatusOK,
			res:  map[string]string{"name": "Alice", "password": ""},
		},
		{
			name: "user already exist",
			body: map[string]string{"name": "Alice", "password": "P@ssw0rd"},
			code: http.StatusConflict,
			res:  map[string]string{"message": "user already exist"},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.body)

			req := httptest.NewRequest(http.MethodPost, "/auth/signup", bytes.NewBuffer(body))
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			if rec.Code != tc.code {
				t.Errorf("status code: want %v but %v", tc.code, rec.Code)
			}

			var actual map[string]string
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

func TestNewRouter_Signin(t *testing.T) {
	repo := newAuthRepo()
	pwd, _ := bcrypt.GenerateFromPassword([]byte("P@ssw0rd"), bcrypt.DefaultCost)
	repo.User().Create(context.Background(), as.UserData{
		ID:       uuid.New().String(),
		Name:     "Alice",
		Password: string(pwd),
	})

	l := newLogger()
	authService := as.NewService(repo, l)
	r := mux.NewRouter()
	NewRouter(r.PathPrefix("/auth").Subrouter(), authService)

	testcases := []struct {
		name    string
		body    map[string]string
		code    int
		res     map[string]string
		cookies int
	}{
		{
			name: "invalid password",
			body: map[string]string{"name": "Alice", "password": "p@SSW0RD"},
			code: http.StatusUnauthorized,
			res:  map[string]string{"message": "signin failed"},
		},
		{
			name: "user does not exist",
			body: map[string]string{"name": "Bob", "password": "P@ssw0rd"},
			code: http.StatusUnauthorized,
			res:  map[string]string{"message": "signin failed"},
		},
		{
			name:    "signin",
			body:    map[string]string{"name": "Alice", "password": "P@ssw0rd"},
			code:    http.StatusOK,
			res:     map[string]string{"message": "signin"},
			cookies: 1,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.body)

			req := httptest.NewRequest(http.MethodPost, "/auth/signin", bytes.NewBuffer(body))
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			if rec.Code != tc.code {
				t.Errorf("status code: want %v but %v", tc.code, rec.Code)
			}

			var actual map[string]string
			if err := json.Unmarshal(rec.Body.Bytes(), &actual); err != nil {
				t.Errorf("invalid response body: %v", rec.Body.String())
			}
			expected := tc.res

			if d := cmp.Diff(expected, actual, testutils.IgnoreKey("id")); d != "" {
				t.Errorf("invalid response body: \n%v", d)
			}

			if len(rec.Result().Cookies()) != tc.cookies {
				t.Errorf("cookies: want %v but %v", expected, actual)
			}
		})
	}
}

func TestNewRouter_Signout(t *testing.T) {
	repo := newAuthRepo()
	_, sessionID := makeSession(repo)

	l := newLogger()
	authService := as.NewService(repo, l)
	r := mux.NewRouter()
	NewRouter(r.PathPrefix("/auth").Subrouter(), authService)

	testcases := []struct {
		name   string
		cookie *http.Cookie
		code   int
		res    map[string]string
	}{
		{
			name:   "no cookie",
			cookie: nil,
			code:   http.StatusUnauthorized,
			res:    map[string]string{"message": "signout failed"},
		},
		{
			name: "invalid cookie",
			cookie: &http.Cookie{
				Name:  "session_id",
				Value: uuid.New().String(),
			},
			code: http.StatusUnauthorized,
			res:  map[string]string{"message": "signout failed"},
		},
		{
			name: "signout",
			cookie: &http.Cookie{
				Name:  "session_id",
				Value: sessionID,
			},
			code: http.StatusOK,
			res:  map[string]string{"message": "signout"},
		},
		{
			name: "second signout",
			cookie: &http.Cookie{
				Name:  "session_id",
				Value: sessionID,
			},
			code: http.StatusUnauthorized,
			res:  map[string]string{"message": "signout failed"},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/auth/signout", nil)
			if tc.cookie != nil {
				req.AddCookie(tc.cookie)
			}
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			if rec.Code != tc.code {
				t.Errorf("status code: want %v but %v", tc.code, rec.Code)
			}

			var actual map[string]string
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
