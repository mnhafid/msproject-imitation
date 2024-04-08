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
		2024, 01, 10, 5, 00, 00, 000000, time.UTC)
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
				ID:           "2",
				Type:         FinishToStart,
				Lag:          "0",
				TaskUniqueID: "2",
			},
		}
		result := taskParent.RecalculateDate(taskPredessor, taskParent.Predecessors[0])
		// End date of predessor must be equal to start date of task
		if result.StartDate != taskParent.EndDate {
			t.Errorf("got %v want %v", result.EndDate, taskParent.StartDate)
		}
	})

	t.Run("Test Success Set Predessor Start to Finish", func(t *testing.T) {

		taskParent.Predecessors = []Predecessor{
			{
				ID:           "2",
				Type:         StartToFinish,
				Lag:          "0",
				TaskUniqueID: "2",
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
				ID:           "2",
				Type:         StartToStart,
				Lag:          "0",
				TaskUniqueID: "2",
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
				ID:           "2",
				Type:         FinishToFinish,
				Lag:          "0",
				TaskUniqueID: "2",
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
		2024, 01, 02, 8, 00, 00, 000000, time.UTC)
	taskParent.EndDate = time.Date(
		2024, 01, 10, 5, 00, 00, 000000, time.UTC)
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
	t.Run("Test Success Set Predecessors FS with Lag", func(t *testing.T) {
		taskParent.Predecessors = []Predecessor{
			{
				ID:           "2",
				Type:         FinishToStart,
				Lag:          "24d",
				TaskUniqueID: "2",
			},
		}
		result := taskParent.RecalculateDate(taskPredessor, taskParent.Predecessors[0])
		// End date of predessor must be equal to start date of task
		if result.StartDate != taskParent.EndDate.Add(time.Duration(24*time.Hour)) {
			t.Errorf("got %v want %v", result.StartDate, taskParent.EndDate.Add(time.Duration(24*time.Hour)))
		}
	})
}
