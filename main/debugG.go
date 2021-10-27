package main

import "task-offloading/client"

func main() {
	for i := 1; i <= 100; i++ {
		go client.StartGrpcClient(i)
	}
	select {}

}
