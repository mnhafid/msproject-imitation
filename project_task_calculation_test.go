package main

import (
	"testing"
	"time"
)

func TestProjectTaskCalculation(t *testing.T) {
	var taskParent Task
	taskParent.ID = "1"
	taskParent.Description = "Task 1"
	taskParent.StartDate = time.Date(
		2024, 04, 25, 8, 00, 00, 000000, time.UTC)
	taskParent.EndDate = time.Date(
		2024, 04, 25, 16, 00, 00, 000000, time.UTC)
	taskParent.Duration = "2.0d"
	taskParent.UniqueID = "1"

	var taskPredessor Task
	taskPredessor.ID = "2"
	taskPredessor.Description = "Task 2"
	taskPredessor.StartDate = time.Date(
		2024, 04, 25, 8, 00, 00, 000000, time.UTC)
	taskPredessor.EndDate = time.Date(
		2024, 04, 25, 16, 00, 00, 000000, time.UTC)
	taskPredessor.Duration = "13.0d"
	taskPredessor.UniqueID = "2"
	t.Run("Test Success Set Predecessors Finish To Start", func(t *testing.T) {
		taskParent.Predecessors = []Predecessor{
			{
				ID:       "2",
				Type:     FinishToStart,
				Lag:      "0.0d",
				UniqueID: "2",
			},
		}
		want := taskParent.EndDate.Add(24 * time.Hour)
		result := RecalculateDate(taskParent, taskPredessor)
		// End date of predessor must be equal to start date of task
		if result.StartDate != want {
			t.Errorf("got %v want %v", result.StartDate, want)
		}
	})

	t.Run("Test Success Set Predecessors Start To Finish", func(t *testing.T) {
		taskParent.Predecessors = []Predecessor{
			{
				ID:       "2",
				Type:     StartToFinish,
				Lag:      "0.0d",
				UniqueID: "2",
			},
		}
		want := taskPredessor.StartDate
		result := RecalculateDate(taskParent, taskPredessor)
		// End date of predessor must be equal to start date of task
		if result.EndDate != want {
			t.Errorf("EndDate got %v want %v", result.EndDate, want)
		}

		if result.StartDate != want.Add(-24*time.Hour) {
			t.Errorf("StartDate got %v want %v", result.StartDate, want.Add(-24*time.Hour))
		}

	})

	t.Run("Test Success Set Predessor Start To Start", func(t *testing.T) {
		taskParent.Predecessors = []Predecessor{
			{
				ID:       "2",
				Type:     StartToStart,
				Lag:      "0.0d",
				UniqueID: "2",
			},
		}
		want := taskPredessor.StartDate
		result := RecalculateDate(taskParent, taskPredessor)
		// End date of predessor must be equal to start date of task
		if result.StartDate != want {
			t.Errorf("StartDate got %v want %v", result.StartDate, want)
		}
	})
	t.Run("Test Success Set Predessor Finish To Finsih", func(t *testing.T) {
		taskParent.Predecessors = []Predecessor{
			{
				ID:       "2",
				Type:     FinishToFinish,
				Lag:      "0.0d",
				UniqueID: "2",
			},
		}
		result := RecalculateDate(taskParent, taskPredessor)
		// FinishtToFinish[FF]
		// Start date of predessor must be equal to end date of task
		if result.EndDate != taskParent.EndDate {
			t.Errorf("got %v want %v", result.EndDate, taskParent.EndDate)
		}

	})
}

func TestProjectTaskCalculationWithLag(t *testing.T) {
	var taskParent Task
	taskParent.ID = "1"
	taskParent.Description = "Task 1"
	taskParent.StartDate = time.Date(
		2024, 01, 03, 8, 00, 00, 000000, time.UTC)
	taskParent.EndDate = time.Date(
		2024, 01, 15, 16, 00, 00, 000000, time.UTC)
	taskParent.Duration = "2.0d"
	taskParent.UniqueID = "1"

	var taskPredessor Task
	taskPredessor.ID = "2"
	taskPredessor.Description = "Task 2"
	taskPredessor.StartDate = time.Date(
		2024, 01, 03, 8, 00, 00, 000000, time.UTC)
	taskPredessor.EndDate = time.Date(
		2024, 01, 15, 16, 00, 00, 000000, time.UTC)
	taskPredessor.Duration = "2.0d"
	taskPredessor.UniqueID = "2"
	t.Run("Test Success Set Predecessors Finish To Start", func(t *testing.T) {
		taskParent.Predecessors = []Predecessor{
			{
				ID:       "2",
				Type:     FinishToStart,
				Lag:      "2.0d",
				UniqueID: "2",
			},
		}
		want := time.Date(
			2024, 01, 15, 16, 00, 00, 000000, time.UTC)
		taskParent.EndDate = want
		result := RecalculateDate(taskParent, taskPredessor)
		want = want.Add(72 * time.Hour)
		// End date of predessor must be equal to start date of task
		if result.StartDate != want {
			t.Errorf("got %v want %v", result.StartDate, want)
		}
	})
}
func TestProjectTaskCalculationMultiple(t *testing.T) {
	t.Run("Test Success Multiple Predecessors Finish To Start", func(t *testing.T) {
		var taskParent Task
		taskParent.ID = "1"
		taskParent.Description = "Task 1"
		taskParent.StartDate = time.Date(
			2024, 04, 25, 8, 00, 00, 000000, time.UTC)
		taskParent.EndDate = time.Date(
			2024, 04, 25, 16, 00, 00, 000000, time.UTC)
		taskParent.Duration = "2.0d"
		taskParent.UniqueID = "1"

		taskPredecessor := []Task{
			{
				ID:          "3",
				Description: "Task 3",
				StartDate: time.Date(
					2024, 04, 25, 8, 0, 00, 000000, time.UTC),
				EndDate: time.Date(
					2024, 01, 25, 16, 0, 00, 000000, time.UTC),
				Duration: "1.0d",
				UniqueID: "3",
			},
			{
				ID:          "4",
				Description: "Task 4",
				StartDate: time.Date(
					2024, 04, 27, 8, 0, 00, 000000, time.UTC),
				EndDate: time.Date(
					2024, 04, 27, 16, 0, 00, 000000, time.UTC),
				Duration: "1.0d",
				UniqueID: "4",
			},
		}

		taskParent.Predecessors = []Predecessor{
			{
				ID:       "3",
				Type:     FinishToStart,
				Lag:      "0.0d",
				UniqueID: "3",
			},
			{
				ID:       "4",
				Type:     FinishToStart,
				Lag:      "0.0d",
				UniqueID: "4",
			},
		}

		startDate := time.Date(
			2023, 04, 27, 8, 00, 00, 000000, time.UTC)
		taskParent.EndDate = startDate.Add(8 * time.Hour)

		// Add first predessor
		result := RecalculateDate(taskParent, taskPredecessor[0])
		result = RecalculateDate(*result, taskPredecessor[1])
		// End date of predessor must be equal to start date of task
		if result.StartDate != taskPredecessor[1].EndDate.Add(24*time.Hour) {
			t.Errorf("got %v want %v", result.StartDate, taskPredecessor[1].EndDate.Add(24*time.Hour))
		}

	})
}
