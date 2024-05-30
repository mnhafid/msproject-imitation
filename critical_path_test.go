package main

import (
	"testing"
)

func TestCriticalPath(t *testing.T) {

	t.Run("Test Success Critical Path FS", func(t *testing.T) {
		TaskCriticalPath := []Task{
			{UniqueID: "1",
				Duration:     10,
				Successors:   []Successor{{UniqueID: "2", Lag: "0"}, {UniqueID: "6", Lag: "0"}},
				Predecessors: []Predecessor{},
			},
			{UniqueID: "2",
				Duration:     5,
				Successors:   []Successor{{UniqueID: "3", Lag: "0"}, {UniqueID: "8", Lag: "0"}},
				Predecessors: []Predecessor{{UniqueID: "1", Lag: "0"}},
			},
			{UniqueID: "3",
				Duration:     7,
				Successors:   []Successor{{UniqueID: "4", Lag: "0"}},
				Predecessors: []Predecessor{{UniqueID: "2", Lag: "0"}},
			},
			{UniqueID: "4",
				Duration:     6,
				Successors:   []Successor{{UniqueID: "5", Lag: "0"}},
				Predecessors: []Predecessor{{UniqueID: "3", Lag: "0"}},
			},
			{UniqueID: "5",
				Duration:   4,
				Successors: []Successor{},
				Predecessors: []Predecessor{{UniqueID: "4", Lag: "0"},
					{UniqueID: "7", Lag: "0"},
					{UniqueID: "8", Lag: "0"}},
			},
			{UniqueID: "6",
				Duration:     8,
				Successors:   []Successor{{UniqueID: "7", Lag: "0"}, {UniqueID: "8", Lag: "0"}},
				Predecessors: []Predecessor{{UniqueID: "1", Lag: "0"}},
			},
			{UniqueID: "7",
				Duration:     4,
				Successors:   []Successor{{UniqueID: "5", Lag: "0"}},
				Predecessors: []Predecessor{{UniqueID: "6", Lag: "0"}},
			},
			{UniqueID: "8",
				Duration:     14,
				Successors:   []Successor{{UniqueID: "5", Lag: "0"}},
				Predecessors: []Predecessor{{UniqueID: "2", Lag: "0"}, {UniqueID: "6", Lag: "0"}},
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

func TestCriticalPathMultiDependency(t *testing.T) {

	t.Run("Test Success Critical Path Multi Dependency", func(t *testing.T) {
		TaskCriticalPath := []Task{
			{UniqueID: "1",
				Duration:     2,
				Successors:   []Successor{{UniqueID: "2", Lag: "0"}},
				Predecessors: []Predecessor{},
			},
			{UniqueID: "2",
				Duration:     10,
				Successors:   []Successor{{UniqueID: "3", Lag: "0"}, {UniqueID: "7", Lag: "0"}},
				Predecessors: []Predecessor{{UniqueID: "1", Lag: "0"}},
			},
			{UniqueID: "3",
				Duration:     14,
				Successors:   []Successor{{UniqueID: "4", Lag: "0"}, {UniqueID: "6", Lag: "0"}, {UniqueID: "7", Lag: "0"}},
				Predecessors: []Predecessor{{UniqueID: "2", Lag: "0"}},
			},
			{UniqueID: "4",
				Duration:     3,
				Successors:   []Successor{{UniqueID: "5", Lag: "0"}, {UniqueID: "6", Lag: "0"}, {UniqueID: "7", Lag: "0"}},
				Predecessors: []Predecessor{{UniqueID: "4", Lag: "0"}},
			},
			{UniqueID: "5",
				Duration:     4,
				Successors:   []Successor{},
				Predecessors: []Predecessor{{UniqueID: "4", Lag: "0"}},
			},
			{UniqueID: "6",
				Duration:     5,
				Successors:   []Successor{},
				Predecessors: []Predecessor{{UniqueID: "3", Lag: "0"}, {UniqueID: "4", Lag: "0"}},
			},
			{UniqueID: "7",
				Duration:     7,
				Successors:   []Successor{},
				Predecessors: []Predecessor{{UniqueID: "2", Lag: "0"}, {UniqueID: "3", Lag: "0"}, {UniqueID: "4", Lag: "0"}},
			},
		}

		// Expected Result
		expected := float64(2)

		taskCpm := calculateCriticalComponent(TaskCriticalPath)
		if taskCpm[0].TotalSlack != expected {
			t.Errorf("Expected %v but got %v", expected, taskCpm[0].TotalSlack)
		}

		if taskCpm[6].TotalSlack != float64(0) {
			t.Errorf("Expected %v but got %v", float64(0), taskCpm[6].TotalSlack)
		}
	})
}
