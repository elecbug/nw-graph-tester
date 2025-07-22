package main

import (
	"fmt"

	"github.com/elecbug/nw-graph-tester/internal/network"
)

func main() {
	n := network.GenrateRandomNetwork(10, 20, 100, 50)
	if n == nil {
		fmt.Println("Failed to generate network")
		return
	}
	n.Print()
}
