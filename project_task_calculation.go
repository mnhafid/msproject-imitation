package main

import (
	"fmt"
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
	Predecessors []Predecessor `json:"predecessors"`
	Successors   []Successor   `json:"successors"`
	DurationType string
	PlanStart    time.Time
	PlanFinish   time.Time
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

func RecalculateDate(pTask Task, predTask Task) Task {
	if len(pTask.Predecessors) == 0 {
		return pTask
	}
	fmt.Printf("Task %s Name %s\n", pTask.ID, pTask.Description)
	predecessor := pTask.Predecessors[0]
	Lag, _ := ParseDuration(predecessor.Lag)
	switch predecessor.Type {
	case FinishToStart:
		fmt.Printf("Found Predessesor Type [FS]: Start Date Before: %v\n", pTask.StartDate)
		duration, _ := ParseDuration(pTask.Duration)
		pTask.EndDate = pTask.StartDate.Add(duration)
		if pTask.StartDate.Before(predTask.EndDate) {
			pTask.StartDate = predTask.EndDate.Add(Lag)
			pTask.EndDate = pTask.StartDate.Add(duration)
		}
		if DefaultDuration == DurationTypeDay {
			pTask.StartDate = predTask.EndDate.AddDate(0, 0, 1).Add(Lag)
		}
		fmt.Printf("Found Predessesor Type [FS]: Start Date After: %v\n", pTask.StartDate)
	case StartToFinish:
		pTask.EndDate = predTask.StartDate
		if DefaultDuration == DurationTypeDay {
			pTask.StartDate = pTask.EndDate.AddDate(0, 0, -1).Add(Lag)
		}
	case StartToStart:
		pTask.StartDate = predTask.StartDate.Add(Lag)
	case FinishToFinish:
	default:
		return pTask
	}
	return pTask
}

func CalculateStartFinish(pTask Task, predTask Task) (start time.Time, finish time.Time) {
	if len(pTask.Predecessors) == 0 {
		return pTask.StartDate, pTask.EndDate
	}
	if pTask.StartDate.Before(predTask.StartDate) {
		return pTask.StartDate, pTask.EndDate
	}
	start = pTask.StartDate
	finish = pTask.EndDate
	predecessor := pTask.Predecessors[0]
	Lag, _ := ParseDuration(predecessor.Lag)
	switch predecessor.Type {
	case FinishToStart:
		duration, _ := ParseDuration(pTask.Duration)
		finish = pTask.StartDate.Add(duration)
		if pTask.StartDate.Before(predTask.EndDate) {
			start = predTask.EndDate.Add(Lag)
			finish = pTask.StartDate.Add(duration)
		}
		if DefaultDuration == DurationTypeDay {
			start = predTask.EndDate.AddDate(0, 0, 1).Add(Lag)
		}
	case StartToFinish:
		finish = predTask.StartDate
		if DefaultDuration == DurationTypeDay {
			start = pTask.EndDate.AddDate(0, 0, -1).Add(Lag)
		}
	case StartToStart:
		finish = predTask.StartDate.Add(Lag)
	case FinishToFinish:
		pTask.EndDate = predTask.EndDate.Add(Lag)
	default:
		return
	}
	return
}
