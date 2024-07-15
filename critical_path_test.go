package main

import (
	"testing"
	"time"
)

func TestCriticalPath(t *testing.T) {

	t.Run("Test Success Critical Path FS", func(t *testing.T) {
		TaskCriticalPath := []Task{
			{UniqueID: "1",
				Duration:     10,
				Successors:   []Successor{{UniqueID: "2", Lag: "0", Type: "FS"}, {UniqueID: "6", Lag: "0", Type: "FS"}},
				Predecessors: []Predecessor{},
			},
			{UniqueID: "2",
				Duration:     5,
				Successors:   []Successor{{UniqueID: "3", Lag: "0", Type: "FS"}, {UniqueID: "8", Lag: "0", Type: "FS"}},
				Predecessors: []Predecessor{{UniqueID: "1", Lag: "0", Type: "FS"}},
			},
			{UniqueID: "3",
				Duration:     7,
				Successors:   []Successor{{UniqueID: "4", Lag: "0", Type: "FS"}},
				Predecessors: []Predecessor{{UniqueID: "2", Lag: "0", Type: "FS"}},
			},
			{UniqueID: "4",
				Duration:     6,
				Successors:   []Successor{{UniqueID: "5", Lag: "0", Type: "FS"}},
				Predecessors: []Predecessor{{UniqueID: "3", Lag: "0", Type: "FS"}},
			},
			{UniqueID: "5",
				Duration:   4,
				Successors: []Successor{},
				Predecessors: []Predecessor{{UniqueID: "4", Lag: "0", Type: "FS"},
					{UniqueID: "7", Lag: "0", Type: "FS"},
					{UniqueID: "8", Lag: "0", Type: "FS"}},
			},
			{UniqueID: "6",
				Duration:     8,
				Successors:   []Successor{{UniqueID: "7", Lag: "0", Type: "FS"}, {UniqueID: "8", Lag: "0", Type: "FS"}},
				Predecessors: []Predecessor{{UniqueID: "1", Lag: "0", Type: "FS"}},
			},
			{UniqueID: "7",
				Duration:     4,
				Successors:   []Successor{{UniqueID: "5", Lag: "0", Type: "FS"}},
				Predecessors: []Predecessor{{UniqueID: "6", Lag: "0", Type: "FS"}},
			},
			{UniqueID: "8",
				Duration:     14,
				Successors:   []Successor{{UniqueID: "5", Lag: "0", Type: "FS"}},
				Predecessors: []Predecessor{{UniqueID: "2", Lag: "0", Type: "FS"}, {UniqueID: "6", Lag: "0", Type: "FS"}},
			},
		}

		// Expected Result
		expected := float64(0)

		taskCpm := calculateCriticalComponent(TaskCriticalPath)
		if taskCpm[0].TotalSlack != expected {
			t.Errorf("Expected %v but got %v", expected, taskCpm[0].TotalSlack)
		}

		if taskCpm[6].TotalSlack != float64(10) {
			t.Errorf("Expected %v but got %v", float64(10), taskCpm[6].TotalSlack)
		}
	})
}

