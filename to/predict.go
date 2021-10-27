package to

import (
	"math"
	"task-offloading/structs"
)

// return -1 -> can't execute, 0 -> blank execution, 1 -> ready execution
func predExecutable(typo string, node structs.Node) float64 {
	r := TypeResourceDict[typo]
	// kernels (K)
	x := node.RestComputing[0]
	// frequency (GHz)
	y := node.RestComputing[1]
	// utilization (%)
	u1 := node.RestComputing[2]
	// storage (MB)
	s := node.RestStorage[0]
	// utilization (%)
	u2 := node.RestStorage[1]
	b1 := r[0] <= x*(1-u1) && r[1] <= y && r[2] <= s*(1-u2)
	if !b1 {
		return math.MaxInt16
	}
	for _, cg := range node.CGs {
		if cg.Type == typo {
			// Busy status must be not included in the ConstDelay slice. I think this can be optimized.
			for i := len(ConstDelay); i > 0; i-- {
				if cg.Utility[i-1] > 0 {
					return ConstDelay[i]
				}
			}
			return ConstDelay[0]
		}
	}
	return 0
}
