package rest

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/x-color/calendar/logging"
	"github.com/x-color/calendar/repogitory/inmem"
	"github.com/x-color/calendar/service/auth"
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
				t.Errorf("want %v but %v", tc.code, rec.Code)
			}

			actual := strings.TrimRight(rec.Body.String(), "\n")
			expected := string(resBody)
			if actual != expected {
				t.Errorf("want %v but %v", expected, actual)
			}
		})
	}
}