func TestCriticalPathLag(t *testing.T) {

	t.Run("Test Success Critical Path FS", func(t *testing.T) {
		TaskCriticalPath := []Task{
			{UniqueID: "1",
				Duration:     10,
				Successors:   []Successor{{UniqueID: "2", Lag: "0", Type: "FS"}},
				Predecessors: []Predecessor{},
			},
			{UniqueID: "2",
				Duration:     5,
				Successors:   []Successor{{UniqueID: "3", Lag: "10.0d", Type: "FS"}},
				Predecessors: []Predecessor{{UniqueID: "1", Lag: "0", Type: "FS"}},
			},
			{UniqueID: "3",
				Duration:     10,
				Successors:   []Successor{{UniqueID: "4", Lag: "0", Type: "SS"}},
				Predecessors: []Predecessor{{UniqueID: "2", Lag: "10.0d", Type: "FS"}},
			},
			{UniqueID: "4",
				Duration:     1,
				Successors:   []Successor{},
				Predecessors: []Predecessor{{UniqueID: "3", Lag: "0", Type: "SS"}},
			},
			{UniqueID: "5",
				Duration:     2,
				Successors:   []Successor{{UniqueID: "6", Lag: "2.0d", Type: "SS"}},
				Predecessors: []Predecessor{},
			},
			{UniqueID: "6",
				Duration:     3,
				Successors:   []Successor{},
				Predecessors: []Predecessor{{UniqueID: "5", Lag: "2.0d", Type: "SS"}},
			},
		}

		// Expected Result
		expected := float64(0)

		taskCpm := calculateCriticalComponent(TaskCriticalPath)
		if taskCpm[0].TotalSlack != expected {
			t.Errorf("Expected %v but got %v", expected, taskCpm[0].TotalSlack)
		}

		if taskCpm[1].TotalSlack != float64(0) {
			t.Errorf("Expected %v but got %v", float64(0), taskCpm[1].TotalSlack)
		}

		// if taskCpm[2].EarlyFinish != float64(35) {
		// 	t.Errorf("Expected %v but got %v", float64(35), taskCpm[2].EarlyFinish)
		// }
		// if taskCpm[3].LateFinish != float64(35) {
		// 	t.Errorf("Expected %v but got %v", float64(35), taskCpm[3].LateFinish)
		// }
		// if taskCpm[3].LateStart != float64(34) {
		// 	t.Errorf("Expected %v but got %v", float64(34), taskCpm[3].LateStart)
		// }
		if taskCpm[3].EarlyFinish != float64(26) {
			t.Errorf("Expected %v but got %v", float64(26), taskCpm[3].EarlyFinish)
		}
		if taskCpm[3].TotalSlack != float64(9) {
			t.Errorf("Expected %v but got %v", float64(9), taskCpm[3].TotalSlack)
		}

		if taskCpm[5].TotalSlack != float64(30) {
			t.Errorf("Expected %v but got %v", float64(30), taskCpm[5].TotalSlack)
		}
		if taskCpm[5].EarlyFinish != float64(5) {
			t.Errorf("Expected %v but got %v", float64(5), taskCpm[5].EarlyFinish)
		}
		if taskCpm[5].LateFinish != float64(35) {
			t.Errorf("Expected %v but got %v", float64(35), taskCpm[5].LateFinish)
		}
	})
}

func TestCriticalPathWithHeader(t *testing.T) {
	t.Run("Test Success Critical Path with header", func(t *testing.T) {
		TaskCriticalPath := []Task{
			{UniqueID: "1",
				Duration:       2,
				DataType:       "header",
				Wbs:            "1",
				Successors:     []Successor{},
				Predecessors:   []Predecessor{},
				ChildUniqueIDs: []string{"2", "7", "10"},
			},
			{UniqueID: "2",
				Duration:       1,
				Wbs:            "1.1",
				DataType:       "header",
				Successors:     []Successor{},
				Predecessors:   []Predecessor{},
				ChildUniqueIDs: []string{"3", "4", "5"},
			},
			{UniqueID: "3",
				Duration:       1,
				Wbs:            "1.1.1",
				Successors:     []Successor{},
				Predecessors:   []Predecessor{},
				DataType:       "item",
				ChildUniqueIDs: []string{},
			},
			{UniqueID: "4",
				Duration:       1,
				Wbs:            "1.1.2",
				Successors:     []Successor{},
				DataType:       "item",
				Predecessors:   []Predecessor{},
				ChildUniqueIDs: []string{},
			},
			{UniqueID: "5",
				Duration:       1,
				Wbs:            "1.1.3",
				DataType:       "header",
				Successors:     []Successor{},
				Predecessors:   []Predecessor{},
				ChildUniqueIDs: []string{"6"},
			},
			{UniqueID: "6",
				Duration:       1,
				Wbs:            "1.1.3.1",
				DataType:       "item",
				Successors:     []Successor{},
				Predecessors:   []Predecessor{},
				ChildUniqueIDs: []string{},
			},
			{UniqueID: "7",
				Duration:       1,
				Wbs:            "1.2",
				DataType:       "header",
				Successors:     []Successor{},
				Predecessors:   []Predecessor{},
				ChildUniqueIDs: []string{"8", "9"},
			},
			{UniqueID: "8",
				Duration:       1,
				Wbs:            "1.2.1",
				DataType:       "item",
				Successors:     []Successor{},
				Predecessors:   []Predecessor{},
				ChildUniqueIDs: []string{},
			},
			{UniqueID: "9",
				Duration:       1,
				DataType:       "item",
				Wbs:            "1.2.2",
				Successors:     []Successor{},
				Predecessors:   []Predecessor{},
				ChildUniqueIDs: []string{},
			},
			{UniqueID: "10",
				Duration:       2,
				Wbs:            "1.3",
				Predecessors:   []Predecessor{},
				DataType:       "header",
				Successors:     []Successor{},
				ChildUniqueIDs: []string{"11", "12"},
			},
			{UniqueID: "11",
				Duration:       1,
				Wbs:            "1.3.1",
				DataType:       "item",
				Predecessors:   []Predecessor{},
				Successors:     []Successor{},
				ChildUniqueIDs: []string{},
			},
			{UniqueID: "12",
				Duration:       2,
				Wbs:            "1.3.2",
				DataType:       "item",
				Predecessors:   []Predecessor{},
				Successors:     []Successor{},
				ChildUniqueIDs: []string{},
			},
		}

		// Expected Result
		expected := float64(0)

		taskCpm := calculateCriticalPath(TaskCriticalPath)
		if taskCpm[0].TotalSlack != expected {
			t.Errorf("Expected %v but got %v", expected, taskCpm[0].TotalSlack)
		}
		if taskCpm[10].TotalSlack != float64(1) {
			t.Errorf("[TASK]: %v Expected %v but got %v", taskCpm[10].UniqueID, 1, taskCpm[10].TotalSlack)
		}
	})

}

