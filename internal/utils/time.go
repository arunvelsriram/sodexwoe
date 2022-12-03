package utils

import (
	"fmt"
	"strings"
	"time"
)

var months = []time.Month{time.January, time.February, time.March, time.April,
	time.May, time.June, time.July, time.August, time.September, time.October,
	time.November, time.December}

func GetMonthByName(name string) (time.Month, error) {
	if len(name) < 3 {
		return 0, fmt.Errorf("month name has less than 3 characters: %v", name)
	}

	for _, month := range months {
		if strings.EqualFold(name, month.String()) ||
			(len(name) == 3 && strings.EqualFold(name, month.String()[0:3])) {
			return month, nil
		}
	}

	return 0, fmt.Errorf("unable to understand month: %v", name)
}

func StartOfMonth(year int, month time.Month) time.Time {
	return time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
}

func EndOfMonth(year int, month time.Month) time.Time {
	temp := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	return temp.AddDate(0, 1, -temp.Day())
}
