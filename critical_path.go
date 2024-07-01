package main

import (
	"fmt"

	"github.com/dominikbraun/graph"
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

	// file, _ := os.Create("./gTree.gv")
	// _ = draw.DOT(gTree, file)
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
				longestFinish = tasks[taskIndices[currentPredID[j]]].EarlyFinish
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
					}
					// fmt.Println(longestFinish, "longestEF")
					if tasks[taskIndices[currentPredID[j]]].EarlyFinish > longestFinish {
						longestFinish = tasks[taskIndices[currentPredID[j]]].EarlyFinish
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
			// fmt.Println("[FIRSTINIT]", tasks[taskIndices[currentPredID[j]]].UniqueID, tasks[taskIndices[currentPredID[j]]].LateFinish, tasks[taskIndices[currentPredID[j]]].LateStart, tasks[taskIndices[currentPredID[j]]].Duration)
			if tasks[i].DataType != "item" {
				if len(tasks[taskIndices[currentPredID[j]]].Successors) == 0 {
					// fmt.Println("[LASTTASK]", longestFinish)
					tasks[taskIndices[currentPredID[j]]].LateFinish = longestFinish
					tasks[taskIndices[currentPredID[j]]].LateStart = tasks[taskIndices[currentPredID[j]]].LateFinish - tasks[taskIndices[currentPredID[j]]].Duration
					// fmt.Println("[SECONDINIT]", tasks[taskIndices[currentPredID[j]]].UniqueID, tasks[taskIndices[currentPredID[j]]].LateFinish, tasks[taskIndices[currentPredID[j]]].LateStart, tasks[taskIndices[currentPredID[j]]].Duration)
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

							if tasks[taskIndices[currentPredID[j]]].LateFinish < 0 {
								tasks[taskIndices[currentPredID[j]]].LateFinish = 0
							}
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
			}
		}
	}

	for i := range tasks {
		if tasks[i].DataType == "header" {
			for _, child := range tasks[i].ChildUniqueIDs {
				// fmt.Println("[CHILD]", child, tasks[taskIndices[child]].LateFinish, tasks[taskIndices[child]].LateStart, tasks[taskIndices[child]].Duration)
				if tasks[taskIndices[child]].LateFinish > tasks[i].LateFinish {
					tasks[i].LateFinish = tasks[taskIndices[child]].LateFinish
				}
				if tasks[taskIndices[child]].LateStart > tasks[i].LateStart {
					tasks[i].LateStart = tasks[taskIndices[child]].LateStart
				}
			}
		}
	}

	// Step 3: Calculate Total Slack
	for i := range tasks {
		tasks[i].TotalSlack = tasks[i].LateFinish - tasks[i].EarlyFinish
		if tasks[i].TotalSlack < 0 {
			tasks[i].TotalSlack = 0
		}
		if tasks[i].TotalSlack == 0 {
			tasks[i].CriticalPath = true
		}
		// fmt.Println("[RESULT]", "Task: ", tasks[i].UniqueID, "ES", tasks[i].EarlyStart, "EF", tasks[i].EarlyFinish, "LS: ", tasks[i].LateStart, "LF: ", tasks[i].LateFinish, "Duration: ", tasks[i].Duration, "TotalSlack: ", tasks[i].TotalSlack)
	}
	// fmt.Println(longestFinish, "longestSlack")

	return tasks
}

