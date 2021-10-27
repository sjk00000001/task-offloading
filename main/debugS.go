package main

import (
	"fmt"
	"task-offloading/to"
)

func main() {
	fmt.Println(" recv message: ", to.ExecuteSimulation(50,100))
	fmt.Scanln()
}
