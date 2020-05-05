package model

import (
	"fmt"

	cerror "github.com/x-color/calendar/model/error"
)

type Color string

const (
	RED    Color = "red"
	BLUE   Color = "blue"
	YELLOW Color = "yellow"
	GREEN  Color = "green"
)

func ConvertToColor(c string) (Color, error) {
	switch Color(c) {
	case RED:
		return RED, nil
	case BLUE:
		return BLUE, nil
	case YELLOW:
		return YELLOW, nil
	case GREEN:
		return GREEN, nil
	}
	return Color(""), cerror.NewInvalidContentError(
		nil,
		fmt.Sprintf("invalid color(%v)", c),
	)
}
