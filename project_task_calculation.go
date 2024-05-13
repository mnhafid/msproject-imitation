package main

import (
	"fmt"
	"time"

	"google.golang.org/genproto/googleapis/type/decimal"
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
	ActualStart  time.Time
	ActualFinish time.Time
	PlanProgress decimal.Decimal
	CriticalPath bool
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
	DefaultRestHour  = 1
	DefaultWorkHour  = float64(DefaultEndHour - DefaultStartHour - DefaultRestHour)
	HoursPerWeek     = 40
	DaysPerMonth     = 20
)

func CalculateStartFinish(pTask Task, predTask Task) (start time.Time, finish time.Time) {
	if len(pTask.Predecessors) == 0 {
		return pTask.StartDate, pTask.EndDate
	}
	if pTask.StartDate.After(predTask.StartDate) {
		return pTask.StartDate, pTask.EndDate
	}
	start = pTask.StartDate
	finish = pTask.EndDate
	predecessor := pTask.Predecessors[0]
	Lag, _ := ParseDuration(predecessor.Lag)
	workingDay := CalculateWorkingDay(Lag)

	switch predecessor.Type {
	case FinishToStart:
		duration, _ := ParseDuration(pTask.Duration)
		finish = pTask.StartDate.Add(duration)
		if pTask.StartDate.Before(predTask.EndDate) {
			start = predTask.EndDate.Add(Lag)
			finish = pTask.StartDate.Add(duration)
		}
		if DefaultDuration == DurationTypeDay {
			start = time.Date(predTask.StartDate.Year(), predTask.StartDate.Month(), predTask.StartDate.Day(), DefaultStartHour, 0, 0, 0, time.UTC).AddDate(0, 0, 1).Add(workingDay)
			finish = time.Date(predTask.StartDate.Year(), predTask.StartDate.Month(), predTask.StartDate.Day(), DefaultEndHour, 0, 0, 0, time.UTC).Add(duration)
			if pTask.ID == "4" {
				fmt.Println("StartToFinish", start, finish, pTask.Description, Lag, pTask.StartDate, pTask.StartDate.Before(predTask.StartDate), predTask.Description, predTask.StartDate, predTask.EndDate)
			}
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

func CalculateWorkingDay(lag time.Duration) time.Duration {
	var workingDay time.Duration
	if DefaultDuration == DurationTypeDay {
		lagWork := lag.Hours() / DefaultWorkHour * 24
		workingDay = time.Duration(int64(lagWork) * int64(time.Hour))
	}

	return workingDay

}
