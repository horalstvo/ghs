package util

import (
	"testing"
	"reflect"
	"time"
)

func TestWorkHoursSameTime(t *testing.T) {
	workHours := WorkHours(parse("2006-01-02T00:00:00Z"), parse("2006-01-02T00:00:00Z"))
	assertEqual(t, workHours, 0)
}

func TestWorkHoursSameDay(t *testing.T) {
	workHours := WorkHours(parse("2018-10-02T12:00:00Z"), parse("2018-10-02T16:00:00Z"))
	assertEqual(t, workHours, 4)

	workHours = WorkHours(parse("2018-10-02T2:00:00Z"), parse("2018-10-02T16:00:00Z"))
	assertEqual(t, workHours, 7)

	workHours = WorkHours(parse("2018-10-02T12:00:00Z"), parse("2018-10-02T20:00:00Z"))
	assertEqual(t, workHours, 6)
}

func TestWorkHoursDays(t *testing.T) {
	workHours := WorkHours(parse("2018-10-02T00:00:00Z"), parse("2018-10-03T00:00:00Z"))
	assertEqual(t, workHours, 9)

	// Tue -> Mon -> 4 working days.
	workHours = WorkHours(parse("2018-10-02T00:00:00Z"), parse("2018-10-08T00:00:00Z"))
	assertEqual(t, workHours, 4*9)

	workHours = WorkHours(parse("2018-10-04T23:00:00Z"), parse("2018-10-05T10:00:00Z"))
	assertEqual(t, workHours, 1)

	workHours = WorkHours(parse("2018-10-04T16:00:00Z"), parse("2018-10-05T07:00:00Z"))
	assertEqual(t, workHours, 2)
}

func TestSpecialCases(t *testing.T) {
	workHours := WorkHours(parse("2018-10-24T17:00:00Z"), parse("2018-10-25T13:00:00Z"))
	assertEqual(t, workHours, 5)
}

func parse(val string) time.Time {
	tim, err := time.Parse(time.RFC3339, val)
	Check(err)
	return tim
}

func assertEqual(t *testing.T, actual interface{}, expected interface{}) {
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("%s != %s", actual, expected)
	}
}
