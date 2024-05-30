package main

import (
	"github.com/dominikbraun/graph"
)

type TaskCpm struct {
	Description    string
	UniqueID       string
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

func calculateCriticalComponent(tasks []TaskCpm) []TaskCpm {
	taskIndices := make(map[string]int)
	for i, task := range tasks {
		taskIndices[task.UniqueID] = i
	}
	g := graph.New(graph.StringHash, graph.Directed(), graph.Acyclic())
	for i := range tasks {
		_ = g.AddVertex(tasks[i].UniqueID)
	}
	for i := range tasks {
		for j := range tasks[i].Successor {
			_ = g.AddEdge(tasks[i].UniqueID, tasks[i].Successor[j].UniqueID)
		}
	}

	for i := range tasks {
		var currentPredID []string
		_ = graph.BFSWithDepth(g, tasks[i].UniqueID, func(value string, depth int) bool {
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

	// Step 3: Calculate Total Slack
	for i := range tasks {
		tasks[i].TotalSlack = tasks[i].LF - tasks[i].EF
		if tasks[i].TotalSlack == 0 {
			tasks[i].IsCriticalPath = true
		}
	}
	return tasks
}
