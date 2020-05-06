package service

import (
	"github.com/x-color/calendar/logging"
)

type Service struct {
	repo Repogitory
	log  logging.Logger
}

func NewService(repo Repogitory, log logging.Logger) Service {
	return Service{
		repo: repo,
		log:  log,
	}
}
