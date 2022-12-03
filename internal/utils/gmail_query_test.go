package utils_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/arunvelsriram/sodexwoe/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestAnyLabelQ(t *testing.T) {
	params := []struct {
		labelNames    []string
		expectedQuery string
	}{
		{[]string{"airtel"}, "(label:\"airtel\")"},
		{[]string{"Jio", "Postpaid Bill / Airtel"}, "(label:\"Jio\" OR label:\"Postpaid Bill / Airtel\")"},
	}

	for _, param := range params {
		t.Run(fmt.Sprintf("LabelNames=%s", strings.Join(param.labelNames, ",")), func(t *testing.T) {
			actualQuery := utils.AnyLabelQ(param.labelNames...)

			assert.Equal(t, param.expectedQuery, actualQuery)
		})
	}
}

func TestWithinMonthQ(t *testing.T) {
	params := []struct {
		year          int
		month         time.Month
		expectedQuery string
	}{
		{2022, time.April, "after:2022/04/01 before:2022/04/30"},
		{2022, time.August, "after:2022/08/01 before:2022/08/31"},
		{2022, time.February, "after:2022/02/01 before:2022/02/28"},
		{2020, time.February, "after:2020/02/01 before:2020/02/29"},
	}

	for _, param := range params {
		t.Run(fmt.Sprintf("Year=%d Month=%s", param.year, param.month.String()), func(t *testing.T) {
			actualQuery := utils.WithinMonthQ(param.year, param.month)

			assert.Equal(t, param.expectedQuery, actualQuery)
		})
	}
}
