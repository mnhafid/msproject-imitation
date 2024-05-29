package main

import (
	"testing"
	"time"
)

func TestTaskWithLagMonth(t *testing.T) {

	TaskResponse := make([]TaskResponse, 2)

	t.Run("Test Success Calculte With Lag = 3 months ", func(t *testing.T) {
		newStart := time.Date(2024, 11, 06, 8, 0, 0, 0, time.UTC)
		TaskResponse[0].Start = newStart.Format("2006-01-02T15:04")
		TaskResponse[0].Duration = "2.0d"
		newDuration, _ := ParseDuration(TaskResponse[0].Duration)
		TaskResponse[0].Finish = CalculateWorkingDay(newStart, newDuration).Format("2006-01-02T15:04")
		TaskResponse[0].Name = "Task 1"
		TaskResponse[0].UniqueID = "1"
		TaskResponse[0].ID = "1"
		TaskResponse[1].Name = "Task 2"
		TaskResponse[1].UniqueID = "2"
		TaskResponse[1].ID = "2"
		TaskResponse[1].Start = time.Now().Format("2006-01-02T15:04")
		TaskResponse[1].Duration = "2.0d"
		newDuration, _ = ParseDuration(TaskResponse[1].Duration)
		TaskResponse[1].Finish = CalculateWorkingDay(time.Now(), newDuration).Format("2006-01-02T15:04")
		TaskResponse[1].Predecessors = []Predecessor{
			{
				ID:       "1",
				UniqueID: "1",
				Type:     "FS",
				Lag:      "3.5mo",
			},
		}
		result := CalculateAllProjectTask(TaskResponse)
		if result[1].StartDate.Truncate(24*time.Hour) != result[0].EndDate.AddDate(0, 3, 1).Truncate(24*time.Hour) {
			t.Errorf("Test %v: Expected  %v but got %v", result[1].Description, result[0].EndDate.AddDate(0, 3, 1).Truncate(24*time.Hour), result[1].StartDate.Truncate(24*time.Hour))
		}

	})
}

func TestTaskWithLagWeek(t *testing.T) {

	TaskResponse := make([]TaskResponse, 2)

	t.Run("Test Success Calculte With Lag = 1 week ", func(t *testing.T) {
		newStart := time.Date(2024, 11, 06, 8, 0, 0, 0, time.UTC)
		TaskResponse[0].Start = newStart.Format("2006-01-02T15:04")
		TaskResponse[0].Duration = "2.0d"
		newDuration, _ := ParseDuration(TaskResponse[0].Duration)
		TaskResponse[0].Finish = CalculateWorkingDay(newStart, newDuration).Format("2006-01-02T15:04")
		TaskResponse[0].Name = "Task 1"
		TaskResponse[0].UniqueID = "1"
		TaskResponse[0].ID = "1"
		TaskResponse[1].Name = "Task 2"
		TaskResponse[1].UniqueID = "2"
		TaskResponse[1].ID = "2"
		TaskResponse[1].Start = time.Now().Format("2006-01-02T15:04")
		TaskResponse[1].Duration = "2.0d"
		newDuration, _ = ParseDuration(TaskResponse[1].Duration)
		TaskResponse[1].Finish = CalculateWorkingDay(time.Now(), newDuration).Format("2006-01-02T15:04")
		TaskResponse[1].Predecessors = []Predecessor{
			{
				ID:       "1",
				UniqueID: "1",
				Type:     "FS",
				Lag:      "1.0w",
			},
		}
		result := CalculateAllProjectTask(TaskResponse)
		if result[1].StartDate.Truncate(24*time.Hour) != result[0].EndDate.AddDate(0, 0, 8).Truncate(24*time.Hour) {
			t.Errorf("Test %v: Expected  %v but got %v", result[1].Description, result[0].EndDate.AddDate(0, 0, 8).Truncate(24*time.Hour), result[1].StartDate.Truncate(24*time.Hour))
		}
	})
}
