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

	// Normalize start.
	if start.Hour() > workdayEnd {
		start = time.Date(start.Year(), start.Month(), start.Day() + 1, workdayStart, 0, 0, 0, start.Location())
	} else if start.Hour() < workdayStart {
		start = time.Date(start.Year(), start.Month(), start.Day(), workdayStart, 0, 0, 0, start.Location())
	}

	for weekend(start) {
		start = start.AddDate(0, 0, 1)
		start = time.Date(start.Year(), start.Month(), start.Day(), workdayStart, 0, 0, 0, start.Location())
	}

	// Normalize end.
	if end.Hour() < workdayStart {
		end = time.Date(end.Year(), end.Month(), end.Day() - 1, workdayEnd, 0, 0, 0, end.Location())
	} else if end.Hour() > workdayEnd {
		end = time.Date(end.Year(), end.Month(), end.Day(), workdayEnd, 0, 0, 0, end.Location())
	}

	for weekend(end) {
		end = end.AddDate(0, 0, -1)
		end = time.Date(end.Year(), end.Month(), end.Day(), workdayEnd, 0, 0, 0, end.Location())
	}

	if end.Sub(start).Hours() < 24 && start.Weekday() == end.Weekday() {
		return end.Hour() - start.Hour()
	}

	if start.After(end) {
		return 0
	}

	hours := workdayEnd - start.Hour()

	day := start.AddDate(0, 0, 1)
	for day.Before(end.AddDate(0, 0, -1)) {
		// Sunday, Saturday.
		if !weekend(day) {
			hours += dayLength
		}
		day = day.AddDate(0, 0, 1)
	}

	hours += end.Hour() - workdayStart

	return hours
}

func weekend(day time.Time) bool {
	return day.Weekday() == 0 || day.Weekday() == 6
}
