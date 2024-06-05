package main

import (
	"fmt"
	"os"

	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
)

func calculateCriticalComponent(tasks []Task) []Task {
	taskIndices := make(map[string]int)
	for i, task := range tasks {
		taskIndices[task.UniqueID] = i
	}
	g := graph.New(graph.StringHash, graph.Directed(), graph.Acyclic())
	for i := range tasks {
		_ = g.AddVertex(tasks[i].UniqueID)
	}
	for i := range tasks {
		for j := range tasks[i].Successors {
			_ = g.AddEdge(tasks[i].UniqueID, tasks[i].Successors[j].UniqueID)
		}
	}
	file, _ := os.Create("./fp.gv")
	_ = draw.DOT(g, file)
	longestFinish := float64(0)
	for i := range tasks {
		var currentPredID []string
		_ = graph.BFS(g, tasks[i].UniqueID, func(value string) bool {
			currentPredID = append(currentPredID, value)
			return false
		})
		if tasks[i].Successors != nil {
			// Step 1: Forward Pass
			for j := range currentPredID {
				// define ES and EF currentPredID the first task
				tasks[taskIndices[currentPredID[j]]].EarlyStart = 0
				tasks[taskIndices[currentPredID[j]]].EarlyFinish = tasks[taskIndices[currentPredID[j]]].Duration
				if len(tasks[taskIndices[currentPredID[j]]].Predecessors) > 0 {
					for _, predecessor := range tasks[taskIndices[currentPredID[j]]].Predecessors {
						lag := float64(0)
						if predecessor.Lag != "0" {
							unit, value := splitLag(predecessor.Lag)
							lag = float64(value)
							if unit == "mo" {
								lag = float64(value) * 30
							}
						}
						switch predecessor.Type {
						case "FS":
							maxEF := float64(0)
							tasks[taskIndices[currentPredID[j]]].EarlyStart = tasks[taskIndices[predecessor.UniqueID]].EarlyFinish + float64(lag)
							if maxEF > tasks[taskIndices[predecessor.UniqueID]].EarlyFinish {
								tasks[taskIndices[predecessor.UniqueID]].EarlyFinish = maxEF
							} else {
								tasks[taskIndices[currentPredID[j]]].EarlyFinish = tasks[taskIndices[currentPredID[j]]].EarlyStart + tasks[taskIndices[currentPredID[j]]].Duration
								maxEF = tasks[taskIndices[predecessor.UniqueID]].EarlyFinish
							}
						case "FF":
						case "SF":
							tasks[taskIndices[currentPredID[j]]].EarlyStart = tasks[taskIndices[predecessor.UniqueID]].EarlyStart + float64(lag) - 1
							tasks[taskIndices[currentPredID[j]]].EarlyFinish = tasks[taskIndices[currentPredID[j]]].EarlyStart + tasks[taskIndices[currentPredID[j]]].Duration
						case "SS":
							tasks[taskIndices[currentPredID[j]]].EarlyStart = tasks[taskIndices[predecessor.UniqueID]].EarlyStart + float64(lag)
							tasks[taskIndices[currentPredID[j]]].EarlyFinish = tasks[taskIndices[currentPredID[j]]].EarlyStart + tasks[taskIndices[currentPredID[j]]].Duration
						}
						if tasks[taskIndices[currentPredID[j]]].EarlyFinish > longestFinish {
							longestFinish = tasks[taskIndices[currentPredID[j]]].EarlyFinish
						}
					}
				}
				// fmt.Println("[ITERATION]", j, "CurrentTask", tasks[j].UniqueID, "Task: ", tasks[0].UniqueID, "ES: ", tasks[0].EarlyStart, "EF: ", tasks[0].EarlyFinish, "Duration: ", tasks[0].Duration)
				// fmt.Println("[ITERATION]", j, "CurrentTask", tasks[j].UniqueID, "Task: ", tasks[1].UniqueID, "ES: ", tasks[1].EarlyStart, "EF: ", tasks[1].EarlyFinish, "Duration: ", tasks[1].Duration)
				// fmt.Println("[ITERATION]", j, "CurrentTask", tasks[j].UniqueID, "Task: ", tasks[2].UniqueID, "ES: ", tasks[2].EarlyStart, "EF: ", tasks[2].EarlyFinish, "Duration: ", tasks[2].Duration)
				// fmt.Println("[ITERATION]", j, "CurrentTask", tasks[j].UniqueID, "Task: ", tasks[4].UniqueID, "ES: ", tasks[4].EarlyStart, "EF: ", tasks[4].EarlyFinish, "Duration: ", tasks[4].Duration)
				// fmt.Println("[ITERATION]", j, "CurrentTask", tasks[j].UniqueID, "Task: ", tasks[5].UniqueID, "ES: ", tasks[5].EarlyStart, "EF: ", tasks[5].EarlyFinish, "Duration: ", tasks[5].Duration)
			}
		}
		// _ = graph.BFSWithDepth(bp, tasks[i].UniqueID, func(value string, depth int) bool {
		// 	currentPredID = append(currentPredID, value)
		// 	return depth > 0
		// })
		// fmt.Println(currentPredID)
	}
	for i := range tasks {
		fmt.Println("[RESULT]", "Task: ", tasks[i].UniqueID, "ES: ", tasks[i].EarlyStart, "EF: ", tasks[i].EarlyFinish, "Duration: ", tasks[i].Duration)
	}
	// Step 2: Backward Pass
	// LS = LF - duration
	// LF = min(ES of all successors)
	for i := range tasks {
		var currentPredID []string
		_ = graph.BFS(g, tasks[i].UniqueID, func(value string) bool {
			currentPredID = append(currentPredID, value)
			return false
		})
		// Step 2: backward Pass
		for j := len(currentPredID) - 1; j >= 0; j-- {
			// define LS and LF currentPredID the last task
			tasks[taskIndices[currentPredID[j]]].LateFinish = tasks[taskIndices[currentPredID[j]]].EarlyFinish
			tasks[taskIndices[currentPredID[j]]].LateStart = tasks[taskIndices[currentPredID[j]]].LateFinish - tasks[taskIndices[currentPredID[j]]].Duration
			if len(tasks[taskIndices[currentPredID[j]]].Successors) == 0 {
				tasks[taskIndices[currentPredID[j]]].LateFinish = longestFinish
				tasks[taskIndices[currentPredID[j]]].LateStart = tasks[taskIndices[currentPredID[j]]].LateFinish - tasks[taskIndices[currentPredID[j]]].Duration
			}
			if len(tasks[taskIndices[currentPredID[j]]].Predecessors) > 0 {
				for _, predecessor := range tasks[taskIndices[currentPredID[j]]].Predecessors {
					switch predecessor.Type {
					case "FS":
					case "FF":
					case "SF":
					case "SS":
						lag := float64(0)
						if predecessor.Lag != "0" {
							unit, value := splitLag(predecessor.Lag)
							lag = float64(value)
							if unit == "mo" {
								lag = float64(value) * 30
							}
						}
						if tasks[taskIndices[predecessor.UniqueID]].EarlyFinish > tasks[taskIndices[currentPredID[j]]].LateFinish {
							tasks[taskIndices[currentPredID[j]]].LateFinish = tasks[taskIndices[predecessor.UniqueID]].EarlyFinish - float64(lag)
						}
						tasks[taskIndices[currentPredID[j]]].LateStart = tasks[taskIndices[currentPredID[j]]].LateFinish - tasks[taskIndices[currentPredID[j]]].Duration
					}
				}
			}
			if len(tasks[taskIndices[currentPredID[j]]].Successors) > 0 {
				for _, successor := range tasks[taskIndices[currentPredID[j]]].Successors {
					lag := float64(0)
					if successor.Lag != "0" {
						unit, value := splitLag(successor.Lag)
						lag = float64(value)
						if unit == "mo" {
							lag = float64(value) * 30
						}
					}
					switch successor.Type {
					case "FS":
						tasks[taskIndices[currentPredID[j]]].LateFinish = tasks[taskIndices[successor.UniqueID]].LateStart - float64(lag)
						tasks[taskIndices[currentPredID[j]]].LateStart = tasks[taskIndices[currentPredID[j]]].LateFinish - tasks[taskIndices[currentPredID[j]]].Duration
					case "FF":
					case "SF":
						tasks[taskIndices[currentPredID[j]]].LateFinish = tasks[taskIndices[successor.UniqueID]].LateFinish
						if lag > 0 {
							tasks[taskIndices[currentPredID[j]]].LateFinish = tasks[taskIndices[successor.UniqueID]].LateFinish - float64(lag) - 1
						}
						tasks[taskIndices[currentPredID[j]]].LateStart = tasks[taskIndices[currentPredID[j]]].LateFinish - tasks[taskIndices[currentPredID[j]]].Duration

					case "SS":
					}
				}
			}
			fmt.Println("[ITERATION]", j, "CurrentTask", tasks[j].UniqueID, "Task: ", tasks[0].UniqueID, "LS: ", tasks[0].LateStart, "LF: ", tasks[0].LateFinish, "Duration: ", tasks[0].Duration)
			fmt.Println("[ITERATION]", j, "CurrentTask", tasks[j].UniqueID, "Task: ", tasks[1].UniqueID, "LS: ", tasks[1].LateStart, "LF: ", tasks[1].LateFinish, "Duration: ", tasks[1].Duration)
			fmt.Println("[ITERATION]", j, "CurrentTask", tasks[j].UniqueID, "Task: ", tasks[2].UniqueID, "LS: ", tasks[2].LateStart, "LF: ", tasks[2].LateFinish, "Duration: ", tasks[2].Duration)
		}
	}

	// Step 3: Calculate Total Slack
	for i := range tasks {
		tasks[i].TotalSlack = tasks[i].LateFinish - tasks[i].EarlyFinish
		fmt.Println("[RESULT]", "Task: ", tasks[i].UniqueID, "EF", tasks[i].EarlyFinish, "LS: ", tasks[i].LateStart, "LF: ", tasks[i].LateFinish, "Duration: ", tasks[i].Duration, "TotalSlack: ", tasks[i].TotalSlack)
		if tasks[i].TotalSlack == 0 {
			tasks[i].CriticalPath = true
		}
	}
	fmt.Println(longestFinish, "longestSlack")

	return tasks
}
