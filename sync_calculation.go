package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
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
	maps := make(map[string]Task)
	for i := 0; i < len(tasks); i++ {
		if len(tasks[i].Predecessors) != 0 {
			for j := 0; j < len(tasks[i].Predecessors); j++ {
				for k := 0; k < len(tasks); k++ {
					// if tasks[i].ID == "16" || tasks[i].ID == "15" {
					if tasks[k].ID == tasks[i].Predecessors[j].ID {
						tasks[i].StartDate, tasks[i].EndDate = CalculateStartFinish(tasks[i], tasks[k])
					}
					// }
				}
			}
		}
		maps[tasks[i].ID] = tasks[i]
	}
	g := graph.New(graph.IntHash, graph.Directed(), graph.PreventCycles())
	gw := graph.New(graph.StringHash, graph.Weighted())
	for i := 0; i < len(tasks); i++ {
		if len(tasks[i].Predecessors) != 0 {
			for j := 0; j < len(tasks[i].Predecessors); j++ {
				for k := 0; k < len(tasks); k++ {
					if tasks[k].ID == tasks[i].Predecessors[j].ID {
						ID, err := strconv.Atoi(tasks[i].ID)
						if err != nil {
							panic(err)
						}
						g.AddVertex(ID)
						gw.AddVertex(tasks[i].ID)
						idPredessesor, err := strconv.Atoi(tasks[k].ID)
						if err != nil {
							panic(err)
						}
						g.AddVertex(idPredessesor)
						gw.AddVertex(tasks[k].ID)
						err = g.AddEdge(ID, idPredessesor, graph.EdgeAttribute("color", "red"))
						if err != nil {
							panic(err)
						}
					}
				}
			}
		}
	}
	err := graph.DFS(g, 2, func(value int) bool {
		id := strconv.Itoa(value)
		task := maps[id]
		fmt.Println(id, task.Predecessors)
		return false
	})
	if err != nil {
		panic(err)
	}
	mst, err := graph.MinimumSpanningTree(gw)
	if err != nil {
		panic(err)
	}
	fmt.Println(mst.Size())
	file, _ := os.Create("my-graph.gv")
	_ = draw.DOT(g, file)
	return tasks
}
