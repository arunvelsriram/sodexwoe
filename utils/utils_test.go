package utils_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/arunvelsriram/sodexwoe/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetMonth(t *testing.T) {
	params := []struct {
		monthName     string
		expectedMonth time.Month
		expectedErr   error
	}{
		{"March", time.March, nil},
		{"JaNuAry", time.January, nil},
		{"Abcd", 0, fmt.Errorf("unable to understand month: Abcd")},
		{"Ja", 0, fmt.Errorf("month name has less than 3 characters: Ja")},
		{"janxxx", 0, fmt.Errorf("unable to understand month: janxxx")},
	}

	for _, param := range params {
		t.Run(fmt.Sprintf("MonthName=%s", param.monthName), func(t *testing.T) {
			months, err := utils.GetMonth(param.monthName)

			assert.Equal(t, param.expectedErr, err)
			assert.Equal(t, param.expectedMonth, months)
		})
	}
}

func TestStartOfMonth(t *testing.T) {
	params := []struct {
		year     int
		month    time.Month
		expected time.Time
	}{
		{2022, time.March, time.Date(2022, time.March, 1, 0, 0, 0, 0, time.Local)},
	}

	for _, param := range params {
		t.Run(fmt.Sprintf("Year=%d Month=%s", param.year, param.month), func(t *testing.T) {
			actual := utils.StartOfMonth(param.year, param.month)

			assert.Equal(t, param.expected, actual)
		})
	}
}

func TestEndOfMonth(t *testing.T) {
	params := []struct {
		year     int
		month    time.Month
		expected time.Time
	}{
		{2022, time.April, time.Date(2022, time.April, 30, 0, 0, 0, 0, time.Local)},
		{2022, time.February, time.Date(2022, time.February, 28, 0, 0, 0, 0, time.Local)},
		{2020, time.February, time.Date(2020, time.February, 29, 0, 0, 0, 0, time.Local)},
		{2022, time.December, time.Date(2022, time.December, 31, 0, 0, 0, 0, time.Local)},
	}

	for _, param := range params {
		t.Run(fmt.Sprintf("Year=%d Month=%s", param.year, param.month), func(t *testing.T) {
			actual := utils.EndOfMonth(param.year, param.month)

			assert.Equal(t, param.expected, actual)
		})
	}
}
