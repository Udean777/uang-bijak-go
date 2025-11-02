package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDashboardQuery_GetDateRange(t *testing.T) {
	// Dapatkan lokasi untuk konsistensi
	loc, _ := time.LoadLocation("Asia/Jakarta")

	t.Run("Default (Bulan Ini)", func(t *testing.T) {
		q := DashboardQuery{} // Query kosong

		// Kita perlu tahu bulan ini untuk validasi
		now := time.Now().In(loc)
		expectedStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)
		expectedEnd := expectedStart.AddDate(0, 1, 0).Add(-1 * time.Nanosecond)

		start, end := q.GetDateRange()

		assert.Equal(t, expectedStart, start)
		assert.Equal(t, expectedEnd, end)
	})

	t.Run("Bulan & Tahun Spesifik (Oktober 2025)", func(t *testing.T) {
		q := DashboardQuery{
			Month: 10,
			Year:  2025,
		}

		// 1 Oktober 2025, 00:00:00
		expectedStart := time.Date(2025, time.October, 1, 0, 0, 0, 0, loc)
		// 31 Oktober 2025, 23:59:59...
		expectedEnd := time.Date(2025, time.October, 31, 23, 59, 59, 999999999, loc)

		start, end := q.GetDateRange()

		assert.Equal(t, expectedStart, start)
		assert.Equal(t, expectedEnd, end)
	})
}
