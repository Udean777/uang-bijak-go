package models

import (
	"log"
	"time"
)

type DashboardSummary struct {
	TotalBalance int64 `json:"total_balance"`
	TotalIncome  int64 `json:"total_income"`
	TotalExpense int64 `json:"total_expense"`
}

type DashboardQuery struct {
	Month int `form:"month"`
	Year  int `form:"year"`
}

func (q *DashboardQuery) GetDateRange() (time.Time, time.Time) {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Println("Peringatan: Gagal load timezone, menggunakan UTC")
		loc = time.UTC
	}

	now := time.Now().In(loc)
	year := q.Year
	month := q.Month

	if year == 0 {
		year = now.Year()
	}
	if month == 0 {
		month = int(now.Month())
	}

	startTime := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, loc)
	endTime := startTime.AddDate(0, 1, 0).Add(-1 * time.Nanosecond)

	return startTime, endTime
}
