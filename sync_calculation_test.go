package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"
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
	// t.Run("Test Success Calculate Predecessors Finish To Start", func(t *testing.T) {
	// 	// Open our jsonFile
	// 	jsonFile, err := os.Open("minimum_test.json")
	// 	// if we os.Open returns an error then handle it
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	// defer the closing of our jsonFile so that we can parse it later on
	// 	defer jsonFile.Close()
	// 	byteReaader, err := io.ReadAll(jsonFile)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}

	// 	TaskResponse := TaskResponses{}
	// 	json.Unmarshal(byteReaader, &TaskResponse)
	// 	// fmt.Println(tasks.Count, tasks, len(tasks.Tasks))
	// 	// fmt.Print(tasks[364])
	// 	TaskResponse.Tasks[3].Start = time.Date(2022, 10, 13, 16, 0, 0, 0, time.UTC).String()
	// 	TaskResponse.Tasks[3].Finish = time.Date(2022, 10, 13, 16, 0, 0, 0, time.UTC).String()
	// 	result := CalculateAllProjectTask(TaskResponse.Tasks)
	// 	if result[16].StartDate != result[3].EndDate.AddDate(0, 0, 1) {
	// 		t.Errorf("Expected %v but got %v", result[3].EndDate.AddDate(0, 0, 1), result[16].StartDate)
	// 	}
	// 	if result[15].StartDate != result[3].EndDate.AddDate(0, 0, 1) {
	// 		t.Errorf("Expected %v but got %v", result[3].EndDate.AddDate(0, 0, 1), result[15].StartDate)
	// 	}
	// })

	t.Run("Test Success Calculte multi dependance Predecessors Finish To Start", func(t *testing.T) {

		TaskResponse := TaskResponses{}
		json.Unmarshal(byteReaader, &TaskResponse)
		// fmt.Println(tasks.Count, tasks, len(tasks.Tasks))
		// fmt.Print(tasks[364]) Predecessors Finish To Start 175, 341, 12
		// TaskResponse.Tasks[175].Start = time.Date(2022, 10, 13, 16, 0, 0, 0, time.UTC).String()
		// TaskResponse.Tasks[175].Finish = time.Date(2022, 10, 13, 16, 0, 0, 0, time.UTC).String()
		// fmt.Println(TaskResponse.Tasks[364].Start, TaskResponse.Tasks[364].Finish, "Task 364")
		// fmt.Println(TaskResponse.Tasks[175].Start, TaskResponse.Tasks[175].Finish, "Task 175")
		// fmt.Println(TaskResponse.Tasks[341].Start, TaskResponse.Tasks[341].Finish, "Task 341")
		// fmt.Println(TaskResponse.Tasks[12].Start, TaskResponse.Tasks[12].Finish, "Task 12")
		_ = CalculateAllProjectTask(TaskResponse.Tasks)
		// if result[364].StartDate != result[3].EndDate.AddDate(0, 0, 1) {
		// 	t.Errorf("Expected %v but got %v", result[3].EndDate.AddDate(0, 0, 1), result[16].StartDate)
		// }
		// if result[15].StartDate != result[3].EndDate.AddDate(0, 0, 1) {
		// 	t.Errorf("Expected %v but got %v", result[3].EndDate.AddDate(0, 0, 1), result[15].StartDate)
		// }
	})
}
