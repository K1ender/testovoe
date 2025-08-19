package utils

import (
	"fmt"
	"time"
)

func MonthsOverlap(aStart, aEnd, bStart, bEnd time.Time) int {
	if aEnd.Before(aStart) || bEnd.Before(bStart) {
		return 0
	}
	s := aStart
	if bStart.After(s) {
		s = bStart
	}
	e := aEnd
	if bEnd.Before(e) {
		e = bEnd
	}
	if e.Before(s) {
		return 0
	}
	months := int((e.Year()-s.Year())*12+int(e.Month())-int(s.Month())) + 1
	if months < 0 {
		return 0
	}
	return months
}

func ParseMonthYear(s string) (time.Time, error) {
	var m, y int
	if _, err := fmt.Sscanf(s, "%02d-%04d", &m, &y); err != nil {
		return time.Time{}, err
	}
	return time.Date(y, time.Month(m), 1, 0, 0, 0, 0, time.UTC), nil
}
