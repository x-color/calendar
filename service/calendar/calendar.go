package calendar

import "github.com/x-color/calendar/service"

type Repogitory interface {
}

type Service struct {
	repo Repogitory
	log  service.Logger
}

func NewService(repo Repogitory, log service.Logger) Service {
	return Service{
		repo: repo,
		log:  log,
	}
}