func TestCriticalPathWithHeaderDependancy(t *testing.T) {

	t.Run("Test Success Critical Path Header with FS", func(t *testing.T) {
		TaskCriticalPath := []Task{
			{UniqueID: "1",
				Duration:       3,
				Wbs:            "1",
				Successors:     []Successor{},
				Predecessors:   []Predecessor{},
				ChildUniqueIDs: []string{"2", "7", "10"},
			},
			{UniqueID: "2",
				Duration:       3,
				Wbs:            "1.1",
				Successors:     []Successor{},
				Predecessors:   []Predecessor{},
				ChildUniqueIDs: []string{"3", "4", "5"},
			},
			{UniqueID: "3",
				Duration:       1,
				Wbs:            "1.1.1",
				Successors:     []Successor{},
				Predecessors:   []Predecessor{},
				ChildUniqueIDs: []string{},
			},
			{UniqueID: "4",
				Duration:       1,
				Wbs:            "1.1.2",
				Successors:     []Successor{},
				Predecessors:   []Predecessor{{UniqueID: "4", Lag: "0.0d", Type: "FS"}},
				ChildUniqueIDs: []string{},
			},
			{UniqueID: "5",
				Duration:       1,
				Wbs:            "1.1.3",
				Successors:     []Successor{},
				Predecessors:   []Predecessor{},
				ChildUniqueIDs: []string{"6"},
			},
			{UniqueID: "6",
				Duration:       1,
				Wbs:            "1.1.3.1",
				Successors:     []Successor{},
				Predecessors:   []Predecessor{},
				ChildUniqueIDs: []string{},
			},
			{UniqueID: "7",
				Duration:       1,
				Wbs:            "1.2",
				Successors:     []Successor{},
				Predecessors:   []Predecessor{},
				ChildUniqueIDs: []string{"8", "9"},
			},
			{UniqueID: "8",
				Duration:       1,
				Wbs:            "1.2.1",
				Successors:     []Successor{},
				Predecessors:   []Predecessor{},
				ChildUniqueIDs: []string{},
			},
			{UniqueID: "9",
				Duration:       1,
				Wbs:            "1.2.2",
				Successors:     []Successor{},
				Predecessors:   []Predecessor{},
				ChildUniqueIDs: []string{},
			},
			{UniqueID: "10",
				Duration:       2,
				Wbs:            "1.3",
				Predecessors:   []Predecessor{},
				Successors:     []Successor{{UniqueID: "4", Lag: "0.0d", Type: "FS"}},
				ChildUniqueIDs: []string{"11", "12"},
			},
			{UniqueID: "11",
				Duration:       1,
				Wbs:            "1.3.1",
				Predecessors:   []Predecessor{},
				Successors:     []Successor{},
				ChildUniqueIDs: []string{},
			},
			{UniqueID: "12",
				Duration:       2,
				Wbs:            "1.3.2",
				Predecessors:   []Predecessor{},
				Successors:     []Successor{},
				ChildUniqueIDs: []string{},
			},
		}

		slackTaskOne := float64(0)
		slackTaskFive := float64(2)
		slackTaskEleven := float64(1)
		taskCpm := calculateCriticalPath(TaskCriticalPath)
		if taskCpm[0].TotalSlack != slackTaskOne {
			t.Errorf("Expected %v but got %v", slackTaskOne, taskCpm[0].TotalSlack)
		}
		if taskCpm[4].TotalSlack != slackTaskFive {
			t.Errorf("[TASK]: %v Expected %v but got %v", taskCpm[4].UniqueID, slackTaskFive, taskCpm[4].TotalSlack)
		}
		if taskCpm[10].TotalSlack != slackTaskEleven {
			t.Errorf("[TASK]: %v Expected %v but got %v", taskCpm[10].UniqueID, slackTaskEleven, taskCpm[10].TotalSlack)
		}
	})
}

