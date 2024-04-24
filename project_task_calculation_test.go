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
		2024, 01, 02, 8, 00, 00, 000000, time.UTC)
	taskParent.EndDate = time.Date(
		2024, 01, 04, 5, 00, 00, 000000, time.UTC)
	taskParent.Duration = 2
	taskParent.UniqueID = "1"

	var taskPredessor Task
	taskPredessor.ID = "2"
	taskPredessor.Description = "Task 2"
	taskPredessor.StartDate = time.Date(
		2024, 01, 03, 8, 00, 00, 000000, time.UTC)
	taskPredessor.EndDate = time.Date(
		2024, 01, 15, 5, 00, 00, 000000, time.UTC)
	taskPredessor.Duration = 13
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
		result := taskParent.RecalculateDate(taskPredessor, taskParent.Predecessors[0])
		// End date of predessor must be equal to start date of task
		if result.StartDate != taskParent.EndDate.Add(24*time.Hour) {
			t.Errorf("got %v want %v", result.StartDate, taskParent.EndDate)
		}
	})

	t.Run("Test Success Set Predessor Start to Finish", func(t *testing.T) {

		taskParent.Predecessors = []Predecessor{
			{
				ID:       "2",
				Type:     StartToFinish,
				Lag:      "0.0d",
				UniqueID: "2",
			},
		}

		result := taskParent.RecalculateDate(taskPredessor, taskParent.Predecessors[0])
		if result.EndDate != taskParent.StartDate {
			t.Errorf("got %v want %v", result.EndDate, taskParent.StartDate)
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
		result := taskParent.RecalculateDate(taskPredessor, taskParent.Predecessors[0])
		// FinishToStart[FS]
		// Start date of predessor must be equal to end date of task
		if result.StartDate != taskParent.StartDate {
			t.Errorf("got %v want %v", result.StartDate, taskParent.StartDate)
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
		result := taskParent.RecalculateDate(taskPredessor, taskParent.Predecessors[0])
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
		2023, 01, 02, 8, 00, 00, 000000, time.UTC)
	taskParent.EndDate = time.Date(
		2023, 01, 02, 5, 00, 00, 000000, time.UTC)
	taskParent.Duration = 2
	taskParent.UniqueID = "1"

	var taskPredessor Task
	taskPredessor.ID = "2"
	taskPredessor.Description = "Task 2"
	taskPredessor.StartDate = time.Date(
		2024, 01, 03, 8, 00, 00, 000000, time.UTC)
	taskPredessor.EndDate = time.Date(
		2024, 01, 15, 5, 00, 00, 000000, time.UTC)
	taskPredessor.Duration = 2
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
			2023, 04, 03, 5, 00, 00, 000000, time.UTC)
		taskParent.EndDate = want
		result := taskParent.RecalculateDate(taskPredessor, taskParent.Predecessors[0])
		want = want.Add(72 * time.Hour)
		// End date of predessor must be equal to start date of task
		if result.StartDate != want {
			t.Errorf("got %v want %v", result.StartDate, want)
		}
	})

	t.Run("Test Success Multiple Predecessors Finish To Start", func(t *testing.T) {
		taskPredecessor := []Task{
			{
				ID:          "3",
				Description: "Task 3",
				StartDate: time.Date(
					2024, 01, 10, 5, 8, 00, 000000, time.UTC),
				EndDate: time.Date(
					2024, 01, 15, 5, 5, 00, 000000, time.UTC),
				Duration: 5,
				UniqueID: "3",
			},
			{
				ID:          "4",
				Description: "Task 4",
				StartDate: time.Date(
					2024, 01, 20, 5, 8, 00, 000000, time.UTC),
				EndDate: time.Date(
					2024, 01, 25, 5, 5, 00, 000000, time.UTC),
				Duration: 5,
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
		// Add first predessor
		result := taskParent.RecalculateDate(taskPredecessor[0], taskParent.Predecessors[0])
		want := time.Date(
			2023, 04, 17, 8, 00, 00, 000000, time.UTC)
		// End date of predessor must be equal to start date of task
		if result.StartDate != want {
			t.Errorf("got %v want %v", result.StartDate, want)
		}
		// Second predessor
		resultFlnal := taskParent.RecalculateDate(taskPredecessor[1], taskParent.Predecessors[0])
		wantFinal := time.Date(
			2023, 04, 27, 8, 00, 00, 000000, time.UTC)
		// End date of predessor must be equal to start date of task
		if resultFlnal.StartDate != wantFinal {
			t.Errorf("got %v want %v", result.StartDate, wantFinal)
		}
	})
}