func calculateCriticalPath(tasks []Task) []Task {
	taskIndices := make(map[string]int)
	for i, task := range tasks {
		taskIndices[task.UniqueID] = i
	}
	g := graph.New(graph.StringHash, graph.Directed(), graph.Acyclic())
	gTree := graph.New(graph.StringHash, graph.Directed(), graph.Acyclic())
	for i := range tasks {
		_ = g.AddVertex(tasks[i].UniqueID)
		_ = gTree.AddVertex(tasks[i].UniqueID)
	}
	for i := range tasks {
		for j := range tasks[i].Successors {
			_ = g.AddEdge(tasks[i].UniqueID, tasks[i].Successors[j].UniqueID)
		}
		for k := range tasks[i].ChildUniqueIDs {
			_ = gTree.AddEdge(tasks[i].UniqueID, tasks[i].ChildUniqueIDs[k])
		}
	}

	longestFinish := float64(0)
	for j := len(tasks) - 1; j >= 0; j-- {
		var CurrentChild []string
		_ = graph.DFS(gTree, tasks[j].UniqueID, func(value string) bool {
			CurrentChild = append(CurrentChild, value)
			return false
		})
		tasks[taskIndices[tasks[j].UniqueID]].EarlyStart = 0
		tasks[taskIndices[tasks[j].UniqueID]].EarlyFinish = tasks[taskIndices[tasks[j].UniqueID]].Duration

		if longestFinish < tasks[j].EarlyFinish {
			longestFinish = tasks[j].EarlyFinish
		}
	}

	for j := len(tasks) - 1; j >= 0; j-- {
		var CurrentChild []string
		_ = graph.DFS(gTree, tasks[j].UniqueID, func(value string) bool {
			CurrentChild = append(CurrentChild, value)
			return false
		})

		tasks[taskIndices[tasks[j].UniqueID]].LateFinish = longestFinish
		tasks[taskIndices[tasks[j].UniqueID]].LateStart = tasks[taskIndices[tasks[j].UniqueID]].LateFinish - tasks[taskIndices[tasks[j].UniqueID]].Duration
		if len(tasks[j].ChildUniqueIDs) > 0 {
			for _, child := range tasks[j].ChildUniqueIDs {
				if tasks[taskIndices[child]].LateFinish > tasks[j].LateFinish {
					tasks[j].LateFinish = tasks[taskIndices[child]].LateFinish
				}
				if tasks[taskIndices[child]].LateStart < tasks[j].LateStart {
					tasks[j].LateStart = tasks[taskIndices[child]].LateStart
				}
			}
		}
		fmt.Println("[RESULT]", "Task: ", tasks[j].UniqueID, "ES", tasks[j].EarlyStart, "EF", tasks[j].EarlyFinish, "LS: ", tasks[j].LateStart, "LF: ", tasks[j].LateFinish, "Duration: ", tasks[j].Duration, "TotalSlack: ", tasks[j].TotalSlack)
	}

	for i := range tasks {
		var currentPredID []string
		_ = graph.BFS(g, tasks[i].UniqueID, func(value string) bool {
			currentPredID = append(currentPredID, value)
			return false
		})
		if len(tasks[i].Successors) != 0 {
			// Step 1: Forward Pass
			for j := range currentPredID {
				if tasks[taskIndices[currentPredID[j]]].UniqueID != tasks[i].UniqueID {
					tasks[taskIndices[currentPredID[j]]].EarlyStart = tasks[i].EarlyFinish
					tasks[taskIndices[currentPredID[j]]].EarlyFinish = tasks[taskIndices[currentPredID[j]]].EarlyStart + tasks[taskIndices[currentPredID[j]]].Duration
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
								tasks[taskIndices[currentPredID[j]]].EarlyStart = tasks[taskIndices[predecessor.UniqueID]].EarlyStart + float64(lag)
								tasks[taskIndices[currentPredID[j]]].EarlyFinish = tasks[taskIndices[currentPredID[j]]].EarlyStart + tasks[taskIndices[currentPredID[j]]].Duration
							case "SF":
								tasks[taskIndices[currentPredID[j]]].EarlyStart = tasks[taskIndices[predecessor.UniqueID]].EarlyStart + float64(lag) - 1
								tasks[taskIndices[currentPredID[j]]].EarlyFinish = tasks[taskIndices[currentPredID[j]]].EarlyStart + tasks[taskIndices[currentPredID[j]]].Duration
							case "SS":
								tasks[taskIndices[currentPredID[j]]].EarlyStart = tasks[taskIndices[predecessor.UniqueID]].EarlyStart + float64(lag)
								tasks[taskIndices[currentPredID[j]]].EarlyFinish = tasks[taskIndices[currentPredID[j]]].EarlyStart + tasks[taskIndices[currentPredID[j]]].Duration
							}
						}
						if tasks[taskIndices[currentPredID[j]]].EarlyFinish > longestFinish {
							longestFinish = tasks[taskIndices[currentPredID[j]]].EarlyFinish
						}
					}
				}
			}
		}

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
			if len(tasks[taskIndices[currentPredID[j]]].Successors) > 0 {
				fmt.Println("BACKWARD", tasks[taskIndices[currentPredID[j]]].UniqueID, tasks[i].UniqueID, tasks[i].EarlyFinish, currentPredID[j])
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

						if tasks[taskIndices[currentPredID[j]]].LateFinish < 0 {
							tasks[taskIndices[currentPredID[j]]].LateFinish = 0
						}
						fmt.Println("FS", tasks[taskIndices[currentPredID[j]]].UniqueID, tasks[taskIndices[successor.UniqueID]].UniqueID, tasks[taskIndices[currentPredID[j]]].LateFinish)
					case "FF":
						// tasks[taskIndices[currentPredID[j]]].LateFinish = tasks[taskIndices[successor.UniqueID]].LateFinish
						fmt.Println("FF", tasks[taskIndices[currentPredID[j]]].UniqueID, tasks[taskIndices[successor.UniqueID]].UniqueID, tasks[taskIndices[currentPredID[j]]].LateFinish)
					case "SF":
						tasks[taskIndices[currentPredID[j]]].LateFinish = tasks[taskIndices[successor.UniqueID]].LateFinish
						// if lag > 0 {
						// 	tasks[taskIndices[currentPredID[j]]].LateFinish = tasks[taskIndices[successor.UniqueID]].LateFinish - float64(lag) - 1
						// }
						// tasks[taskIndices[currentPredID[j]]].LateStart = tasks[taskIndices[currentPredID[j]]].LateFinish - tasks[taskIndices[currentPredID[j]]].Duration
					case "SS":
						// tasks[taskIndices[currentPredID[j]]].LateStart = tasks[taskIndices[successor.UniqueID]].LateStart
						// tasks[taskIndices[currentPredID[j]]].LateFinish = tasks[taskIndices[successor.UniqueID]].LateFinish
						fmt.Println("SS", tasks[taskIndices[currentPredID[j]]].UniqueID, tasks[taskIndices[successor.UniqueID]].UniqueID, tasks[taskIndices[currentPredID[j]]].LateFinish)
					}
				}
				if len(tasks[taskIndices[currentPredID[j]]].ChildUniqueIDs) != 0 {
					for _, child := range tasks[taskIndices[currentPredID[j]]].ChildUniqueIDs {
						tasks[taskIndices[child]].LateFinish = tasks[i].LateFinish
						tasks[taskIndices[child]].LateStart = tasks[i].LateFinish - tasks[taskIndices[child]].Duration
					}
				}
			}
		}
	}

	// Step 3: Calculate Total Slack
	for i := range tasks {
		tasks[i].TotalSlack = tasks[i].LateFinish - tasks[i].EarlyFinish
		if tasks[i].TotalSlack < 0 {
			tasks[i].TotalSlack = 0
		}
		if tasks[i].TotalSlack == 0 {
			tasks[i].CriticalPath = true
		}
		fmt.Println("[RESULT]", "Task: ", tasks[i].UniqueID, "ES", tasks[i].EarlyStart, "EF", tasks[i].EarlyFinish, "LS: ", tasks[i].LateStart, "LF: ", tasks[i].LateFinish, "Duration: ", tasks[i].Duration, "TotalSlack: ", tasks[i].TotalSlack)
	}
	fmt.Println(longestFinish)
	return tasks
}