func TestCriticalPathWithHeaderDependancyAndChild(t *testing.T) {

	t.Run("Test Success Critical Path Header with FS", func(t *testing.T) {
		TaskCriticalPath := []Task{
			{UniqueID: "1",
				Duration:       3,
				Wbs:            "1",
				Successors:     []Successor{},
				Predecessors:   []Predecessor{},
				ChildUniqueIDs: []string{"2", "7", "10"},
			},
			{UniqueID: "2",
				Duration:       3,
				Wbs:            "1.1",
				Successors:     []Successor{},
				Predecessors:   []Predecessor{},
				ChildUniqueIDs: []string{"3", "4", "5"},
			},
			{UniqueID: "3",
				Duration:       1,
				Wbs:            "1.1.1",
				Successors:     []Successor{},
				Predecessors:   []Predecessor{},
				ChildUniqueIDs: []string{},
			},
			{UniqueID: "4",
				Duration:       1,
				Wbs:            "1.1.2",
				Successors:     []Successor{},
				Predecessors:   []Predecessor{{UniqueID: "4", Lag: "0.0d", Type: "FS"}},
				ChildUniqueIDs: []string{},
			},
			{UniqueID: "5",
				Duration:       1,
				Wbs:            "1.1.3",
				Successors:     []Successor{},
				Predecessors:   []Predecessor{},
				ChildUniqueIDs: []string{"6"},
			},
			{UniqueID: "6",
				Duration:       1,
				Wbs:            "1.1.3.1",
				Successors:     []Successor{},
				Predecessors:   []Predecessor{},
				ChildUniqueIDs: []string{},
			},
			{UniqueID: "7",
				Duration:       1,
				Wbs:            "1.2",
				Successors:     []Successor{},
				Predecessors:   []Predecessor{},
				ChildUniqueIDs: []string{"8", "9"},
			},
			{UniqueID: "8",
				Duration:       1,
				Wbs:            "1.2.1",
				Successors:     []Successor{{UniqueID: "11", Lag: "0.0d", Type: "FS"}},
				Predecessors:   []Predecessor{},
				ChildUniqueIDs: []string{},
			},
			{UniqueID: "9",
				Duration:       1,
				Wbs:            "1.2.2",
				Successors:     []Successor{},
				Predecessors:   []Predecessor{},
				ChildUniqueIDs: []string{},
			},
			{UniqueID: "10",
				Duration:       2,
				Wbs:            "1.3",
				Predecessors:   []Predecessor{},
				Successors:     []Successor{{UniqueID: "4", Lag: "0.0d", Type: "FS"}},
				ChildUniqueIDs: []string{"11", "12"},
			},
			{UniqueID: "11",
				Duration:       1,
				Wbs:            "1.3.1",
				Predecessors:   []Predecessor{{UniqueID: "8", Lag: "0.0d", Type: "FS"}},
				Successors:     []Successor{},
				ChildUniqueIDs: []string{},
			},
			{UniqueID: "12",
				Duration:       2,
				Wbs:            "1.3.2",
				Predecessors:   []Predecessor{},
				Successors:     []Successor{},
				ChildUniqueIDs: []string{},
			},
		}

		slackTaskOne := float64(0)
		slackTaskFive := float64(2)
		slackTaskEleven := float64(0)
		taskCpm := calculateCriticalPath(TaskCriticalPath)
		if taskCpm[0].TotalSlack != slackTaskOne {
			t.Errorf("Expected %v but got %v", slackTaskOne, taskCpm[0].TotalSlack)
		}
		if taskCpm[4].TotalSlack != slackTaskFive {
			t.Errorf("[TASK]: %v Expected %v but got %v", taskCpm[4].UniqueID, slackTaskFive, taskCpm[4].TotalSlack)
		}
		if taskCpm[10].TotalSlack != slackTaskEleven {
			t.Errorf("[TASK]: %v Expected %v but got %v", taskCpm[10].UniqueID, slackTaskEleven, taskCpm[10].TotalSlack)
		}
	})
}

