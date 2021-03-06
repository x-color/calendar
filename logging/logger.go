package logging

type Logger interface {
	Uniq(id string) Logger
	Info(msg string)
	Error(msg string)
}
