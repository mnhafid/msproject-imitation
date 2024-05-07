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
	t.Run("Test Success Calculate Predecessors Finish To Start", func(t *testing.T) {
		// Open our jsonFile
		testFile, err := os.Open("test_file_js.json")
		// if we os.Open returns an error then handle it
		if err != nil {
			fmt.Println(err)
		}
		// defer the closing of our jsonFile so that we can parse it later on
		defer testFile.Close()
		byteReader, err := io.ReadAll(testFile)
		if err != nil {
			fmt.Println(err)
		}

		TaskResponses := TaskResponses{}
		json.Unmarshal(byteReader, &TaskResponses)
		// fmt.Println(tasks.Count, tasks, len(tasks.Tasks))
		// fmt.Print(tasks[364])
		TaskResponses.Tasks[3].Start = time.Date(2022, 10, 13, 16, 0, 0, 0, time.UTC).String()
		TaskResponses.Tasks[3].Finish = time.Date(2022, 10, 13, 16, 0, 0, 0, time.UTC).String()
		result := CalculateAllProjectTask(TaskResponses.Tasks)
		if result[16].StartDate != result[3].EndDate.AddDate(0, 0, 1) {
			t.Errorf("Expected %v but got %v", result[3].EndDate.AddDate(0, 0, 1), result[16].StartDate)
		}
		if result[15].StartDate != result[3].EndDate.AddDate(0, 0, 1) {
			t.Errorf("Expected %v but got %v", result[3].EndDate.AddDate(0, 0, 1), result[15].StartDate)
		}
	})

	t.Run("Test Success Calculte multi dependance Predecessors Finish To Start", func(t *testing.T) {

		result := CalculateAllProjectTask(TaskResponse.Tasks)
		if result[2].PlanStart != result[10].EndDate.AddDate(0, 0, 1) {
			t.Errorf("Expected %v but got %v", result[2].EndDate.AddDate(0, 0, 1), result[10].StartDate)
		}
		if result[3].StartDate != result[10].EndDate.AddDate(0, 0, 1) {
			t.Errorf("Expected %v but got %v", result[3].EndDate.AddDate(0, 0, 1), result[10].StartDate)
		}
	})
}