func TestCriticalPathWithHeaderDependancyAndChildDependancy(t *testing.T) {

	t.Run("Test Success Critical Path Header with FS", func(t *testing.T) {
		TaskCriticalPath := []Task{
			{UniqueID: "1",
				Duration:       3,
				Wbs:            "1",
				Successors:     []Successor{},
				Predecessors:   []Predecessor{},
				ChildUniqueIDs: []string{"2", "7", "10"},
			},
			{UniqueID: "2",
				Duration:       3,
				Wbs:            "1.1",
				Successors:     []Successor{},
				Predecessors:   []Predecessor{},
				ChildUniqueIDs: []string{"3", "4", "5"},
			},
			{UniqueID: "3",
				Duration:       1,
				Wbs:            "1.1.1",
				Successors:     []Successor{},
				Predecessors:   []Predecessor{},
				ChildUniqueIDs: []string{},
			},
			{UniqueID: "4",
				Duration:       1,
				Wbs:            "1.1.2",
				Successors:     []Successor{},
				Predecessors:   []Predecessor{{UniqueID: "4", Lag: "0.0d", Type: "FS"}},
				ChildUniqueIDs: []string{},
			},
			{UniqueID: "5",
				Duration:       1,
				Wbs:            "1.1.3",
				Successors:     []Successor{},
				Predecessors:   []Predecessor{},
				ChildUniqueIDs: []string{"6"},
			},
			{UniqueID: "6",
				Duration:       1,
				Wbs:            "1.1.3.1",
				Successors:     []Successor{},
				Predecessors:   []Predecessor{},
				ChildUniqueIDs: []string{},
			},
			{UniqueID: "7",
				Duration:       1,
				Wbs:            "1.2",
				Successors:     []Successor{},
				Predecessors:   []Predecessor{},
				ChildUniqueIDs: []string{"8", "9"},
			},
			{UniqueID: "8",
				Duration:       1,
				Wbs:            "1.2.1",
				Successors:     []Successor{},
				Predecessors:   []Predecessor{{UniqueID: "12", Lag: "0.0d", Type: "FS"}},
				ChildUniqueIDs: []string{},
			},
			{UniqueID: "9",
				Duration:       1,
				Wbs:            "1.2.2",
				Successors:     []Successor{},
				Predecessors:   []Predecessor{},
				ChildUniqueIDs: []string{},
			},
			{UniqueID: "10",
				Duration:       2,
				Wbs:            "1.3",
				Predecessors:   []Predecessor{},
				Successors:     []Successor{{UniqueID: "4", Lag: "0.0d", Type: "FS"}},
				ChildUniqueIDs: []string{"11", "12"},
			},
			{UniqueID: "11",
				Duration:       1,
				Wbs:            "1.3.1",
				Predecessors:   []Predecessor{},
				Successors:     []Successor{},
				ChildUniqueIDs: []string{},
			},
			{UniqueID: "12",
				Duration:       2,
				Wbs:            "1.3.2",
				Predecessors:   []Predecessor{},
				Successors:     []Successor{{UniqueID: "8", Lag: "0.0d", Type: "FS"}},
				ChildUniqueIDs: []string{},
			},
		}

		slackTaskOne := float64(0)
		slackTaskEight := float64(0)
		slackTaskTwelve := float64(0)
		taskCpm := calculateCriticalPath(TaskCriticalPath)
		if taskCpm[0].TotalSlack != slackTaskOne {
			t.Errorf("Expected %v but got %v", slackTaskOne, taskCpm[0].TotalSlack)
		}
		if taskCpm[7].TotalSlack != slackTaskEight {
			t.Errorf("[TASK]: %v Expected %v but got %v", taskCpm[4].UniqueID, slackTaskEight, taskCpm[4].TotalSlack)
		}
		if taskCpm[11].TotalSlack != slackTaskTwelve {
			t.Errorf("[TASK]: %v Expected %v but got %v", taskCpm[10].UniqueID, slackTaskTwelve, taskCpm[10].TotalSlack)
		}
	})
}

