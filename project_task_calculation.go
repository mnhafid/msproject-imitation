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

func (t Task) RecalculateDate(predecessorsTask Task, predecessor Predecessor) *Task {
	Lag := time.Duration(0 * time.Hour)
	Lag, _ = ParseDuration(predecessor.Lag)
	switch predecessor.Type {
	case StartToFinish:
		predecessorsTask.EndDate = t.StartDate.Add(Lag)
	case FinishToStart:
		predecessorsTask.StartDate = t.EndDate.Add(Lag + 24*time.Hour)
		predecessorsTask.EndDate = predecessorsTask.StartDate.Add(time.Duration(predecessorsTask.Duration) * time.Hour)
	case StartToStart:
		predecessorsTask.StartDate = t.StartDate.Add(Lag)
	case FinishToFinish:
		predecessorsTask.EndDate = t.EndDate.Add(Lag)
	default:
		return &predecessorsTask
	}

	predecessorsTask.Successor = append(predecessorsTask.Successor, Successor{
		ID:           t.ID,
		TaskUniqueID: t.UniqueID,
		Type:         predecessor.Type,
		Lag:          predecessor.Lag,
	})

	return &predecessorsTask
}
