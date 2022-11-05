package utils

import (
	"fmt"
	"strings"
	"time"
)

var months = []time.Month{
	time.January,
	time.February,
	time.March,
	time.April,
	time.May,
	time.June,
	time.July,
	time.August,
	time.September,
	time.October,
	time.November,
	time.December,
}

func GetMonth(monthName string) (time.Month, error) {
	if len(monthName) < 3 {
		return 0, fmt.Errorf("month name has less than 3 characters: %v", monthName)
	}

	for _, it := range months {
		if strings.EqualFold(monthName, it.String()) ||
			(len(monthName) == 3 && strings.EqualFold(monthName, it.String()[0:3])) {
			return it, nil
		}
	}

	return 0, fmt.Errorf("unable to understand month: %v", monthName)
}

func StartOfMonth(year int, month time.Month) time.Time {
	return time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
}

func EndOfMonth(year int, month time.Month) time.Time {
	temp := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	return temp.AddDate(0, 1, -temp.Day())
}
