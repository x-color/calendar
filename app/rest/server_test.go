package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/x-color/calendar/logging"
	mauth "github.com/x-color/calendar/model/auth"
	"github.com/x-color/calendar/repogitory/inmem"
	"github.com/x-color/calendar/service/auth"
	"golang.org/x/crypto/bcrypt"
)

func resetRepo() auth.Repogitory {
	r := inmem.NewRepogitory()
	return &r
}

func resetService(r auth.Repogitory) (auth.Service, auth.Logger) {
	l := logging.NewLogger(ioutil.Discard)
	return auth.NewService(r, &l), &l
}

func TestNewRouter_Signup(t *testing.T) {
	r := newRouter(resetService(resetRepo()))

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
			resBody, _ := json.Marshal(tc.res)

			req := httptest.NewRequest(http.MethodPost, "/auth/signup", bytes.NewBuffer(body))
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			if rec.Code != tc.code {
				t.Errorf("status code: want %v but %v", tc.code, rec.Code)
			}

			actual := strings.TrimRight(rec.Body.String(), "\n")
			expected := string(resBody)
			if actual != expected {
				t.Errorf("response body: want %v but %v", expected, actual)
			}
		})
	}
}

func TestNewRouter_Signin(t *testing.T) {
	repo := resetRepo()
	pwd, _ := bcrypt.GenerateFromPassword([]byte("P@ssw0rd"), bcrypt.DefaultCost)
	repo.User().Create(context.Background(), mauth.User{
		ID:       uuid.New().String(),
		Name:     "Alice",
		Password: string(pwd),
	})

	r := newRouter(resetService(repo))

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
			resBody, _ := json.Marshal(tc.res)

			req := httptest.NewRequest(http.MethodPost, "/auth/signin", bytes.NewBuffer(body))
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			if rec.Code != tc.code {
				t.Errorf("status code: want %v but %v", tc.code, rec.Code)
			}

			actual := strings.TrimRight(rec.Body.String(), "\n")
			expected := string(resBody)
			if actual != expected {
				t.Errorf("response body: want %v but %v", expected, actual)
			}

			if len(rec.Result().Cookies()) != tc.cookies {
				t.Errorf("cookies: want %v but %v", expected, actual)
			}
		})
	}
}
