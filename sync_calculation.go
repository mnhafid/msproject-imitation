package main

import (
	"fmt"
	"time"

	"github.com/dominikbraun/graph"
)

type TaskResponses struct {
	Count int
	Tasks []TaskResponse `json:"task"`
}

type TaskResponse struct {
	ID           string
	Name         string
	UniqueID     string
	Start        string
	Finish       string
	PlanFinish   string
	Duration     string
	Work         int
	Cost         int
	Milestone    bool
	Predecessors []Predecessor `json:"predecessors"`
	Successors   []Successor   `json:"successors"`
	DurationType string
}

const (
	formatTime = "2006-01-02T15:04"
	formatDate = "2006-01-02"
)

func CalculateAllProjectTask(tasksReponse []TaskResponse) []Task {
	var tasks []Task
	// Calculate all project task
	for i := 0; i < len(tasksReponse); i++ {
		if tasksReponse[i].Milestone {
			continue
		}
		start, _ := time.Parse(formatTime, tasksReponse[i].Start)
		finish, _ := time.Parse(formatTime, tasksReponse[i].Finish)
		task := Task{
			ID:           tasksReponse[i].ID,
			Description:  tasksReponse[i].Name,
			UniqueID:     tasksReponse[i].UniqueID,
			Duration:     tasksReponse[i].Duration,
			Work:         tasksReponse[i].Work,
			Cost:         tasksReponse[i].Cost,
			DurationType: tasksReponse[i].DurationType,
			Predecessors: tasksReponse[i].Predecessors,
			Successors:   tasksReponse[i].Successors,
			StartDate:    start,
			EndDate:      finish,
		}
		tasks = append(tasks, task)
	}
	// maps := make(map[string]Task)
	for i := 0; i < len(tasks); i++ {
		if len(tasks[i].Predecessors) != 0 {
			for j := 0; j < len(tasks[i].Predecessors); j++ {
				for k := 0; k < len(tasks); k++ {
					if tasks[k].ID == tasks[i].Predecessors[j].ID {
						tasks[i].PlanStart, tasks[i].PlanFinish = CalculateStartFinish(tasks[i], tasks[k])
					}
				}
			}
		}
	}
	// g := graph.New(graph.StringHash, graph.Directed(), graph.PreventCycles())
	// // gw := graph.New(graph.StringHash, graph.Weighted())
	// for i := 0; i < len(tasks); i++ {
	// 	if len(tasks[i].Predecessors) != 0 {
	// 		for j := 0; j < len(tasks[i].Predecessors); j++ {
	// 			for k := 0; k < len(tasks); k++ {
	// 				if tasks[k].ID == tasks[i].Predecessors[j].ID {
	// 					tasks[i].PlanStart = time.Date(2022, 10, 13, 16, 0, 0, 0, time.UTC)
	// 					tasks[i].PlanFinish = time.Date(2022, 10, 13, 16, 0, 0, 0, time.UTC)
	// 					g.AddVertex(tasks[i].ID)
	// 					g.AddVertex(tasks[k].ID)
	// 					err := g.AddEdge(tasks[i].ID, tasks[k].ID, graph.EdgeAttribute("color", "red"))
	// 					if err != nil {
	// 						panic(err)
	// 					}
	// 				}
	// 			}
	// 		}
	// 	}
	// }
	// file, _ := os.Create("my-graph.gv")
	// _ = draw.DOT(g, file)
	return tasks
}

func DFS[K comparable, T any](g graph.Graph[K, T], start K, visit func(K) bool) error {
	adjacencyMap, err := g.AdjacencyMap()
	if err != nil {
		return fmt.Errorf("could not get adjacency map: %w", err)
	}

	if _, ok := adjacencyMap[start]; !ok {
		return fmt.Errorf("could not find start vertex with hash %v", start)
	}

	stack := make([]K, 0)
	visited := make(map[K]bool)

	stack = append(stack, start)

	for len(stack) > 0 {
		currentHash := stack[len(stack)-1]

		stack = stack[:len(stack)-1]

		if _, ok := visited[currentHash]; !ok {
			// Stop traversing the graph if the visit function returns true.
			if stop := visit(currentHash); stop {
				break
			}
			visited[currentHash] = true

			for adjacency := range adjacencyMap[currentHash] {
				stack = append(stack, adjacency)
			}
		}
	}

	return nil
}
