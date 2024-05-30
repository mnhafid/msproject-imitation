package main

import (
	"math"
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
	PlanProgress *decimal.Decimal
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
	UniqueID     string
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
	HourInDay        = 24
	DefaultDuration  = DurationTypeDay
	DefaultStartHour = 8
	DefaultEndHour   = 17
	DefaultRestHour  = float64(1)
	DefaultWorkHour  = float64(DefaultEndHour-DefaultStartHour) - DefaultRestHour
	HoursPerWeek     = 40
	DaysPerMonth     = 20
)

func CalculateStartFinish(pTask Task, predTask Task, predecessor Predecessor) (start time.Time, finish time.Time) {
	if len(pTask.Predecessors) == 0 {
		return pTask.StartDate, pTask.EndDate
	}

	start = pTask.StartDate
	finish = pTask.EndDate
	duration, _ := ParseDuration(pTask.Duration)
	switch predecessor.Type {
	case FinishToStart:
		start = time.Date(predTask.EndDate.Year(), predTask.EndDate.Month(), predTask.EndDate.Day(), DefaultStartHour, 0, 0, 0, time.UTC)
		if DefaultDuration == DurationTypeDay {
			start = start.AddDate(0, 0, 1)
		}
		finish = CalculateWorkingDay(start, duration)
		if predecessor.Lag != "" {
			start, finish = CalculateLag(start, finish, predecessor.Lag)
		}
		if pTask.StartDate.After(start) {
			start = pTask.StartDate
		}
		return
	case StartToFinish:
		finish = predTask.StartDate
		if DefaultDuration == DurationTypeDay {
			start = pTask.EndDate.AddDate(0, 0, -1)
		}
		if predecessor.Lag != "" {
			start, finish = CalculateLag(start, finish, predecessor.Lag)
		}
	case StartToStart:
		finish = predTask.StartDate
		if predecessor.Lag != "" {
			start, finish = CalculateLag(start, finish, predecessor.Lag)
		}
	case FinishToFinish:
		start = predTask.EndDate
		finish = CalculateWorkingDay(start, duration)
		if predecessor.Lag != "" {
			start, finish = CalculateLag(start, finish, predecessor.Lag)
		}
	default:
		return
	}
	return
}

func CalculateWorkingDay(start time.Time, duration time.Duration) time.Time {
	var workingEndDay time.Time
	switch DefaultDuration {
	case DurationTypeDay:
		workingEndDay = time.Date(start.Year(), start.Month(), start.Day(), DefaultEndHour, 0, 0, 0, time.UTC)
		if duration.Hours() < DefaultWorkHour {
			if duration.Hours() > DefaultWorkHour/2 {
				duration = duration + time.Duration(DefaultRestHour*float64(time.Hour))
			}
			workingEndDay = time.Date(start.Year(), start.Month(), start.Day(), DefaultStartHour, 0, 0, 0, time.UTC).Add(duration)
		}
		if duration.Hours() >= DefaultWorkHour {
			roundDuration := math.Ceil(float64(duration.Hours()/24 - 1))
			workingEndDay = time.Date(start.Year(), start.Month(), start.Day(), DefaultEndHour, 0, 0, 0, time.UTC).AddDate(0, 0, int(roundDuration))
		}
	}

	return workingEndDay

}

func CalculateLag(start time.Time, finish time.Time, lagDuration string) (time.Time, time.Time) {
	switch DefaultDuration {
	case DurationTypeDay:
		duration := finish.Truncate(24 * time.Hour).Sub(start.Truncate(24 * time.Hour))
		unit, value := splitLag(lagDuration)
		if unit == "mo" {
			duration := finish.Truncate(24 * time.Hour).Sub(start.Truncate(24 * time.Hour))
			start = time.Date(start.Year(), start.Month(), start.Day(), DefaultStartHour, 0, 0, 0, time.UTC).AddDate(0, int(value), 0)
			finish = CalculateWorkingDay(start, duration)
			return start, finish
		}
		lag, _ := ParseDuration(lagDuration)
		deltaLagDuration := math.Mod(lag.Hours(), DefaultWorkHour)
		if math.Mod(lag.Hours(), float64(HourInDay)) == 0 {
			start = time.Date(start.Year(), start.Month(), start.Day(), DefaultStartHour, 0, 0, 0, time.UTC).AddDate(0, 0, int(lag.Hours()/float64(HourInDay)))
			finish = CalculateWorkingDay(start, duration)
			return start, finish
		}

		if lag.Hours() >= DefaultWorkHour && deltaLagDuration == 0 {
			start = time.Date(start.Year(), start.Month(), start.Day(), DefaultStartHour, 0, 0, 0, time.UTC).AddDate(0, 0, int(lag.Hours()/DefaultWorkHour))
			finish = CalculateWorkingDay(start, time.Duration(math.Mod(lag.Hours(), DefaultWorkHour)))
			return start, finish
		}
		if lag.Hours() > DefaultWorkHour && deltaLagDuration != 0 {
			roundLag := math.Ceil(float64(lag.Hours() / DefaultWorkHour))
			start = time.Date(start.Year(), start.Month(), start.Day(), DefaultStartHour+int(deltaLagDuration), 0, 0, 0, time.UTC).AddDate(0, 0, int(roundLag))
			finish = CalculateWorkingDay(start.AddDate(0, 0, 1), time.Duration(math.Mod(lag.Hours(), DefaultWorkHour)))
			return start, finish
		}
		if lag.Hours() < DefaultWorkHour {
			start = time.Date(start.Year(), start.Month(), start.Day(), DefaultStartHour+int(deltaLagDuration), 0, 0, 0, time.UTC)
			finish = CalculateWorkingDay(start.AddDate(0, 0, 1), time.Duration(math.Mod(lag.Hours(), DefaultWorkHour)))
			return start, finish
		}

	}

	return start, finish

}
