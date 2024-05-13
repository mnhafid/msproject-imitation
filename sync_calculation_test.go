package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"
	"time"
)

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
	t.Run("Test Success Calculte multi dependance Predecessors Finish To Start", func(t *testing.T) {
		TaskResponse.Tasks[10].Finish = time.Date(2024, 05, 06, 17, 0, 0, 0, time.UTC).Format("2006-01-02T15:04")
		TaskResponse.Tasks[10].Duration = "1.0d"
		result := CalculateAllProjectTask(TaskResponse.Tasks)

		if result[2].ActualStart != result[10].StartDate.AddDate(0, 0, 6) {
			t.Errorf("Expected from StartDate %v but got ActualStart %v", result[10].StartDate.AddDate(0, 0, 6), result[2].ActualStart)
		}
		if result[3].ActualStart != result[10].StartDate.AddDate(0, 0, 7) {
			t.Errorf("Expected %v but got %v", result[10].StartDate.AddDate(0, 0, 7), result[3].ActualStart)
		}
		if result[4].ActualStart != result[10].StartDate.AddDate(0, 0, 7) {
			t.Errorf("Expected %v but got %v", result[10].StartDate.AddDate(0, 0, 7), result[4].ActualStart)
		}
		if result[5].ActualStart != result[10].StartDate.AddDate(0, 0, 1) {
			t.Errorf("Expected from EndDate %v but got PlanStart %v", result[10].StartDate.AddDate(0, 0, 1), result[5].ActualStart)
		}
		if result[6].ActualStart != result[8].ActualStart {
			t.Errorf("Expected %v but got %v", result[8].ActualStart, result[6].ActualStart)
		}
		if result[7].ActualStart != result[10].StartDate.AddDate(0, 0, 8) {
			t.Errorf("Expected %v but got %v", result[10].StartDate.AddDate(0, 0, 8), result[7].ActualStart)
		}
		if result[8].EndDate != result[10].EndDate {
			t.Errorf("Expected %v but got %v", result[10].EndDate, result[8].ActualFinish)
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
		TaskResponse.Tasks[10].Finish = time.Date(2024, 05, 06, 17, 0, 0, 0, time.UTC).Format("2006-01-02T15:04")
		TaskResponse.Tasks[10].Duration = "1.0d"
		result := CalculateAllProjectTask(TaskResponse.Tasks)
		if result[5].ActualStart != result[10].EndDate.AddDate(0, 0, 1) {
			t.Errorf("Expected from EndDate %v but got PlanStart %v", result[10].EndDate.AddDate(0, 0, 1), result[5].ActualStart)
		}
		if result[3].ActualStart != result[10].EndDate.AddDate(0, 0, 7) {
			t.Errorf("Expected %v but got %v", result[10].EndDate.AddDate(0, 0, 7), result[3].ActualStart)
		}
	})
}
