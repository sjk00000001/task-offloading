package simulate

import (
	"math"
	"math/cmplx"
	"math/rand"
	"task-offloading/structs"
	"task-offloading/util"
)

/* Task Feature parameters */
// bases workload of task (10^6 cpu cycles)
var baseWorkload = 200.0
var timesWorkload = 4.0

// data-workload ratio
var dataRatio = 0.2

// upper and lower deadlines (ms)
var deadlineLimit = [2]float64{10, 20}

/* Network Performance parameters */
var computingResourceLimit = [2]float64{50, 150}
var storageResourceLimit = [2]float64{30, 120}

// number of containers allowed per node
var maxDockerPerNode = 2

// allowable amount of overflow resources
var overflowResource = [2]float64{5, 15}

/* R parameter */
// bandwidth (Hz)
var B = 2e6

// channel loss
var pathLoss = 3.5

// size of space (m^2)
var square = 100.0

// power of white noise (Additive Gauss White Noise)
var noise = 1e-15

// transmission power threshold
var transPowerLimit = [2]float64{10.0, 20.0}

/* two-dimension coordinate */
type Coordinate2D struct {
	// X-axis
	X float64
	// Y-axis
	Y float64
}

/* simulate TFD (F, C, E) */
func GetSiTask(M int) []structs.SiTask {
	taskData := make([]structs.SiTask, M)
	for m := 0; m < M; m++ {
		var siTask structs.SiTask
		siTask.WorkLoad = math.Ceil(rand.Float64()*timesWorkload) * baseWorkload
		siTask.DataSize = dataRatio * siTask.WorkLoad
		siTask.Deadline = util.IntervalRandGenerator(deadlineLimit)
		taskData[m] = siTask
	}
	return taskData
}

/* simulate NPD (W, D, T) */
func GetSiNet(N int) []structs.SiNet {
	netData := make([]structs.SiNet, N)
	for n := 0; n < N; n++ {
		var siNet structs.SiNet
		siNet.RestComputing = util.IntervalRandGenerator(computingResourceLimit) * (1 - rand.Float64())
		siNet.RestStorage = util.IntervalRandGenerator(storageResourceLimit) * (1 - rand.Float64())
		siNet.LimitPerDocker = [2]float64{computingResourceLimit[0]/float64(maxDockerPerNode) - util.IntervalRandGenerator(overflowResource),
			computingResourceLimit[1]/float64(maxDockerPerNode) - util.IntervalRandGenerator(overflowResource)}
		netData[n] = siNet
	}
	return netData
}

/* simulate R */
func GetR(N int) [][]float64 {
	RMatrix := util.MakeSlicesFloat64(N, N)
	RMatrix[0][0] = 0
	coordinate := make([]Coordinate2D, N)
	for i := 0; i < N; i++ {
		coordinate[i] = Coordinate2D{X: rand.Float64() * math.Sqrt(square), Y: rand.Float64() * math.Sqrt(square)}
	}
	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			h := complex(rand.Float64(), rand.Float64())
			X := 20 * math.Log10(cmplx.Abs(h))
			d := math.Sqrt(math.Pow(coordinate[i].X-coordinate[j].X, 2) + math.Pow(coordinate[i].Y-coordinate[j].Y, 2))
			if d == 0 {
				RMatrix[i][j] = 0
			} else {
				Pl := 10 * pathLoss * math.Log10(d)
				Pw := util.IntervalRandGenerator(transPowerLimit)
				Pt := 10 * math.Log10(Pw)
				Pr := Pt - Pl + X
				Pr = math.Pow(10, (Pr-30)/10)
				RMatrix[i][j] = B * math.Log2(1+Pr/(B*noise)) * 1e-6
			}
		}
	}
	return RMatrix
}
