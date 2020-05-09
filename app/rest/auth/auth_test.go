package auth_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	. "github.com/x-color/calendar/app/rest/auth"
	"github.com/x-color/calendar/app/rest/testutils"
	as "github.com/x-color/calendar/auth/service"
	"golang.org/x/crypto/bcrypt"
)

func TestNewRouter_Signup(t *testing.T) {
	repo := testutils.NewAuthRepo()
	l := testutils.NewLogger()
	authService := as.NewService(repo, l)
	r := mux.NewRouter()
	NewRouter(r.PathPrefix("/auth").Subrouter(), authService)

	testcases := []struct {
		name string
		body map[string]string
		code int
	}{
		{
			name: "invalid password",
			body: map[string]string{"name": "Alice", "password": "Password"},
			code: http.StatusBadRequest,
		},
		{
			name: "invalid name",
			body: map[string]string{"name": "", "password": "P@ssw0rd"},
			code: http.StatusBadRequest,
		},
		{
			name: "signup new user",
			body: map[string]string{"name": "Alice", "password": "P@ssw0rd"},
			code: http.StatusNoContent,
		},
		{
			name: "user already exist",
			body: map[string]string{"name": "Alice", "password": "P@ssw0rd"},
			code: http.StatusConflict,
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
		})
	}
}

func TestNewRouter_Signin(t *testing.T) {
	repo := testutils.NewAuthRepo()
	pwd, _ := bcrypt.GenerateFromPassword([]byte("P@ssw0rd"), bcrypt.DefaultCost)
	repo.User().Create(context.Background(), as.UserData{
		ID:       uuid.New().String(),
		Name:     "Alice",
		Password: string(pwd),
	})

	l := testutils.NewLogger()
	authService := as.NewService(repo, l)
	r := mux.NewRouter()
	NewRouter(r.PathPrefix("/auth").Subrouter(), authService)

	testcases := []struct {
		name    string
		body    map[string]string
		code    int
		cookies int
	}{
		{
			name: "invalid password",
			body: map[string]string{"name": "Alice", "password": "p@SSW0RD"},
			code: http.StatusUnauthorized,
		},
		{
			name: "user does not exist",
			body: map[string]string{"name": "Bob", "password": "P@ssw0rd"},
			code: http.StatusUnauthorized,
		},
		{
			name:    "signin",
			body:    map[string]string{"name": "Alice", "password": "P@ssw0rd"},
			code:    http.StatusNoContent,
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

			if len(rec.Result().Cookies()) != tc.cookies {
				t.Errorf("cookies: want %v but %v", tc.cookies, rec.Result().Cookies())
			}
		})
	}
}

func TestNewRouter_Signout(t *testing.T) {
	repo := testutils.NewAuthRepo()
	_, sessionID := testutils.MakeSession(repo)

	l := testutils.NewLogger()
	authService := as.NewService(repo, l)
	r := mux.NewRouter()
	NewRouter(r.PathPrefix("/auth").Subrouter(), authService)

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
			name: "invalid cookie",
			cookie: &http.Cookie{
				Name:  "session_id",
				Value: uuid.New().String(),
			},
			code: http.StatusUnauthorized,
		},
		{
			name: "signout",
			cookie: &http.Cookie{
				Name:  "session_id",
				Value: sessionID,
			},
			code: http.StatusNoContent,
		},
		{
			name: "second signout",
			cookie: &http.Cookie{
				Name:  "session_id",
				Value: sessionID,
			},
			code: http.StatusUnauthorized,
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
		})
	}
}
