package utils

import (
	"fmt"
	"strings"
	"time"
)

func AnyLabelQ(labelNames ...string) string {
	labelQs := make([]string, 0, len(labelNames))
	for _, labelName := range labelNames {
		labelQs = append(labelQs, fmt.Sprintf("label:\"%s\"", labelName))
	}
	return fmt.Sprintf("(%s)", strings.Join(labelQs, " OR "))
}

func WithinMonthQ(year int, month time.Month) string {
	layout := "2006/01/02"
	afterDate := StartOfMonth(year, month)
	beforeDate := EndOfMonth(year, month)
	return fmt.Sprintf("after:%s before:%s", afterDate.Format(layout), beforeDate.Format(layout))
}
