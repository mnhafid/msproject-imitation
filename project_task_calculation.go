package main

import (
	"time"
)

const (
	StartToFinish  = "SF"
	FinishToStart  = "FS"
	StartToStart   = "SS"
	FinishToFinish = "FF"
)

type Task struct {
	ID           string
	Description  string
	UniqueID     string
	StartDate    time.Time
	EndDate      time.Time
	Duration     string
	Work         int
	Cost         int
	Predecessors []Predecessor
	Successor    []Successor
	DurationType string
}

type Predecessor struct {
	ID       string
	UniqueID string
	Type     string
	Lag      string
}

type Successor struct {
	ID           string
	TaskUniqueID string
	Type         string
	Lag          string
}

const (
	// DurationTypeDay is a constant for duration type day
	DurationTypeDay = "d"
	// DurationTypeHour is a constant for duration type hour
	DurationTypeHour = "h"
	// DurationTypeMinute is a constant for duration type minute
	DurationTypeMinute = "m"
	// DurationTypeSecond is a constant for duration type second
	DurationTypeSecond = "s"

	// DurationTypeWeek is a constant for duration type week
	DurationTypeWeek = "w"
	// DurationTypeMonth is a constant for duration type month
	DurationTypeMonth = "mo"
	// DurationTypeYear is a constant for duration type year
	DurationTypeYear = "y"
)

var (
	DefaultDuration  = DurationTypeDay
	DefaultStartHour = 8
	DefaultEndHour   = 17

	DefaultWorkHour = 8
	HoursPerWeek    = 40
	DaysPerMonth    = 20
)

func RecalculateDate(pTask Task, predTask Task) *Task {
	if len(pTask.Predecessors) == 0 {
		return &pTask
	}

	predecessor := pTask.Predecessors[0]
	Lag, _ := ParseDuration(predecessor.Lag)
	switch predecessor.Type {
	case FinishToStart:
		duration, _ := ParseDuration(pTask.Duration)
		pTask.EndDate = pTask.StartDate.Add(duration)
		if pTask.StartDate.Before(predTask.EndDate) {
			pTask.StartDate = predTask.EndDate.Add(Lag)
			pTask.EndDate = pTask.StartDate.Add(duration)
		}
		if DefaultDuration == DurationTypeDay {
			pTask.StartDate = predTask.EndDate.AddDate(0, 0, 1).Add(Lag)
		}
	case StartToFinish:
		pTask.EndDate = predTask.StartDate
		if DefaultDuration == DurationTypeDay {
			pTask.StartDate = pTask.EndDate.AddDate(0, 0, -1).Add(Lag)
		}
	case StartToStart:
		pTask.StartDate = predTask.StartDate.Add(Lag)
	case FinishToFinish:
	default:
		return &pTask
	}
	return &pTask
}

// 	Lag := time.Duration(0 * time.Hour)
// 	Lag, _ = ParseDuration(pTask.Lag)
// 	switch predecessor.Type {
// 	case StartToFinish:
// 		pTask.EndDate = t.StartDate.Add(Lag)
// 	case FinishToStart:
// 		pTask.StartDate = t.EndDate.Add(Lag)
// 		duration, _ := ParseDuration(pTask.Duration)
// 		pTask.EndDate = pTask.StartDate.Add(duration)
// 	case StartToStart:
// 		pTask.StartDate = t.StartDate.Add(Lag)
// 	case FinishToFinish:
// 		pTask.EndDate = t.EndDate.Add(Lag)
// 	default:
// 		return &pTask
// 	}

// 	pTask.Successor = append(pTask.Successor, Successor{
// 		ID:           t.ID,
// 		TaskUniqueID: t.UniqueID,
// 		Type:         predecessor.Type,
// 		Lag:          predecessor.Lag,
// 	})

// 	return &pTask
// }

func (t Task) CalculateDateTask() *Task {
	duration, _ := ParseDuration(t.Duration)
	t.EndDate = t.StartDate.Add(duration)
	return &t
}
