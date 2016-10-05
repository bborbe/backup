package timeparser

import (
	"time"

	"github.com/bborbe/backup/constants"
)

type TimeParser interface {
	TimeByName(name string) (time.Time, error)
}

type timeParser struct{}

func New() *timeParser {
	return new(timeParser)
}

func (t *timeParser) TimeByName(name string) (time.Time, error) {
	return timeByName(name)
}

func timeByName(name string) (time.Time, error) {
	return time.Parse(constants.DATEFORMAT, name)
}
