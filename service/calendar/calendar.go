package calendar

type Repogitory interface {
}

type Logger interface {
	Uniq(id string) Logger
	Info(msg string)
	Error(msg string)
}

type Service struct {
	repo Repogitory
	log  Logger
}
