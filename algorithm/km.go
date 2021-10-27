package algorithm

import (
	"math"
	"task-offloading/util"
	"time"
)

var wx []float64
var wy []float64
var cx []int
var cy []int
var slack []float64
var vx []bool
var vy []bool

var iteration int

/* KM algorithm, an algorithm with Graph Theory with O(N^3) time complexity */
func KM(graph [][]float64) map[int]int {
	E := len(graph)
	res := make(map[int]int)
	val := 0.0
	wx = make([]float64, E)
	wy = make([]float64, E)
	cx = make([]int, E)
	cy = make([]int, E)
	slack = make([]float64, 2*E)
	iteration = 0
	// array initialization
	for i := 0; i < E; i++ {
		cx[i] = -1
		cy[i] = -1
		for j := 0; j < E; j++ {
			wx[i] = math.Max(wx[i], graph[i][j])
		}
	}
	// execution time statistics
	startTime := time.Now()
	for i := 0; i < E; i++ {
		for j := 0; j < 2*E; j++ {
			slack[j] = math.MaxFloat64
		}
		for {
			vx = make([]bool, 2*E)
			vy = make([]bool, 2*E)
			// execute Hungary algorithm in a loop, adjust the value of label by the value of cursor if the augmenting path is not found
			if !dfs(i, graph) {
				cursor := math.MaxFloat64
				for j := 0; j < E; j++ {
					if !vy[j] && cursor > slack[j] {
						cursor = slack[j]
					}
				}
				// decrease the value of X label by the value of cursor
				for j := 0; j < E; j++ {
					if vx[j] {
						wx[j] -= cursor
					}
				}
				// increase the value of Y label by the value of cursor
				for j := 0; j < E; j++ {
					if vy[j] {
						wy[j] += cursor
					} else {
						slack[j] -= cursor
					}
				}
			} else {
				break
			}
		}
	}
	duration := int(time.Since(startTime) / 1e6)
	// clean the result
	for i := 0; i < E; i++ {
		res[i] = cx[i]
		val -= graph[i][cx[i]]
	}
	util.PrintResult(res, val, iteration, duration)
	return res
}

/* Hungary algorithm, an algorithm to find augmenting path implemented by dfs */
func dfs(u int, graph [][]float64) bool {
	iteration++
	vx[u] = true
	for v := 0; v < len(cy); v++ {
		if !vy[v] {
			temp := wx[u] + wy[v] - graph[u][v]
			// match successfully when the edge weight is equal to the sum of the value of labels
			if temp < 1e-8 {
				vy[v] = true
				if cy[v] == -1 || dfs(cy[v], graph) {
					cx[u] = v
					cy[v] = u
					return true
				}
			} else {
				slack[v] = math.Min(slack[v], temp)
			}
		}
	}
	return false
}
