package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"
	"time"
)

func TestSingleTask(t *testing.T) {
	// Open our jsonFile
	jsonFile, err := os.Open("single_test.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	byteReaader, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
	}

	TaskResponse := TaskResponses{}
	json.Unmarshal(byteReaader, &TaskResponse)
	t.Run("Test Success duration", func(t *testing.T) {
		newStart := time.Date(2024, 05, 07, 8, 0, 0, 0, time.UTC)
		TaskResponse.Tasks[10].Start = newStart.Format("2006-01-02T15:04")
		TaskResponse.Tasks[10].Duration = "1.0d"
		newDuration, _ := ParseDuration(TaskResponse.Tasks[10].Duration)
		TaskResponse.Tasks[10].Finish = time.Date(2024, 05, 07, 17, 0, 0, 0, time.UTC).Add(newDuration).Format("2006-01-02T15:04")
		result := CalculateWorkingDay(newStart, newDuration)
		if result != time.Date(2024, 05, 7, 17, 0, 0, 0, time.UTC) {
			t.Errorf("Expected %v but got %v", time.Date(2024, 05, 7, 17, 0, 0, 0, time.UTC), result)
		}
	})
	t.Run("Test Success 0 duration", func(t *testing.T) {
		newStart := time.Date(2024, 05, 07, 8, 0, 0, 0, time.UTC)
		TaskResponse.Tasks[10].Start = newStart.Format("2006-01-02T15:04")
		TaskResponse.Tasks[10].Duration = "0.0d"
		newDuration, _ := ParseDuration(TaskResponse.Tasks[10].Duration)
		TaskResponse.Tasks[10].Finish = time.Date(2024, 05, 07, 17, 0, 0, 0, time.UTC).Add(newDuration).Format("2006-01-02T15:04")
		result := CalculateWorkingDay(newStart, newDuration)
		if result != time.Date(2024, 05, 7, 8, 0, 0, 0, time.UTC) {
			t.Errorf("Expected %v but got %v", time.Date(2024, 05, 7, 8, 0, 0, 0, time.UTC), result)
		}
	})

	t.Run("Test Success 7.5 duration", func(t *testing.T) {
		newStart := time.Date(2024, 05, 07, 8, 0, 0, 0, time.UTC)
		TaskResponse.Tasks[10].Start = newStart.Format("2006-01-02T15:04")
		TaskResponse.Tasks[10].Duration = "7.5h"
		newDuration, _ := ParseDuration(TaskResponse.Tasks[10].Duration)
		result := CalculateWorkingDay(newStart, newDuration)
		if result != time.Date(2024, 05, 7, 16, 30, 0, 0, time.UTC) {
			t.Errorf("Expected %v but got %v", time.Date(2024, 05, 7, 16, 30, 0, 0, time.UTC), result)
		}
	})
	t.Run("Test Success Calculte Predecessors Finish To Start", func(t *testing.T) {
		newStart := time.Date(2024, 05, 07, 8, 0, 0, 0, time.UTC)
		TaskResponse.Tasks[10].Start = newStart.Format("2006-01-02T15:04")
		TaskResponse.Tasks[10].Duration = "1.0d"
		TaskResponse.Tasks[2].Predecessors[0].Lag = "0.0"
		newDuration, _ := ParseDuration(TaskResponse.Tasks[10].Duration)
		TaskResponse.Tasks[10].Finish = CalculateWorkingDay(newStart, newDuration).Format("2006-01-02T15:04")
		result := CalculateAllProjectTask(TaskResponse.Tasks)
		if result[2].StartDate.Truncate(24*time.Hour) != result[10].EndDate.AddDate(0, 0, 1).Truncate(24*time.Hour) {
			t.Errorf("Test %v: Expected  %v but got %v", result[2].Description, result[10].EndDate.AddDate(0, 0, 1), result[2].StartDate)
		}
		if result[2].EndDate.Truncate(24*time.Hour) != result[10].EndDate.AddDate(0, 0, 1).Truncate(24*time.Hour) {
			t.Errorf("Test %v Finish: Expected  %v but got %v", result[2].Description, result[10].EndDate.AddDate(0, 0, 1), result[2].EndDate)
		}
	})
	t.Run("Test Success Calculte With Lag = 8 hour ", func(t *testing.T) {
		newStart := time.Date(2024, 05, 07, 8, 0, 0, 0, time.UTC)
		TaskResponse.Tasks[10].Start = newStart.Format("2006-01-02T15:04")
		TaskResponse.Tasks[10].Duration = "1.0d"
		newDuration, _ := ParseDuration(TaskResponse.Tasks[10].Duration)
		TaskResponse.Tasks[2].Predecessors[0].Lag = "8.0h"

		TaskResponse.Tasks[10].Finish = CalculateWorkingDay(newStart, newDuration).Format("2006-01-02T15:04")
		result := CalculateAllProjectTask(TaskResponse.Tasks)
		if result[2].StartDate.Truncate(24*time.Hour) != result[10].EndDate.AddDate(0, 0, 2).Truncate(24*time.Hour) {
			t.Errorf("Test %v: Expected  %v but got %v", result[2].Description, result[10].EndDate.AddDate(0, 0, 2), result[2].StartDate)
		}
		if result[2].EndDate.Truncate(24*time.Hour) != result[10].EndDate.AddDate(0, 0, 2).Truncate(24*time.Hour) {
			t.Errorf("Test %v Finish: Expected  %v but got %v", result[2].Description, result[10].EndDate.AddDate(0, 0, 2), result[2].EndDate)
		}
	})

	t.Run("Test Success Calculte With Lag 1 hour ", func(t *testing.T) {
		newStart := time.Date(2024, 05, 07, 8, 0, 0, 0, time.UTC)
		TaskResponse.Tasks[10].Start = newStart.Format("2006-01-02T15:04")
		TaskResponse.Tasks[10].Duration = "1.0d"
		newDuration, _ := ParseDuration(TaskResponse.Tasks[10].Duration)
		TaskResponse.Tasks[2].Predecessors[0].Lag = "1.0h"

		TaskResponse.Tasks[10].Finish = CalculateWorkingDay(newStart, newDuration).Format("2006-01-02T15:04")
		result := CalculateAllProjectTask(TaskResponse.Tasks)

		if result[2].StartDate.Truncate(24*time.Hour) != result[10].EndDate.AddDate(0, 0, 1).Truncate(24*time.Hour) {
			t.Errorf("Test %v: Expected  %v but got %v", result[2].Description, result[10].EndDate.AddDate(0, 0, 1), result[2].StartDate)
		}
		if result[2].EndDate.Truncate(24*time.Hour) != result[10].EndDate.AddDate(0, 0, 2).Truncate(24*time.Hour) {
			t.Errorf("Test %v Finish: Expected  %v but got %v", result[2].Description, result[10].EndDate.AddDate(0, 0, 2), result[2].EndDate)
		}
	})

	t.Run("Test Success Calculte With Lag > 8 hour ", func(t *testing.T) {
		newStart := time.Date(2024, 05, 07, 8, 0, 0, 0, time.UTC)
		TaskResponse.Tasks[10].Start = newStart.Format("2006-01-02T15:04")
		TaskResponse.Tasks[10].Duration = "1.0d"
		newDuration, _ := ParseDuration(TaskResponse.Tasks[10].Duration)
		TaskResponse.Tasks[2].Predecessors[0].Lag = "9.0h"

		TaskResponse.Tasks[10].Finish = CalculateWorkingDay(newStart, newDuration).Format("2006-01-02T15:04")
		result := CalculateAllProjectTask(TaskResponse.Tasks)

		if result[2].StartDate.Truncate(24*time.Hour) != result[10].EndDate.AddDate(0, 0, 3).Truncate(24*time.Hour) {
			t.Errorf("Test %v: Expected  %v but got %v", result[2].Description, result[10].EndDate.AddDate(0, 0, 3), result[2].StartDate)
		}
		if result[2].EndDate.Truncate(24*time.Hour) != result[10].EndDate.AddDate(0, 0, 4).Truncate(24*time.Hour) {
			t.Errorf("Test %v Finish: Expected  %v but got %v", result[2].Description, result[10].EndDate.AddDate(0, 0, 4), result[2].EndDate)
		}
	})
}

