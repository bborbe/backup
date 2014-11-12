package timeparser

import "time"

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
	return time.Parse("2006-01-02T15:04:05", name)
}
