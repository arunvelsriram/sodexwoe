package models

import (
	"time"
)

type BillEmails []BillEmail

type BillEmail struct {
	BillName string
	Year     int
	Month    time.Month
	Bill     Bill
}

type Bill struct {
	Filename string
	Data     []byte
}
