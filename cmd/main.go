package main

import (
	"fmt"
	"time"

	"github.com/elecbug/nw-graph-tester/internal/network"
)

func main() {
	n := network.GenerateGossipSubNetwork(51, 6, 5, 12, 6000, 100)
	if n == nil {
		fmt.Println("Failed to generate network")
		return
	}
	n.Print()

	n.Nodes[0].Relay(1, nil)

	time.Sleep(time.Second * 20)

	sum := 0
	for i := range n.Nodes {
		sum += len(n.Nodes[i].DuplicateMap[1])
		n.Nodes[i].PrintRelayState()
		fmt.Printf("> Degree - DuplicateMsg = %d\n", len(n.Nodes[i].Connections)-len(n.Nodes[i].DuplicateMap[1]))
	}

	fmt.Printf("Total average duplicates for relay 1: %f\n", float64(sum)/float64(len(n.Nodes)-1))
}
