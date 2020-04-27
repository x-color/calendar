package logging

import (
	"io"

	"log"

	"github.com/x-color/calendar/service/auth"
)

type logger struct {
	l   *log.Logger
	uid string
}

func (l *logger) Uniq(id string) auth.Logger {
	return &logger{
		l:   l.l,
		uid: id,
	}
}

func (l *logger) Info(msg string) {
	l.l.Println("[INFO]", l.uid, msg)
}

func (l *logger) Error(msg string) {
	l.l.Println("[ERROR]", l.uid, msg)
}

func NewLogger(output io.Writer) logger {
	l := log.New(output, "", log.Ldate|log.Ltime)
	return logger{l, ""}
}