func TestCriticalPathWithHeaderDependancyFSFFSS(t *testing.T) {

	t.Run("Test Success Critical Path Header with FS", func(t *testing.T) {
		TaskCriticalPath := []Task{
			{UniqueID: "1",
				Duration:       4,
				Wbs:            "1",
				Successors:     []Successor{},
				Predecessors:   []Predecessor{},
				ChildUniqueIDs: []string{"2", "7", "10"},
			},
			{UniqueID: "2",
				Duration:       3,
				Wbs:            "1.1",
				Successors:     []Successor{},
				Predecessors:   []Predecessor{},
				ChildUniqueIDs: []string{"3", "4", "5"},
			},
			{UniqueID: "3",
				Duration:       1,
				Wbs:            "1.1.1",
				Successors:     []Successor{},
				Predecessors:   []Predecessor{},
				ChildUniqueIDs: []string{},
			},
			{UniqueID: "4",
				Duration:       1,
				Wbs:            "1.1.2",
				Successors:     []Successor{},
				Predecessors:   []Predecessor{{UniqueID: "10", Lag: "0.0d", Type: "FS"}},
				ChildUniqueIDs: []string{},
			},
			{UniqueID: "5",
				Duration:       1,
				Wbs:            "1.1.3",
				Successors:     []Successor{},
				Predecessors:   []Predecessor{},
				ChildUniqueIDs: []string{"6"},
			},
			{UniqueID: "6",
				Duration:       1,
				Wbs:            "1.1.3.1",
				Successors:     []Successor{},
				Predecessors:   []Predecessor{},
				ChildUniqueIDs: []string{},
			},
			{UniqueID: "7",
				Duration:       3,
				Wbs:            "1.2",
				Successors:     []Successor{},
				Predecessors:   []Predecessor{},
				ChildUniqueIDs: []string{"8", "9", "13"},
			},
			{UniqueID: "8",
				Duration:       1,
				Wbs:            "1.2.1",
				Successors:     []Successor{},
				Predecessors:   []Predecessor{{UniqueID: "12", Lag: "0.0d", Type: "FS"}},
				ChildUniqueIDs: []string{},
			},
			{UniqueID: "9",
				Duration:       2,
				Wbs:            "1.2.2",
				Successors:     []Successor{},
				Predecessors:   []Predecessor{},
				ChildUniqueIDs: []string{},
			},
			{UniqueID: "10",
				Duration:       2,
				Wbs:            "1.3",
				Predecessors:   []Predecessor{},
				Successors:     []Successor{{UniqueID: "4", Lag: "0.0d", Type: "FS"}},
				ChildUniqueIDs: []string{"11", "12"},
			},
			{UniqueID: "11",
				Duration:       1,
				Wbs:            "1.3.1",
				Predecessors:   []Predecessor{},
				Successors:     []Successor{},
				ChildUniqueIDs: []string{},
			},
			{UniqueID: "12",
				Duration:       2,
				Wbs:            "1.3.2",
				Predecessors:   []Predecessor{},
				Successors:     []Successor{{UniqueID: "8", Lag: "0.0d", Type: "FS"}},
				ChildUniqueIDs: []string{},
			},
			{UniqueID: "13",
				Duration:       1,
				Wbs:            "1.3.2",
				Predecessors:   []Predecessor{},
				Successors:     []Successor{},
				ChildUniqueIDs: []string{},
			},
			{UniqueID: "14",
				Duration:       4,
				Wbs:            "1.3",
				Predecessors:   []Predecessor{},
				Successors:     []Successor{},
				ChildUniqueIDs: []string{"15", "17", "18", "19", "20"},
			},
			{UniqueID: "15",
				Duration:       1,
				Wbs:            "1.3",
				Predecessors:   []Predecessor{},
				Successors:     []Successor{},
				ChildUniqueIDs: []string{"16"},
			},
			{UniqueID: "16",
				Duration:       1,
				Wbs:            "1.3",
				Predecessors:   []Predecessor{{UniqueID: "18", Lag: "0.0d", Type: "FF"}},
				Successors:     []Successor{},
				ChildUniqueIDs: []string{},
			},
			{UniqueID: "17",
				Duration:       1,
				Wbs:            "1.3",
				Successors:     []Successor{},
				ChildUniqueIDs: []string{},
			},
			{UniqueID: "18",
				Duration:     1,
				Wbs:          "1.3",
				Predecessors: []Predecessor{},
				Successors: []Successor{{UniqueID: "16", Lag: "0.0d", Type: "FF"},
					{UniqueID: "19", Lag: "0.0d", Type: "FS"},
					{UniqueID: "22", Lag: "0.0d", Type: "SS"}},
				ChildUniqueIDs: []string{},
			},
			{UniqueID: "19",
				Duration:       1,
				Wbs:            "1.3",
				Predecessors:   []Predecessor{{UniqueID: "18", Lag: "0.0d", Type: "FS"}},
				Successors:     []Successor{{UniqueID: "20", Lag: "0.0d", Type: "FS"}},
				ChildUniqueIDs: []string{},
			},
			{UniqueID: "20",
				Duration:       1,
				Wbs:            "1.3",
				Predecessors:   []Predecessor{{UniqueID: "19", Lag: "0.0d", Type: "FS"}},
				Successors:     []Successor{{UniqueID: "21", Lag: "0.0d", Type: "FS"}},
				ChildUniqueIDs: []string{},
			},
			{UniqueID: "21",
				Duration:       1,
				Wbs:            "1.3",
				Predecessors:   []Predecessor{{UniqueID: "20", Lag: "0.0d", Type: "FS"}},
				Successors:     []Successor{},
				ChildUniqueIDs: []string{},
			},
			{UniqueID: "22",
				Duration:       1,
				Wbs:            "1.3",
				Predecessors:   []Predecessor{{UniqueID: "18", Lag: "0.0d", Type: "SS"}},
				Successors:     []Successor{},
				ChildUniqueIDs: []string{},
			},
		}

		expected := []float64{0, 1, 3, 1, 3, 3, 1, 1, 2, 1, 2, 1, 3, 0, 3, 3, 3, 0, 0, 0, 0, 3}
		taskCpm := calculateCriticalPath(TaskCriticalPath)
		for id := range taskCpm {
			if taskCpm[id].TotalSlack != expected[id] {
				t.Errorf("Task Unique ID: %v Expected %v but got %v", taskCpm[id].UniqueID, expected[id], taskCpm[id].TotalSlack)
			}
		}
	})
}

