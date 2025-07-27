package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"time"

	"github.com/elecbug/p2p-broadcast-tester/internal/network"
	"github.com/elecbug/p2p-broadcast-tester/internal/p2p"
)

func main() {
	for i := 0; i < 10; i++ {
		fmt.Printf("Starting BasicPublish iteration %d\n", i+1)
		Publish((i+1)*1000, p2p.BasicPublish)
		time.Sleep(time.Second * 2)
	}

	for i := 0; i < 10; i++ {
		fmt.Printf("Starting WavePublish iteration %d\n", i+1)
		Publish((i+1)*1000, p2p.WavePublish)
		time.Sleep(time.Second * 2)
	}
}

func Publish(nodeCount int, broadcastType p2p.BroadcastType) {
	n := network.GenerateGossipSubNetwork(network.NetworkConfig{
		NodeCount:    nodeCount,
		DLow:         8,
		D:            10,
		DHigh:        12,
		MaxNodeDelay: 1000,
		MaxLinkDelay: 1,
	})

	if n == nil {
		fmt.Println("Failed to generate network")
		return
	}
	// n.Print()

	n.Nodes[0].Broadcast(1, broadcastType)

	time.Sleep(time.Second * 10)

	recvCount := 0
	dontRecvCount := 0
	for i := range n.Nodes {
		recvCount += len(n.Nodes[i].ReceiveRoute(1))

		if len(n.Nodes[i].ReceiveRoute(1)) == 0 {
			dontRecvCount++
		}
	}

	recvTarget := len(n.Nodes) - 1
	fmt.Printf("duplicate: %f\n", float64(recvCount)/float64(recvTarget)-1)
	fmt.Printf("Total nodes not receiving relay 1: %f\n", float64(recvTarget-dontRecvCount)/float64(recvTarget))

	metric := NetworkMetric{
		NodeCount:     len(n.Nodes),
		Broadcast:     "WavePublish",
		DuplicateRate: float64(recvCount)/float64(recvTarget) - 1,
		ReceivingRate: float64(recvTarget-dontRecvCount) / float64(recvTarget),
	}

	jsonData, err := json.Marshal(metric)

	if err != nil {
		fmt.Printf("Error serializing network metric: %v\n", err)
		return
	}

	file, err := os.OpenFile("network_metric.json", os.O_CREATE|os.O_WRONLY|os.O_APPEND, fs.ModePerm)

	if err != nil {
		fmt.Printf("Error opening network_metric.json: %v\n", err)
		return
	}
	defer file.Close()

	if _, err := file.Write(jsonData); err != nil {
		fmt.Printf("Error writing to network_metric.json: %v\n", err)
		return
	}
}

type NetworkMetric struct {
	NodeCount     int     `json:"node_count"`
	Broadcast     string  `json:"broadcast"`
	DuplicateRate float64 `json:"duplicate_rate"`
	ReceivingRate float64 `json:"receiving_rate"`
}
