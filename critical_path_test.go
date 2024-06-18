package main

import (
	"testing"
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

func TestCriticalPathMultiDependency(t *testing.T) {

	t.Run("Test Success Critical Path Multi Dependency", func(t *testing.T) {
		TaskCriticalPath := []Task{
			{UniqueID: "1",
				Duration:     1,
				Successors:   []Successor{{UniqueID: "2", Lag: "2.d", Type: "SS"}},
				Predecessors: []Predecessor{},
			},
			{UniqueID: "2",
				Duration:     1,
				Successors:   []Successor{{UniqueID: "3", Lag: "14.0d", Type: "FS"}},
				Predecessors: []Predecessor{{UniqueID: "1", Lag: "2.0d", Type: "SS"}},
			},
			{UniqueID: "3",
				Duration:     1,
				Successors:   []Successor{},
				Predecessors: []Predecessor{{UniqueID: "2", Lag: "14.0d", Type: "FS"}},
			},
			{UniqueID: "4",
				Duration:     1,
				Successors:   []Successor{{UniqueID: "6", Lag: "0", Type: "FS"}},
				Predecessors: []Predecessor{{UniqueID: "5", Lag: "0", Type: "FF"}},
			},
			{UniqueID: "5",
				Duration:     1,
				Successors:   []Successor{{UniqueID: "4", Lag: "0", Type: "FF"}},
				Predecessors: []Predecessor{},
			},
			{UniqueID: "6",
				Duration:     1,
				Successors:   []Successor{},
				Predecessors: []Predecessor{{UniqueID: "4", Lag: "0", Type: "FS"}},
			},
			{UniqueID: "7",
				Duration:     1,
				Successors:   []Successor{{UniqueID: "8", Lag: "0", Type: "SF"}},
				Predecessors: []Predecessor{},
			},
			{UniqueID: "8",
				Duration:     1,
				Successors:   []Successor{},
				Predecessors: []Predecessor{{UniqueID: "7", Lag: "0", Type: "SF"}},
			},
			{UniqueID: "9",
				Duration:     1,
				Successors:   []Successor{{UniqueID: "10", Lag: "10.0d", Type: "SF"}},
				Predecessors: []Predecessor{},
			},
			{UniqueID: "10",
				Duration:     1,
				Successors:   []Successor{},
				Predecessors: []Predecessor{{UniqueID: "9", Lag: "10.0d", Type: "SF"}},
			},
		}

		// Expected Result
		expected := float64(0)

		taskCpm := calculateCriticalComponent(TaskCriticalPath)
		if taskCpm[0].TotalSlack != expected {
			t.Errorf("Expected %v but got %v", expected, taskCpm[0].TotalSlack)
		}

		if taskCpm[5].TotalSlack != float64(16) {
			t.Errorf("Expected %v but got %v", float64(16), taskCpm[5].TotalSlack)
		}

		if taskCpm[6].EarlyStart != float64(0) {
			t.Errorf("Expected %v but got %v", float64(0), taskCpm[6].EarlyStart)
		}
		if taskCpm[6].EarlyFinish != float64(1) {
			t.Errorf("Expected %v but got %v", float64(1), taskCpm[6].EarlyFinish)
		}
		if taskCpm[6].LateStart != float64(17) {
			t.Errorf("Expected %v but got %v Late Start", float64(16), taskCpm[6].LateStart)
		}
		if taskCpm[6].LateFinish != float64(18) {
			t.Errorf("Expected %v but got %v Late Finish", float64(17), taskCpm[6].LateFinish)
		}
		if taskCpm[6].TotalSlack != float64(17) {
			t.Errorf("Expected %v but got %v slack task 6", float64(17), taskCpm[6].TotalSlack)
		}

		if taskCpm[7].EarlyStart != float64(-1) {
			t.Errorf("Expected %v but got %v Early start task 8", float64(-1), taskCpm[7].EarlyStart)
		}
		if taskCpm[7].EarlyFinish != float64(0) {
			t.Errorf("Expected %v but got %v early finish task 8", float64(0), taskCpm[7].EarlyFinish)
		}
		if taskCpm[7].LateStart != float64(17) {
			t.Errorf("Expected %v but got %v Late Start", float64(18), taskCpm[7].LateStart)
		}
		if taskCpm[7].LateFinish != float64(18) {
			t.Errorf("Expected %v but got %v Late Finish", float64(18), taskCpm[7].LateFinish)
		}
		if taskCpm[7].TotalSlack != float64(18) {
			t.Errorf("Expected %v but got %v", float64(18), taskCpm[7].TotalSlack)
		}
		if taskCpm[9].TotalSlack != float64(8) {
			t.Errorf("Expected %v but got %v", float64(8), taskCpm[9].TotalSlack)
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

	t.Run("Test Success Critical Path HEader with FS", func(t *testing.T) {
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

		// Expected Result
		expected := float64(0)

		taskCpm := calculateCriticalPath(TaskCriticalPath)
		if taskCpm[0].TotalSlack != expected {
			t.Errorf("Expected %v but got %v", expected, taskCpm[0].TotalSlack)
		}
		if taskCpm[6].TotalSlack != float64(2) {
			t.Errorf("[TASK]: %v Expected %v but got %v", taskCpm[6].UniqueID, 2, taskCpm[6].TotalSlack)
		}
		if taskCpm[10].TotalSlack != float64(1) {
			t.Errorf("[TASK]: %v Expected %v but got %v", taskCpm[10].UniqueID, 1, taskCpm[10].TotalSlack)
		}
	})
}
