package main

import (
	"os"
	"task-offloading/client"
	"task-offloading/comnet"
	"task-offloading/util"
)

func main() {
	portedIPs := os.Args[1:]
	if len(portedIPs) == 2 && util.IsPortedIP(portedIPs) {
		go comnet.StartZmqSubscriber(portedIPs[0])
		go client.StartGrpcServer(portedIPs[1])
	} else {
		util.PrintError("[OS] parameter is invalid, run debug mode now")
		go comnet.StartZmqSubscriber("")
		go client.StartGrpcServer("")
	}
	select {}
}
