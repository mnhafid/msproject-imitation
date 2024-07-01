package main

import (
	"testing"
	"time"
)

func TestEtc(t *testing.T) {

	t.Run("Test Success ETC Path FS", func(t *testing.T) {
		TaskCriticalPath := []Task{
			{UniqueID: "1",
				Duration:       10,
				Successors:     []Successor{{UniqueID: "2", Lag: "0", Type: "FS"}, {UniqueID: "6", Lag: "0", Type: "FS"}},
				Predecessors:   []Predecessor{},
				ActualProgress: 5,
				StartDate:      time.Date(2024, 6, 23, 0, 0, 0, 0, time.UTC),
				EndDate:        time.Date(2024, 8, 3, 8, 0, 0, 0, time.UTC),
				ActualStart:    time.Date(2024, 6, 26, 8, 0, 0, 0, time.UTC),
				ActualFinish:   time.Date(2024, 6, 26, 8, 0, 0, 0, time.UTC),
				StartEtc:       time.Date(2024, 6, 23, 8, 0, 0, 0, time.UTC),
				Etc:            time.Date(2024, 8, 3, 8, 0, 0, 0, time.UTC),
			},
			{UniqueID: "2",
				Duration:     5,
				Successors:   []Successor{{UniqueID: "3", Lag: "0", Type: "FS"}, {UniqueID: "8", Lag: "0", Type: "FS"}},
				Predecessors: []Predecessor{{UniqueID: "1", Lag: "0", Type: "FS"}},
				StartDate:    time.Date(2024, 8, 6, 0, 0, 0, 0, time.UTC),
				EndDate:      time.Date(2024, 9, 3, 8, 0, 0, 0, time.UTC),
				StartEtc:     time.Date(2024, 8, 6, 0, 0, 0, 0, time.UTC),
				Etc:          time.Date(2024, 9, 3, 8, 0, 0, 0, time.UTC),
			},
			{UniqueID: "3",
				Duration:     5,
				Successors:   []Successor{{UniqueID: "3", Lag: "0", Type: "FS"}, {UniqueID: "8", Lag: "0", Type: "FS"}},
				Predecessors: []Predecessor{{UniqueID: "1", Lag: "0", Type: "FS"}},
				StartDate:    time.Date(2024, 9, 6, 0, 0, 0, 0, time.UTC),
				EndDate:      time.Date(2024, 9, 30, 8, 0, 0, 0, time.UTC),
				StartEtc:     time.Date(2024, 9, 6, 0, 0, 0, 0, time.UTC),
				Etc:          time.Date(2024, 9, 30, 8, 0, 0, 0, time.UTC),
			},
		}

		taskIndices := make(map[string]int)
		for i := 0; i < len(TaskCriticalPath); i++ {
			taskIndices[TaskCriticalPath[i].UniqueID] = i
		}
		result := []time.Time{time.Date(2024, 6, 23, 8, 0, 0, 0, time.UTC),
			time.Date(2024, 7, 25, 8, 0, 0, 0, time.UTC),
			time.Date(2024, 8, 25, 8, 0, 0, 0, time.UTC)}

		TaskCriticalPath = PrepareEtc(TaskCriticalPath, taskIndices)
		for i := range TaskCriticalPath {
			if TaskCriticalPath[i].Etc.Truncate(24*time.Hour) != result[i].Truncate(24*time.Hour) {
				t.Errorf("[TASK]: %v Expected %v but got %v", TaskCriticalPath[i].UniqueID, result[i], TaskCriticalPath[i].Etc)
			}
		}

	})
}
