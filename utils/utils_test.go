package utils

import (
	"testing"
	"time"
)

func TestFormatTime(t *testing.T) {
	now := time.Now()

	cases := []struct {
		name string
		time string
		want string
	}{
		{"Test for same day", now.Format("2006-01-02T15:04:05"), "Сегодня в " + now.Format("15:04")},
		{"Test for previous day", now.AddDate(0, 0, -1).Format("2006-01-02T15:04:05"), "Вчера в " + now.AddDate(0, 0, -1).Format("15:04")},
		{"Test for two days ago", "2023-07-30T22:30:00", "30.07.2023 в 22:30"},
		{"Test for same year different month", "2023-07-01T12:15:00", "01.07.2023 в 12:15"},
		{"Test for same year beginning", "2023-01-01T00:00:00", "01.01.2023 в 00:00"},
		{"Test for end of previous year", "2022-12-31T23:59:59", "31.12.2022 в 23:59"},
		{"Test for previous year", "2022-05-15T08:30:00", "15.05.2022 в 08:30"},
		{"Test for beginning of century", "2000-01-01T00:00:00", "01.01.2000 в 00:00"},
		{"Test for leap year February", "2024-02-29T14:30:00", "29.02.2024 в 14:30"},
		{"Test for non-leap year February", "2023-02-28T14:30:00", "28.02.2023 в 14:30"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := FormatTime(tc.time)
			if err != nil {
				t.Fatalf("formatTime() error = %v", err)
			}
			if got != tc.want {
				t.Errorf("formatTime() = %v, want %v", got, tc.want)
			}
		})
	}
}