func TestCriticalPathWithCustomCalendar(t *testing.T) {

	t.Run("Test Success Critical Path ", func(t *testing.T) {
		TaskCriticalPath := []Task{
			{UniqueID: "1",
				Duration:       33,
				Successors:     []Successor{},
				Predecessors:   []Predecessor{},
				ChildUniqueIDs: []string{"2", "5"},
				StartDate:      time.Date(2024, 6, 11, 0, 0, 0, 0, time.UTC),
			},
			{UniqueID: "2",
				Duration:       31,
				Successors:     []Successor{},
				Predecessors:   []Predecessor{},
				ChildUniqueIDs: []string{"3", "4"},
				StartDate:      time.Date(2024, 6, 11, 0, 0, 0, 0, time.UTC),
			},
			{UniqueID: "3",
				Duration:     29,
				Successors:   []Successor{},
				Predecessors: []Predecessor{},
				StartDate:    time.Date(2024, 6, 12, 0, 0, 0, 0, time.UTC),
			},
			{UniqueID: "4",
				Duration:     31,
				Successors:   []Successor{},
				Predecessors: []Predecessor{},
				StartDate:    time.Date(2024, 6, 11, 0, 0, 0, 0, time.UTC),
			},
			{UniqueID: "5",
				Duration:       33,
				Successors:     []Successor{},
				Predecessors:   []Predecessor{},
				ChildUniqueIDs: []string{"6", "7"},
				StartDate:      time.Date(2024, 6, 12, 0, 0, 0, 0, time.UTC),
			},
			{UniqueID: "6",
				Duration:     32,
				Successors:   []Successor{},
				Predecessors: []Predecessor{},
				StartDate:    time.Date(2024, 6, 12, 0, 0, 0, 0, time.UTC),
			},
			{UniqueID: "7",
				Duration:     32,
				Successors:   []Successor{},
				Predecessors: []Predecessor{},
				StartDate:    time.Date(2024, 6, 11, 0, 0, 0, 0, time.UTC),
			},
		}

		// Expected Result
		expected := []float64{0, 2, 4, 2, 0, 1, 1}
		taskCpm := calculateCriticalPath(TaskCriticalPath)
		for id := range taskCpm {
			if taskCpm[id].TotalSlack != expected[id] {
				t.Errorf("Task Unique ID: %v Expected %v but got %v", taskCpm[id].UniqueID, expected[id], taskCpm[id].TotalSlack)
			}
		}

	})
}
