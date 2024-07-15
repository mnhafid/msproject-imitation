package main

import (
	"testing"
	"time"
)

func TestEtc(t *testing.T) {

	t.Run("Test Success ETC Path FS", func(t *testing.T) {
		TaskCriticalPath := []Task{
			{UniqueID: "1",
				Duration:       10,
				Successors:     []Successor{{UniqueID: "2", Lag: "2.0d", Type: "FS"}, {UniqueID: "6", Lag: "0", Type: "FS"}},
				Predecessors:   []Predecessor{},
				ActualProgress: 33.3333,
				StartDate:      time.Date(2024, 6, 23, 0, 0, 0, 0, time.UTC),
				EndDate:        time.Date(2024, 8, 3, 8, 0, 0, 0, time.UTC),
				ActualStart:    time.Date(2024, 6, 26, 8, 0, 0, 0, time.UTC),
				ActualFinish:   time.Date(2024, 6, 26, 8, 0, 0, 0, time.UTC),
				StartEtc:       time.Date(2024, 6, 23, 8, 0, 0, 0, time.UTC),
				Etc:            time.Date(2024, 8, 3, 8, 0, 0, 0, time.UTC),
			},
			{UniqueID: "2",
				Duration:       5,
				Successors:     []Successor{{UniqueID: "3", Lag: "2.0d", Type: "FS"}},
				Predecessors:   []Predecessor{{UniqueID: "1", Lag: "2.0d", Type: "FS"}},
				StartDate:      time.Date(2024, 8, 6, 0, 0, 0, 0, time.UTC),
				EndDate:        time.Date(2024, 9, 3, 8, 0, 0, 0, time.UTC),
				StartEtc:       time.Date(2024, 8, 6, 0, 0, 0, 0, time.UTC),
				Etc:            time.Date(2024, 9, 3, 8, 0, 0, 0, time.UTC),
				ActualProgress: 0,
			},
			{UniqueID: "3",
				Duration:       5,
				Successors:     []Successor{},
				Predecessors:   []Predecessor{{UniqueID: "2", Lag: "2.0d", Type: "FS"}},
				StartDate:      time.Date(2024, 9, 6, 0, 0, 0, 0, time.UTC),
				EndDate:        time.Date(2024, 9, 30, 8, 0, 0, 0, time.UTC),
				StartEtc:       time.Date(2024, 9, 6, 0, 0, 0, 0, time.UTC),
				Etc:            time.Date(2024, 9, 30, 8, 0, 0, 0, time.UTC),
				ActualProgress: 0,
			},
		}

		taskIndices := make(map[string]int)
		for i := 0; i < len(TaskCriticalPath); i++ {
			taskIndices[TaskCriticalPath[i].UniqueID] = i
		}
		result := []time.Time{time.Date(2024, 7, 23, 8, 0, 0, 0, time.UTC),
			time.Date(2024, 8, 23, 8, 0, 0, 0, time.UTC),
			time.Date(2024, 9, 19, 8, 0, 0, 0, time.UTC)}

		TaskCriticalPath = PrepareEtc(TaskCriticalPath, taskIndices)
		for i := range TaskCriticalPath {
			if TaskCriticalPath[i].Etc.Truncate(24*time.Hour) != result[i].Truncate(24*time.Hour) {
				t.Errorf("[TASK]: %v Expected %v but got %v", TaskCriticalPath[i].UniqueID, result[i], TaskCriticalPath[i].Etc)
			}
		}
	})
}

func TestEtcCustomCalendar(t *testing.T) {

	t.Run("Test Success ETC Path FS", func(t *testing.T) {
		ProjectCalendars := []ProjectCalendar{
			{
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
					{
						StartDate: time.Date(2024, 6, 28, 0, 0, 0, 0, time.UTC),
						EndDate:   time.Date(2024, 6, 28, 0, 0, 0, 0, time.UTC),
						Working:   false},
				},
			},
			{
				ID:   "2",
				Name: "Test",
				WorkDays: []WorkDay{
					{Day: time.Monday, Working: true},
					{Day: time.Tuesday, Working: true},
					{Day: time.Wednesday, Working: true},
					{Day: time.Thursday, Working: true},
					{Day: time.Friday, Working: true},
					{Day: time.Saturday, Working: true},
					{Day: time.Sunday, Working: true}},
				Exceptions: []ProjectCalendarExepction{
					{
						StartDate: time.Date(2024, 6, 28, 0, 0, 0, 0, time.UTC),
						EndDate:   time.Date(2024, 6, 28, 0, 0, 0, 0, time.UTC),
						Working:   false},
				},
			},
		}
		TaskCriticalPath := []Task{
			{UniqueID: "1",
				Duration:        10,
				Successors:      []Successor{},
				Predecessors:    []Predecessor{},
				ActualProgress:  33.3333,
				StartDate:       time.Date(2024, 7, 11, 0, 0, 0, 0, time.UTC),
				EndDate:         time.Date(2024, 8, 3, 0, 0, 0, 0, time.UTC),
				ActualStart:     time.Date(2024, 6, 26, 8, 0, 0, 0, time.UTC),
				ActualFinish:    time.Date(2024, 6, 26, 8, 0, 0, 0, time.UTC),
				ProjectCalendar: ProjectCalendars[0],
			},
			{UniqueID: "2",
				Duration:        10,
				Successors:      []Successor{},
				Predecessors:    []Predecessor{},
				StartDate:       time.Date(2024, 8, 6, 0, 0, 0, 0, time.UTC),
				EndDate:         time.Date(2024, 9, 3, 8, 0, 0, 0, time.UTC),
				ActualProgress:  0,
				ProjectCalendar: ProjectCalendars[0]},
			{UniqueID: "3",
				Successors:      []Successor{},
				Predecessors:    []Predecessor{},
				StartDate:       time.Date(2024, 9, 6, 0, 0, 0, 0, time.UTC),
				EndDate:         time.Date(2024, 9, 30, 8, 0, 0, 0, time.UTC),
				ActualProgress:  0,
				ProjectCalendar: ProjectCalendars[1]},
			{UniqueID: "4",
				Successors:      []Successor{},
				Predecessors:    []Predecessor{},
				ActualProgress:  33.3333,
				StartDate:       time.Date(2024, 6, 23, 0, 0, 0, 0, time.UTC),
				EndDate:         time.Date(2024, 8, 3, 8, 0, 0, 0, time.UTC),
				ProjectCalendar: ProjectCalendars[1],
			},
		}

		taskIndices := make(map[string]int)
		for i := 0; i < len(TaskCriticalPath); i++ {
			taskIndices[TaskCriticalPath[i].UniqueID] = i
		}
		result := []time.Time{
			time.Date(2024, 7, 31, 8, 0, 0, 0, time.UTC),
			time.Date(2024, 9, 13, 8, 0, 0, 0, time.UTC),
			time.Date(2024, 9, 30, 8, 0, 0, 0, time.UTC),
			time.Date(2024, 8, 31, 0, 0, 0, 0, time.UTC),
		}

		TaskCriticalPath = PrepareEtc(TaskCriticalPath, taskIndices)
		for i := range TaskCriticalPath {
			if TaskCriticalPath[i].Etc.Truncate(24*time.Hour) != result[i].Truncate(24*time.Hour) {
				t.Errorf("[TASK]: %v Expected %v but got %v", TaskCriticalPath[i].UniqueID, result[i], TaskCriticalPath[i].Etc)
			}
		}
	})
}
