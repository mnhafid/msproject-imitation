package main

import (
	"time"
)

type WorkDay struct {
	Day     time.Weekday
	Working bool
}

type ProjectCalendar struct {
	ID         string
	Name       string
	WorkDays   []WorkDay
	Exceptions []ProjectCalendarExepction
}

type ProjectCalendarExepction struct {
	StartDate time.Time
	EndDate   time.Time
	Working   bool
}

func (pc ProjectCalendar) IsWorkingDay(date time.Time) bool {
	// Check if the date falls within any exceptions
	for _, exception := range pc.Exceptions {
		if (date.After(exception.StartDate) || date.Equal(exception.StartDate)) &&
			(date.Before(exception.EndDate) || date.Equal(exception.EndDate)) {
			return exception.Working
		}
	}
	for _, workDay := range pc.WorkDays {
		if workDay.Day == date.Weekday() && workDay.Working {
			return true
		}
	}
	return false
}

func AddDate(startDate time.Time, duration int, pc ProjectCalendar) time.Time {
	daysAdded := 0
	for daysAdded < duration {
		// Check if the current date is a working day
		if pc.IsWorkingDay(startDate) {
			daysAdded++
		}
		// Move to the next day only if the duration hasn't been met
		if daysAdded < duration {
			startDate = startDate.AddDate(0, 0, 1)
		}
	}
	return startDate
}
