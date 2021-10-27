package simulate

import (
	"math"
	"task-offloading/structs"
	"task-offloading/util"
)

/* solve optimal delay matrix (find the optimal computing resources to minimize the sum of delay) */
func GetSimpleMatrix(taskData []structs.SiTask, netData []structs.SiNet, RMatrix [][]float64) [][]float64 {
	M, N := len(taskData), len(netData)
	delayMatrix := util.MakeSlicesFloat64(N, N)
	// padding the matrix to a square matrix
	if M < N {
		for m := M; m < N; m++ {
			for n := 0; n < N; n++ {
				delayMatrix[m][n] = 0
			}
		}
	}
	if M <= N {
		for m := 0; m < M; m++ {
			mData := netData[m]
			taskData := taskData[m]
			// upper limit is optimal
			fmmax := math.Min(mData.LimitPerDocker[1], mData.RestComputing)
			for n := 0; n < N; n++ {
				nData := netData[n]
				Rmn := RMatrix[m][n]
				if n == m {
					delayMatrix[m][n] = - (taskData.WorkLoad / fmmax)
				} else {
					// it must be executed locally if it cannot meet its deadline when offloading
					if taskData.Deadline-taskData.DataSize/Rmn <= 0 {
						delayMatrix[m][n] = math.Inf(-1)
					} else {
						fnmax := math.Min(nData.LimitPerDocker[1], nData.RestComputing)
						delayMatrix[m][n] = - (taskData.DataSize/Rmn + taskData.WorkLoad/fnmax)
					}
				}
			}
		}
	} else {
		return nil
	}
	return delayMatrix
}
