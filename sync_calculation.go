package main

import (
	"fmt"
	"time"
)

type TaskResponses struct {
	Count int
	Tasks []TaskResponse `json:"task"`
}

type TaskResponse struct {
	ID           string
	Name         string
	UniqueID     string
	Start        string
	Finish       string
	Duration     string
	Work         int
	Cost         int
	Predecessors []Predecessor `json:"predecessors"`
	Successors   []Successor   `json:"successors"`
	DurationType string
}

const (
	formatTime = "2006-01-02T15:04"
)

func CalculateAllProjectTask(tasksReponse []TaskResponse) []Task {
	var tasks []Task
	// Calculate all project task
	for i := 0; i < len(tasksReponse); i++ {
		start, _ := time.Parse(formatTime, tasksReponse[i].Start)
		finish, _ := time.Parse(formatTime, tasksReponse[i].Finish)
		task := Task{
			ID:           tasksReponse[i].ID,
			Description:  tasksReponse[i].Name,
			UniqueID:     tasksReponse[i].UniqueID,
			Duration:     tasksReponse[i].Duration,
			Work:         tasksReponse[i].Work,
			Cost:         tasksReponse[i].Cost,
			DurationType: tasksReponse[i].DurationType,
			Predecessors: tasksReponse[i].Predecessors,
			Successors:   tasksReponse[i].Successors,
			StartDate:    start,
			EndDate:      finish,
		}
		tasks = append(tasks, task)
	}
	tasks[3].EndDate.AddDate(0, 0, 3)
	result := RecalculateDate(tasks[16], tasks[3])
	fmt.Println(result.ID, result.StartDate)
	tasks[16] = *result
	return tasks
}
