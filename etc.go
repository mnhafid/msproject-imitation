package main

import (
	"encoding/json"
	"fmt"
	"math"
	"time"
)

func CalculateEtcTask(tasksReponse []ProjectTasksProgress) []ProjectTasksProgress {
	var tasks []Task
	taskIndices := make(map[string]int)
	for i := 0; i < len(tasksReponse); i++ {
		// calculate duration
		duration := tasksReponse[i].EndDate.Sub(*tasksReponse[i].StartDate)
		durationDays := math.Ceil(duration.Hours()/24) + 1
		var predecessor []Predecessor
		var successors []Successor
		json.Unmarshal(tasksReponse[i].Predecessor.Bytes, &predecessor)
		json.Unmarshal(tasksReponse[i].Successor.Bytes, &successors)
		durationMpp := fmt.Sprintf("%v.0d", durationDays)
		taskIndices[*tasksReponse[i].UniqueID] = i

		// build data for calculate critical path
		tasks = append(tasks, Task{
			ID:           *tasksReponse[i].ID,
			UniqueID:     *tasksReponse[i].UniqueID,
			Duration:     math.Ceil(duration.Hours() / 24),
			DurationMpp:  durationMpp,
			Predecessors: predecessor,
			Successors:   successors,
			StartDate:    *tasksReponse[i].StartETC,
			EndDate:      *tasksReponse[i].ETC,
		})

	}

	tasks = calculateCriticalComponent(tasks)
	for i := range tasks {
		index, ok := taskIndices[tasks[i].UniqueID]
		if ok {
			tasksReponse[index].IsCriticalPath = &tasks[i].CriticalPath
		}
	}
	return tasksReponse
}

func PrepareEtc(tasks []Task, taskIndices map[string]int) []Task {
	for i := range tasks {
		duration := tasks[i].EndDate.Sub(tasks[i].StartDate)
		durationDays := math.Ceil(duration.Hours() / 24)
		// set data for old data
		if tasks[i].StartEtc.IsZero() {
			tasks[i].StartEtc = tasks[i].StartDate
		}

		// Set actual progress and set default if task not have not started
		progresTask := 0.0
		if tasks[i].ActualProgress != 0 {
			progresTask = tasks[i].ActualProgress
		}

		startEtc := tasks[i].StartDate

		for _, pred := range tasks[i].Predecessors {
			lag := float64(0)
			if pred.Lag != "0" {
				unit, value := splitLag(pred.Lag)
				lag = float64(value)
				if unit == "mo" {
					lag = float64(value) * 30
				}
			}
			switch pred.Type {
			case "FS":
				if tasks[taskIndices[pred.UniqueID]].Etc.After(startEtc) {
					startEtc = tasks[taskIndices[pred.UniqueID]].Etc
					continue
				}
				startEtc = tasks[taskIndices[pred.UniqueID]].Etc
				if startEtc.Before(tasks[i].StartDate) {
					startEtc = startEtc.AddDate(0, 0, int(lag))
				}
			case "FF":
			case "SF":
			case "SS":
			}

		}
		ETC := CalculateEtc(durationDays, startEtc, tasks[i].StartDate, progresTask)
		tasks[i].StartEtc = startEtc
		tasks[i].Etc = ETC
	}

	return tasks

}
func CalculateEtc(duration float64, startEtc time.Time, planDate time.Time, actualProgressItd float64) time.Time {
	ratio := math.Pow(10, float64(2))
	actualDurationItd := math.Ceil(time.Since(planDate).Hours() / 24)
	// default pace using pace Plan
	pace := 100 / duration
	fmt.Println("pace", pace, duration)
	// if have actual progress, calculate pace using actual progress
	if actualProgressItd != 0 {
		pace = actualProgressItd / actualDurationItd
	}
	estimateDuration := int(math.Ceil(math.Round((100/pace)*ratio) / ratio))
	// minus 1 because we want start from day 0
	etc := startEtc.AddDate(0, 0, estimateDuration)
	return etc
}
