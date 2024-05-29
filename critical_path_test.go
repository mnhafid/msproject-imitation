package main

import (
	"testing"
)

func TestCriticalPath(t *testing.T) {

	t.Run("Test Success Critical Path FS", func(t *testing.T) {
		TaskCriticalPath := []TaskCpm{
			{ID: "1",
				Duration:     10,
				Successor:    []Successor{{ID: "2", Lag: "0"}, {ID: "6", Lag: "0"}},
				Predecessors: []Predecessor{},
			},
			{ID: "2",
				Duration:     5,
				Successor:    []Successor{{ID: "3", Lag: "0"}, {ID: "8", Lag: "0"}},
				Predecessors: []Predecessor{{ID: "1", Lag: "0"}},
			},
			{ID: "3",
				Duration:     7,
				Successor:    []Successor{{ID: "4", Lag: "0"}},
				Predecessors: []Predecessor{{ID: "2", Lag: "0"}},
			},
			{ID: "4",
				Duration:     6,
				Successor:    []Successor{{ID: "5", Lag: "0"}},
				Predecessors: []Predecessor{{ID: "3", Lag: "0"}},
			},
			{ID: "5",
				Duration:  4,
				Successor: []Successor{},
				Predecessors: []Predecessor{{ID: "4", Lag: "0"},
					{ID: "7", Lag: "0"},
					{ID: "8", Lag: "0"}},
			},
			{ID: "6",
				Duration:     8,
				Successor:    []Successor{{ID: "7", Lag: "0"}, {ID: "8", Lag: "0"}},
				Predecessors: []Predecessor{{ID: "1", Lag: "0"}},
			},
			{ID: "7",
				Duration:     4,
				Successor:    []Successor{{ID: "5", Lag: "0"}},
				Predecessors: []Predecessor{{ID: "6", Lag: "0"}},
			},
			{ID: "8",
				Duration:     14,
				Successor:    []Successor{{ID: "5", Lag: "0"}},
				Predecessors: []Predecessor{{ID: "2", Lag: "0"}, {ID: "6", Lag: "0"}},
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