// https://stackoverflow.com/questions/77753292/issue-with-mpxj-12-5-0-extracting-or-calculating-duration-of-a-split-task
func TestMultipleTask(t *testing.T) {
	// Open our jsonFile
	jsonFile, err := os.Open("minimum_test.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	byteReaader, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
	}

	TaskResponse := TaskResponses{}
	json.Unmarshal(byteReaader, &TaskResponse)
	t.Run("Test Success Calculate multi dependance Predecessors Finish To Start", func(t *testing.T) {
		newStart := time.Date(2024, 05, 8, 8, 0, 0, 0, time.UTC)
		TaskResponse.Tasks[10].Start = newStart.Format("2006-01-02T15:04")
		TaskResponse.Tasks[10].Duration = "1.0d"
		newDuration, _ := ParseDuration(TaskResponse.Tasks[10].Duration)
		TaskResponse.Tasks[10].Finish = CalculateWorkingDay(newStart, newDuration).Format("2006-01-02T15:04")
		result := CalculateAllProjectTask(TaskResponse.Tasks)
		if result[2].StartDate.Truncate(24*time.Hour) != result[10].EndDate.AddDate(0, 0, 6).Truncate(24*time.Hour) {
			t.Errorf("Test %v: Expected  %v but got %v", result[2].Description,
				result[10].EndDate.AddDate(0, 0, 6).Format("2006-01-02T15:04"),
				result[2].StartDate.Format("2006-01-02T15:04"))
		}
		if result[3].StartDate.Truncate(24*time.Hour) != result[2].EndDate.AddDate(0, 0, 1).Truncate(24*time.Hour) {
			t.Errorf("Test %v Start: Expected %v but got %v", result[3].Description, result[2].EndDate.AddDate(0, 0, 1).Truncate(24*time.Hour).Format("2006-01-02T15:04"), result[3].StartDate.Truncate(24*time.Hour).Format("2006-01-02T15:04"))
		}
		if result[3].EndDate.Truncate(24*time.Hour) != result[2].EndDate.Add(24*time.Hour).Truncate(24*time.Hour) {
			t.Errorf("Test %v Finish: Expected %v but got %v", result[3].Description, result[2].EndDate.Add(24*time.Hour), result[3].EndDate)
		}
		if result[4].StartDate != result[10].StartDate.AddDate(0, 0, 7) {
			t.Errorf("Test %v: Expected %v but got %v", result[4].Description, result[10].StartDate.AddDate(0, 0, 7), result[4].StartDate)
		}
		if result[5].StartDate != result[10].StartDate.AddDate(0, 0, -1) {
			t.Errorf("Test %v:Expected from EndDate %v but got PlanStart %v", result[5].Description, result[10].StartDate.AddDate(0, 0, -1), result[5].StartDate)
		}
		if result[6].StartDate != result[8].StartDate {
			t.Errorf("Test %v:Expected %v but got %v", result[6].Description, result[8].StartDate, result[6].StartDate)
		}
		if result[7].StartDate != result[10].StartDate.AddDate(0, 0, 8) {
			t.Errorf("Test %v:Expected %v but got %v", result[7].Description, result[10].StartDate.AddDate(0, 0, 8), result[7].StartDate)
		}
	})
}

func TestPaceTask(t *testing.T) {
	// Open our jsonFile
	jsonFile, err := os.Open("minimum_test.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	byteReaader, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
	}

	TaskResponse := TaskResponses{}
	json.Unmarshal(byteReaader, &TaskResponse)
	t.Run("Test Success Calculte multi dependance Predecessors Finish To Start", func(t *testing.T) {
		TaskResponse.Tasks[10].Start = time.Date(2024, 05, 07, 8, 0, 0, 0, time.UTC).Format("2006-01-02T15:04")
		TaskResponse.Tasks[10].Duration = "1.0d"
		result := CalculateAllProjectTask(TaskResponse.Tasks)

		if result[2].EndDate == result[2].ActualFinish {
			t.Errorf("Expected from EndDate %v but got ActualFinish %v", result[2].EndDate, result[2].ActualFinish)
		}
		if !result[2].CriticalPath {
			t.Errorf("Expected %v but got %v", true, result[2].CriticalPath)
		}
	})
}
