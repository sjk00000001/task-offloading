package to

import (
	"fmt"
	"math"
	"math/rand"
	"task-offloading/algorithm"
	"task-offloading/comnet"
	"task-offloading/simulate"
	"task-offloading/structs"
	"task-offloading/util"
	"time"
)

func ExecuteSimulation(M int, N int) []structs.TaskOut {
	taskOuts := make([]structs.TaskOut, M)
	// set random seed
	rand.Seed(time.Now().UnixNano())
	// get TFD (Task Feature Data)
	taskData := simulate.GetSiTask(M)
	// get NPD (Network Performance Data)
	netData := simulate.GetSiNet(N)
	// get R (channel rate)
	RMatrix := simulate.GetR(N)
	delayMatrix := simulate.GetSimpleMatrix(taskData, netData, RMatrix)
	decisionMap := algorithm.KM(delayMatrix)
	for i := range taskOuts {
		deadline := taskData[i].Deadline
		localDelay := delayMatrix[i][i]
		toDelay := delayMatrix[i][decisionMap[i]]
		taskOuts[i] = structs.TaskOut{
			ID:         string(rune(i)),
			ArriveAt:   string(rune(i)),
			OffloadTo:  string(rune(decisionMap[i])),
			Type:       FaceRkg,
			Deadline:   util.DecimalFormat(deadline),
			LocalDelay: util.DecimalFormat(localDelay),
			LocalEff:   util.DecimalFormat(math.Max(0, 1-localDelay/deadline)),
			TODelay:    util.DecimalFormat(toDelay),
			TOEff:      util.DecimalFormat(math.Max(0, 1-toDelay/deadline)),
		}
	}
	return taskOuts
}

func Execute(taskIns []structs.TaskIn) []structs.TaskOut {
	return execute(taskIns)
}

/* execute task offloading decision */
func execute(taskIns []structs.TaskIn) []structs.TaskOut {
	fmt.Println("[TO] decision recv:", taskIns)
	// pre-pend
	comNet, err := comnet.PullRecentComNet()
	if err != nil {
		util.PrintError(err.Error())
		return nil
	}
	taskOuts := make([]structs.TaskOut, len(taskIns))
	nodes := comNet.Nodes
	links := comNet.Links
	M := len(taskIns)
	N := len(nodes)
	// get core
	pairMap := getPairMap(taskIns, links)
	graph := getGraph(pairMap, M, N)
	fmt.Println("[TO] graph:", graph)
	// offloading
	decisionMap := algorithm.KM(graph)
	reformatDecisionMap(decisionMap, M, N)
	util.PrintMap(decisionMap)
	// post-pend
	for i, taskIn := range taskIns {
		deadline := TypeDeadlineDict[taskIn.Type]
		src := pairMap[taskIn][0].Src
		dest := decisionMap[src]
		fmt.Println("hello:",dest)
		offloadTo := comnet.Decode(dest)
		lcDelay := -graph[src][comnet.Encode(taskIn.ArriveAt)]
		toDelay := -graph[src][dest]
		taskOuts[i] = structs.TaskOut{
			ID:         taskIn.ID,
			ArriveAt:   taskIn.ArriveAt,
			OffloadTo:  offloadTo,
			Type:       taskIn.Type,
			Deadline:   util.DecimalFormat(deadline),
			LocalDelay: util.DecimalFormat(lcDelay),
			LocalEff:   util.DecimalFormat(math.Max(0, 1-lcDelay/deadline) * 100),
			TODelay:    util.DecimalFormat(toDelay),
			TOEff:      util.DecimalFormat(math.Max(0, 1-toDelay/deadline) * 100),
		}
	}
	fmt.Println("[TO] decision result:", taskOuts)
	var sum float64
	for _, taskOut := range taskOuts{
		sum += taskOut.TODelay
	}
	fmt.Println("使用KM算法后所有任务的计算卸载时间总延迟为",sum)
	return taskOuts
}

func getPairMap(taskIns []structs.TaskIn, links []structs.Link) map[structs.TaskIn][]structs.Pair {
	pairMap := make(map[structs.TaskIn][]structs.Pair)
	feed := 0
	// get taskIn-pair map
	var sum2 float64 = 0
	for _, taskIn := range taskIns {
		taskArrivalIP := taskIn.ArriveAt
		dataSize := TypeDataSizeDict[taskIn.Type]
		ProDelay := TypeDeadlineDict[taskIn.Type] / 2
		// offload
		for _, link := range links {
			if link.NodeFrom == taskArrivalIP {
				nodeTo := link.NodeTo
				node := comnet.NodeIPMap[nodeTo]
				UpDelay := dataSize / link.Rate
				pair := genPair(taskIn, node, feed, UpDelay+ProDelay+link.EsDelay)
				pairMap[taskIn] = append(pairMap[taskIn], pair)
			}
		}
		// local
		pair := genPair(taskIn, comnet.NodeIPMap[taskArrivalIP], feed, ProDelay)
		sum2 += pair.Cost
		pairMap[taskIn] = append(pairMap[taskIn], pair)
		feed++
	}
	fmt.Println("未使用KM算法时所有任务的计算卸载时间总延迟为",sum2)
	return pairMap
}

func genPair(taskIn structs.TaskIn, node structs.Node, index int, Cost0 float64) structs.Pair {
	CGDelay := predExecutable(taskIn.Type, node)
	return structs.Pair{
		Src:  index,
		Des:  comnet.Encode(node.IP),
		Cost: CGDelay + Cost0,
	}
}

func getGraph(pairMap map[structs.TaskIn][]structs.Pair, M int, N int) [][]float64 {
	E := N
	if M > N {
		/* Min-Span */
		E = N * ((M-1)/N + 1)
		/* Max-Span */
		//E = M * N
	}
	graph := util.MakeSlicesFloat64(E, E)
	// get bi-part graph
	for _, pairs := range pairMap {
		for _, pair := range pairs {
			graph[pair.Src][pair.Des] = -pair.Cost
		}
	}
	if M > N {
		for m := 0; m < M; m++ {
			for n := N; n < E; n++ {
				graph[m][n] = graph[m][n-N*(n/N)]
			}
		}
	}
	return graph
}

func reformatDecisionMap(decisionMap map[int]int, M int, N int) {
	if M > N {
		for k, v := range decisionMap {
			decisionMap[k] = v - N*(v/N)
		}
	}
}
