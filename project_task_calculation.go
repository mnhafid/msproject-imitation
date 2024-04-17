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
	Duration     int
	Work         int
	Cost         int
	Predecessors []Predecessor
	Successor    []Successor
	DurationType string
}

type Predecessor struct {
	ID           string
	TaskUniqueID string
	Type         string
	Lag          string
}

type Successor struct {
	ID           string
	TaskUniqueID string
	Type         string
	Lag          string
}

func (t Task) RecalculateDate(predecessorsTask Task, predecessor Predecessor) *Task {
	Lag, _ := time.ParseDuration(predecessor.Lag)
	fmt.Println(Lag)
	if predecessorsTask.Successor != nil {
	}
	switch predecessor.Type {
	case StartToFinish:
		predecessorsTask.EndDate = t.StartDate.Add(Lag)
	case FinishToStart:
		fmt.Println(Lag, t.EndDate, t.StartDate, t.Duration)
		predecessorsTask.StartDate = t.EndDate.Add(Lag)
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
