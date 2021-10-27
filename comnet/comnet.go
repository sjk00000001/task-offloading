package comnet

import (
	"encoding/json"
	"errors"
	"fmt"
	"task-offloading/structs"
	"task-offloading/util"
)

var NodeIPMap = make(map[string]structs.Node)

var comNet structs.ComNet
var comNetQueue []structs.ComNet
var nodeEncoder = make(map[string]int)
var nodeDecoder = make(map[int]string)

func Encode(s string) int {
	return nodeEncoder[s]
}

func Decode(i int) string {
	return nodeDecoder[i]
}

func PullRecentComNet() (structs.ComNet, error) {
	if len(comNetQueue) == 0 {
		err := errors.New("[ComNet] no ComNet in the queue")
		return comNet, err
	} else {
		fmt.Println("[ComNet] pull ComNet: ", comNet)
		for i, node := range comNet.Nodes {
			nodeEncoder[node.IP] = i
			nodeDecoder[i] = node.IP
			NodeIPMap[node.IP] = node
		}
		return comNetQueue[len(comNetQueue)-1], nil
	}
}

func StoreComNet(s string) {
	data := util.String2Bytes(s)
	err := json.Unmarshal(data, &comNet)
	if err != nil {
		fmt.Println("[ComNet] store failed", err)
	} else {
		fmt.Println("[ComNet] store ComNet: ", comNet)
		comNetQueue = append(comNetQueue, comNet)
		if len(comNetQueue) > 15 {
			comNetQueue = comNetQueue[1:]
		}
	}
}
