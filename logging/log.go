package logging

import (
	"io"

	"log"
)

type logger struct {
	l *log.Logger
}

func (l *logger) Info(id, msg string) {
	l.l.Println("[INFO]", id, msg)
}

func (l *logger) Error(id, msg string) {
	l.l.Println("[ERROR]", id, msg)
}

func NewLogger(output io.Writer) logger {
	l := log.New(output, "", log.Ldate|log.Ltime)
	return logger{l}
}
