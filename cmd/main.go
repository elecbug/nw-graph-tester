package main

import (
	"fmt"
	"time"

	"github.com/elecbug/p2p-broadcast-tester/internal/network"
	"github.com/elecbug/p2p-broadcast-tester/internal/p2p"
)

func main() {
	n := network.GenerateGossipSubNetwork(network.NetworkConfig{
		NodeCount:    10000,
		DLow:         8,
		D:            10,
		DHigh:        12,
		MaxNodeDelay: 100,
		MaxLinkDelay: 100,
	})

	if n == nil {
		fmt.Println("Failed to generate network")
		return
	}
	// n.Print()

	n.Nodes[0].Broadcast(1, p2p.TikTokPublish)

	time.Sleep(time.Second * 10)

	sum := 0
	for i := range n.Nodes {
		sum += len(n.Nodes[i].ReceiveRoute(1))
		// n.Nodes[i].PrintRelayState()
		// fmt.Printf("> Degree - Receving Count = %d\n", len(n.Nodes[i].Connections())-len(n.Nodes[i].ReceiveRoute(1)))
	}

	fmt.Printf("Total average duplicates for relay 1: %f\n", float64(sum)/float64(len(n.Nodes)-1)-1)
	n.PrintPropagationTree(1)
}
