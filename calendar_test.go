package main

import (
	"testing"
	"time"
)

func TestCalendar(t *testing.T) {
	t.Run("Test Add Date", func(t *testing.T) {
		// AddDate(time.Date(2024, 6, 23, 0, 0, 0, 0, time.UTC), 10, ProjectCalendar{})
		ProjectCalendar := ProjectCalendar{
			ID:   "1",
			Name: "Test",
			WorkDays: []WorkDay{
				{Day: time.Monday, Working: true},
				{Day: time.Tuesday, Working: true},
				{Day: time.Wednesday, Working: true},
				{Day: time.Thursday, Working: true},
				{Day: time.Friday, Working: true},
				{Day: time.Saturday, Working: false},
				{Day: time.Sunday, Working: false}}}
		expected := time.Date(2024, 7, 5, 0, 0, 0, 0, time.UTC)
		result := AddDate(time.Date(2024, 6, 23, 0, 0, 0, 0, time.UTC), 10, ProjectCalendar)
		if result != expected {
			t.Errorf("Expected %v but got %v", expected, result)
		}
	})

	t.Run("Test Add Date full work week", func(t *testing.T) {
		// AddDate(time.Date(2024, 6, 23, 0, 0, 0, 0, time.UTC), 10, ProjectCalendar{})
		ProjectCalendar := ProjectCalendar{
			ID:   "1",
			Name: "Test",
			WorkDays: []WorkDay{
				{Day: time.Monday, Working: true},
				{Day: time.Tuesday, Working: true},
				{Day: time.Wednesday, Working: true},
				{Day: time.Thursday, Working: true},
				{Day: time.Friday, Working: true},
				{Day: time.Saturday, Working: true},
				{Day: time.Sunday, Working: true}}}
		expected := time.Date(2024, 7, 8, 0, 0, 0, 0, time.UTC)
		result := AddDate(time.Date(2024, 6, 23, 0, 0, 0, 0, time.UTC), 16, ProjectCalendar)
		if result != expected {
			t.Errorf("Expected %v but got %v", expected, result)
		}
	})
}

func TestCalendarExceptions(t *testing.T) {
	t.Run("Test Add Date", func(t *testing.T) {
		// AddDate(time.Date(2024, 6, 23, 0, 0, 0, 0, time.UTC), 10, ProjectCalendar{})
		ProjectCalendar := ProjectCalendar{
			ID:   "1",
			Name: "Test",
			WorkDays: []WorkDay{
				{Day: time.Monday, Working: true},
				{Day: time.Tuesday, Working: true},
				{Day: time.Wednesday, Working: true},
				{Day: time.Thursday, Working: true},
				{Day: time.Friday, Working: true},
				{Day: time.Saturday, Working: false},
				{Day: time.Sunday, Working: false}},
			Exceptions: []ProjectCalendarExepction{
				{StartDate: time.Date(2024, 6, 28, 0, 0, 0, 0, time.UTC), EndDate: time.Date(2024, 6, 28, 0, 0, 0, 0, time.UTC), Working: false},
			}}
		expected := time.Date(2024, 7, 8, 0, 0, 0, 0, time.UTC)
		result := AddDate(time.Date(2024, 6, 23, 0, 0, 0, 0, time.UTC), 10, ProjectCalendar)
		if result != expected {
			t.Errorf("Expected %v but got %v", expected, result)
		}
	})

	t.Run("Test Add Date full work week", func(t *testing.T) {
		// AddDate(time.Date(2024, 6, 23, 0, 0, 0, 0, time.UTC), 10, ProjectCalendar{})
		ProjectCalendar := ProjectCalendar{
			ID:   "1",
			Name: "Test",
			WorkDays: []WorkDay{
				{Day: time.Monday, Working: true},
				{Day: time.Tuesday, Working: true},
				{Day: time.Wednesday, Working: true},
				{Day: time.Thursday, Working: true},
				{Day: time.Friday, Working: true},
				{Day: time.Saturday, Working: true},
				{Day: time.Sunday, Working: true}}}
		expected := time.Date(2024, 7, 8, 0, 0, 0, 0, time.UTC)
		result := AddDate(time.Date(2024, 6, 23, 0, 0, 0, 0, time.UTC), 16, ProjectCalendar)
		if result != expected {
			t.Errorf("Expected %v but got %v", expected, result)
		}
	})
}
