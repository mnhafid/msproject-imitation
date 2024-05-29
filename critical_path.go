package main

import (
	"fmt"
	"os"
	"time"

	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
)

type TaskCpm struct {
	ID             string
	Description    string
	UniqueID       string
	StartDate      time.Time
	EndDate        time.Time
	ES             float64
	EF             float64
	LS             float64
	LF             float64
	TotalSlack     float64
	Duration       float64
	Predecessors   []Predecessor
	Successor      []Successor
	IsCriticalPath bool
}

type Dependency struct {
	PredecessorID string
	SuccessorID   string
	Type          string
}

func calculateCriticalComponent(tasks []TaskCpm) []TaskCpm {
	taskIndices := make(map[string]int)
	for i, task := range tasks {
		taskIndices[task.ID] = i
	}
	g := graph.New(graph.StringHash, graph.Directed(), graph.Acyclic())
	for i := range tasks {
		_ = g.AddVertex(tasks[i].ID)
	}
	for i := range tasks {
		for j := range tasks[i].Successor {
			_ = g.AddEdge(tasks[i].ID, tasks[i].Successor[j].ID)
		}
	}
	file, _ := os.Create("./mygraph.gv")
	_ = draw.DOT(g, file)
	for i := range tasks {
		var currentPredID []string
		_ = graph.BFSWithDepth(g, tasks[i].ID, func(value string, depth int) bool {
			currentPredID = append(currentPredID, value)
			return depth > 2
		})
		// Step 1: Forward Pass
		for j := range currentPredID {
			// define ES and EF for the first task
			tasks[taskIndices[currentPredID[j]]].ES = 0
			tasks[taskIndices[currentPredID[j]]].EF = tasks[taskIndices[currentPredID[j]]].Duration
			if len(tasks[taskIndices[currentPredID[j]]].Predecessors) > 0 {
				for _, predecessor := range tasks[taskIndices[currentPredID[j]]].Predecessors {
					tasks[taskIndices[currentPredID[j]]].ES = tasks[taskIndices[predecessor.ID]].EF
					tasks[taskIndices[currentPredID[j]]].EF = tasks[taskIndices[currentPredID[j]]].ES + tasks[taskIndices[currentPredID[j]]].Duration
				}
			}
		}
	}
	// Step 2: Backward Pass
	// LS = LF - duration
	// LF = min(ES of all successors)
	for i := range tasks {
		tasks[i].LF = tasks[i].EF
		if len(tasks[i].Successor) > 0 {
			for _, successor := range tasks[i].Successor {
				tasks[i].LF = tasks[taskIndices[successor.ID]].ES
			}
		}
		tasks[i].LS = tasks[i].LF - tasks[i].Duration
	}
	for i := range tasks {
		tasks[i].TotalSlack = tasks[i].LF - tasks[i].EF
		if tasks[i].TotalSlack == 0 {
			tasks[i].IsCriticalPath = true
		}
	}
	fmt.Println("[RESULT]", "Task: ", tasks[0].ID, "ES: ", tasks[0].ES, "EF: ", tasks[0].EF, "Duration: ", tasks[0].Duration)
	fmt.Println("[RESULT]", "Task: ", tasks[4].ID, "ES: ", tasks[4].ES, "EF: ", tasks[4].EF, "Duration: ", tasks[4].Duration)
	fmt.Println("[RESULT]", "Task: ", tasks[6].ID, "ES: ", tasks[6].ES, "EF: ", tasks[6].EF, "Duration: ", tasks[6].Duration)
	fmt.Println("[RESULT]", "Task: ", tasks[7].ID, "ES: ", tasks[7].ES, "EF: ", tasks[7].EF, "Duration: ", tasks[7].Duration)

	fmt.Println("[RESULT]", "Task: ", tasks[0].ID, "LS: ", tasks[0].LS, "LF: ", tasks[0].LF, "Duration: ", tasks[0].Duration)
	fmt.Println("[RESULT]", "Task: ", tasks[4].ID, "LS: ", tasks[4].LS, "LF: ", tasks[4].LF, "Duration: ", tasks[4].Duration)
	fmt.Println("[RESULT]", "Task: ", tasks[6].ID, "LS: ", tasks[6].LS, "LF: ", tasks[6].LF, "Duration: ", tasks[6].Duration)
	fmt.Println("[RESULT]", "Task: ", tasks[7].ID, "LS: ", tasks[7].LS, "LF: ", tasks[7].LF, "Duration: ", tasks[7].Duration)
	return tasks
}
