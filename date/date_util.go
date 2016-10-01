package date

import "time"

func DayEqual(t time.Time, now time.Time) bool {
	return t.Year() == now.Year() && t.Month() == now.Month() && t.Day() == now.Day()
}
