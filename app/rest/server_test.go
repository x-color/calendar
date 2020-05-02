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
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/x-color/calendar/logging"
	mauth "github.com/x-color/calendar/model/auth"
	"github.com/x-color/calendar/repogitory/inmem"
	"github.com/x-color/calendar/service"
	"github.com/x-color/calendar/service/auth"
	"github.com/x-color/calendar/service/calendar"
	"golang.org/x/crypto/bcrypt"
)

func newAuthRepo() auth.Repogitory {
	r := inmem.NewRepogitory()
	return &r
}

func newCalRepo() calendar.Repogitory {
	r := inmem.NewRepogitory()
	return &r
}

func newLogger() service.Logger {
	l := logging.NewLogger(ioutil.Discard)
	return &l
}

func dummyCalService() calendar.Service {
	return calendar.Service{}
}

func ignoreKey(key string) cmp.Option {
	return cmpopts.IgnoreMapEntries(func(k string, t interface{}) bool {
		return k == key
	})
}

func TestNewRouter_Signup(t *testing.T) {
	repo := newAuthRepo()
	l := newLogger()
	as := auth.NewService(repo, l)
	r := newRouter(as, dummyCalService(), l)

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

			actual := strings.TrimSpace(rec.Body.String())
			expected := string(resBody)
			if actual != expected {
				t.Errorf("response body: want %v but %v", expected, actual)
			}
		})
	}
}

func TestNewRouter_Signin(t *testing.T) {
	repo := newAuthRepo()
	pwd, _ := bcrypt.GenerateFromPassword([]byte("P@ssw0rd"), bcrypt.DefaultCost)
	repo.User().Create(context.Background(), mauth.User{
		ID:       uuid.New().String(),
		Name:     "Alice",
		Password: string(pwd),
	})

	l := newLogger()
	as := auth.NewService(repo, l)
	r := newRouter(as, dummyCalService(), l)

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

			actual := strings.TrimSpace(rec.Body.String())
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

func TestNewRouter_Signout(t *testing.T) {
	repo := newAuthRepo()
	sessionID := uuid.New().String()
	repo.Session().Create(context.Background(), mauth.Session{
		ID:      sessionID,
		UserID:  uuid.New().String(),
		Expires: time.Now().Add(time.Hour),
	})

	l := newLogger()
	as := auth.NewService(repo, l)
	r := newRouter(as, dummyCalService(), l)

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

			actual := strings.TrimSpace(rec.Body.String())
			resBody, _ := json.Marshal(tc.res)
			expected := string(resBody)
			if actual != expected {
				t.Errorf("response body: want %v but %v", expected, actual)
			}
		})
	}
}

func TestNewRouter_MakeCalendar(t *testing.T) {
	authRepo := newAuthRepo()
	userID := uuid.New().String()
	sessionID := uuid.New().String()
	authRepo.Session().Create(context.Background(), mauth.Session{
		ID:      sessionID,
		UserID:  userID,
		Expires: time.Now().Add(time.Hour),
	})
	calRepo := newCalRepo()

	l := newLogger()
	as := auth.NewService(authRepo, l)
	cs := calendar.NewService(calRepo, l)
	r := newRouter(as, cs, l)

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
			res:    map[string]interface{}{"id": "", "name": "My plans", "color": "red", "private": true, "shares": []interface{}{userID}, "plans": []interface{}{}},
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

			if d := cmp.Diff(expected, actual, ignoreKey("id")); d != "" {
				t.Errorf("invalid response body: \n%v", d)
			}
		})
	}
}
