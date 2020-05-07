package testutils

import (
	"context"
	"io/ioutil"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	ar "github.com/x-color/calendar/auth/repogitory/inmem"
	as "github.com/x-color/calendar/auth/service"
	cr "github.com/x-color/calendar/calendar/repogitory/inmem"
	cs "github.com/x-color/calendar/calendar/service"
	"github.com/x-color/calendar/logging"
)

func IgnoreKey(key string) cmp.Option {
	return cmpopts.IgnoreMapEntries(func(k string, t interface{}) bool {
		return k == key
	})
}

func NewAuthRepo() as.Repogitory {
	r := ar.NewRepogitory()
	return &r
}

func NewCalRepo() cs.Repogitory {
	r := cr.NewRepogitory()
	return &r
}

func NewLogger() logging.Logger {
	l := logging.NewLogger(ioutil.Discard)
	return &l
}

func DummyCalService() cs.Service {
	return cs.Service{}
}

func MakeSession(authRepo as.Repogitory) (string, string) {
	userID := uuid.New().String()
	sessionID := uuid.New().String()
	authRepo.Session().Create(context.Background(), as.SessionData{
		ID:      sessionID,
		UserID:  userID,
		Expires: time.Now().Add(time.Hour).Unix(),
	})
	return userID, sessionID
}
