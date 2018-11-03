package util

import (
	"time"
)

const workdayStart = 9
const workdayEnd = 18
const dayLength = workdayEnd - workdayStart

func WorkHours(start time.Time, end time.Time) int {
	if start.After(end) {
		return 0 // Or return error.
	}

	// Normalize start and end.
	if start.Hour() > workdayEnd {
		start = time.Date(start.Year(), start.Month(), start.Day() + 1, workdayStart, 0, 0, 0, start.Location())
	} else if start.Hour() < workdayStart {
		start = time.Date(start.Year(), start.Month(), start.Day(), workdayStart, 0, 0, 0, start.Location())
	}

	if end.Hour() < workdayStart {
		end = time.Date(end.Year(), end.Month(), end.Day() - 1, workdayEnd, 0, 0, 0, end.Location())
	} else if end.Hour() > workdayEnd {
		end = time.Date(end.Year(), end.Month(), end.Day(), workdayEnd, 0, 0, 0, end.Location())
	}

	if end.Sub(start).Hours() < 24 && start.Weekday() == end.Weekday() {
		return end.Hour() - start.Hour()
	}

	hours := workdayEnd - start.Hour()
	day := start.AddDate(0, 0, 1)
	for day.Before(end.AddDate(0, 0, -1)) {
		if day.Weekday() < 6 {
			hours += dayLength
		}
		day = day.AddDate(0, 0, 1)
	}

	hours += end.Hour() - workdayStart

	return hours
}
