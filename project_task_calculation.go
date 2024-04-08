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
	Duration     int
	Work         int
	Cost         int
	Predecessors []Predecessor
	DurationType string
}

type Predecessor struct {
	ID           string
	TaskUniqueID string
	Type         string
	Lag          string
}

func (t Task) RecalculateDate(predecessorsTask Task, predecessor Predecessor) *Task {
	Lag, _ := time.ParseDuration(predecessor.Lag)
	switch predecessor.Type {
	case StartToFinish:
		predecessorsTask.EndDate = t.StartDate.Add(Lag)
	case FinishToStart:
		predecessorsTask.StartDate = t.EndDate.Add(Lag)
	case StartToStart:
		predecessorsTask.StartDate = t.StartDate.Add(Lag)
	case FinishToFinish:
		predecessorsTask.EndDate = t.EndDate.Add(Lag)
	}
	return &predecessorsTask
}
