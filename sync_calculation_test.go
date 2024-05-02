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
	t.Run("Test Success Set Predecessors Finish To Start", func(t *testing.T) {
		// Open our jsonFile
		jsonFile, err := os.Open("test_file_fs.json")
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
		// fmt.Println(tasks.Count, tasks, len(tasks.Tasks))
		// fmt.Print(tasks[364])
		result := CalculateAllProjectTask(TaskResponse.Tasks)
		if result[16].StartDate != result[3].EndDate.AddDate(0, 0, 3) {
			t.Errorf("Expected %v but got %v", result[3].EndDate.AddDate(0, 0, 3), result[16].StartDate)
		}
	})
}
