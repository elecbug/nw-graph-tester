package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"sync"
	"time"

	"github.com/elecbug/p2p-broadcast-tester/internal/network"
	"github.com/elecbug/p2p-broadcast-tester/internal/p2p"
)

var mu sync.Mutex

func main() {
	for d := 0; d < 10; d++ {
		dCoef := 1000
		nCoef := 1000
		w := sync.WaitGroup{}

		for _, p := range p2p.AllBroadcastTypes() {
			for i := 0; i < 10; i++ {
				w.Add(1)
				go func(w *sync.WaitGroup, p p2p.BroadcastType, i, dCoef, nCoef int) {
					defer w.Done()

					fmt.Printf("Starting %s iteration %d\n", p.String(), i+1)
					Publish((i+1)*nCoef, p, (d+1)*dCoef)
					time.Sleep(time.Second * 2)
				}(&w, p, i, dCoef, nCoef)
			}
			w.Wait()
		}
	}

	time.Sleep(time.Second * 60) // Wait for all goroutines to finish
}

func Publish(nodeCount int, broadcastType p2p.BroadcastType, delay int) {
	n := network.GenerateGossipSubNetwork(network.NetworkConfig{
		NodeCount:    nodeCount,
		DLow:         18,
		D:            20,
		DHigh:        22,
		MaxNodeDelay: p2p.Delay(delay),
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

	metric := NetworkMetric{
		NodeCount:     len(n.Nodes),
		Broadcast:     broadcastType.String(),
		Delay:         delay,
		DuplicateRate: float64(recvCount)/float64(recvTarget-dontRecvCount+1) - 1,
		ReceivingRate: float64(recvTarget-dontRecvCount+1) / float64(recvTarget),
	}

	jsonData, err := json.Marshal(metric)
	jsonData = append(jsonData, '\n') // Add newline for better readability

	if err != nil {
		fmt.Printf("Error serializing network metric: %v\n", err)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	file, err := os.OpenFile("results/network_metric.jsonl", os.O_CREATE|os.O_WRONLY|os.O_APPEND, fs.ModePerm)

	if err != nil {
		fmt.Printf("Error opening network_metric.jsonl: %v\n", err)
		return
	}
	defer file.Close()

	if _, err := file.Write(jsonData); err != nil {
		fmt.Printf("Error writing to network_metric.jsonl: %v\n", err)
		return
	}
}

type NetworkMetric struct {
	NodeCount     int     `json:"node_count"`
	Broadcast     string  `json:"broadcast"`
	Delay         int     `json:"delay"`
	DuplicateRate float64 `json:"duplicate_rate"`
	ReceivingRate float64 `json:"receiving_rate"`
}
