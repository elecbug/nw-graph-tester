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

// Global mutex for thread-safe file writing
var mu sync.Mutex

// main function runs broadcast performance tests for different network configurations
func main() {
	// Test with different delay configurations (currently only d=0)
	for d := 0; d < 1; d++ {
		dCoef := 100  // Delay coefficient multiplier
		nCoef := 1000 // Node count coefficient multiplier
		wg := sync.WaitGroup{}

		// Test BasicPublish broadcast method with 10 different network sizes
		for i := 0; i < 10; i++ {
			wg.Add(1)

			go func(w *sync.WaitGroup, p p2p.BroadcastType, i, dCoef, nCoef int) {
				defer w.Done()

				fmt.Printf("Starting %s iteration %d\n", p.String(), i+1)
				Publish((i+1)*nCoef, p, (d+1)*dCoef)
			}(&wg, p2p.BroadcastType{Type: p2p.BasicPublish}, i, dCoef, nCoef)
		}

		wg.Wait()

		// Test WavePublish broadcast method with different levels (5, 10, 15, ..., 100)
		for p := 5; p <= 100; p += 5 {
			for i := 0; i < 10; i++ {
				wg.Add(1)

				go func(w *sync.WaitGroup, p p2p.BroadcastType, i, dCoef, nCoef int) {
					defer w.Done()

					fmt.Printf("Starting %s iteration %d\n", p.String(), i+1)
					Publish((i+1)*nCoef, p, (d+1)*dCoef)
				}(&wg, p2p.BroadcastType{Type: p2p.WavePublish, Level: p}, i, dCoef, nCoef)
			}

			wg.Wait()
		}
	}

	time.Sleep(time.Second * 1) // Wait for all goroutines to finish
}

// Publish creates a network and tests message broadcasting performance
// Parameters:
//   - nodeCount: number of nodes in the network
//   - broadcastType: the broadcast algorithm to test
//   - delay: maximum node processing delay
func Publish(nodeCount int, broadcastType p2p.BroadcastType, delay int) {
	meanDegree := 40 // Target average degree for network nodes

	// Generate a degree-limited network with specified parameters
	n := network.GenerateLimitDegreeNetwork(network.NetworkConfig{
		NodeCount:    nodeCount,
		DLow:         meanDegree - 2, // Minimum allowed degree
		D:            meanDegree,     // Target degree
		DHigh:        meanDegree + 2, // Maximum allowed degree
		MaxNodeDelay: p2p.Delay(delay),
		MaxLinkDelay: 1, // Fixed link delay
	})

	if n == nil {
		fmt.Println("Failed to generate network")
		return
	}
	// n.Print() // Uncomment to print network topology

	// Start broadcast from the first node (node 0)
	wg := &sync.WaitGroup{}
	n.Nodes[0].Broadcast(1, broadcastType, wg)

	wg.Wait() // Wait for broadcast to complete

	// Calculate broadcast performance metrics
	recvCount := 0     // Total number of message receptions (including duplicates)
	dontRecvCount := 0 // Number of nodes that didn't receive the message
	for i := range n.Nodes {
		recvCount += len(n.Nodes[i].ReceiveRoute(1))

		if len(n.Nodes[i].ReceiveRoute(1)) == 0 {
			dontRecvCount++
		}
	}

	recvTarget := len(n.Nodes) - 1 // Expected number of receivers (excluding sender)

	// Create network performance metric
	metric := p2p.NetworkMetric{
		NodeCount:     len(n.Nodes),
		Broadcast:     broadcastType.String(),
		Delay:         delay,
		AvgDegree:     float64(n.AvgDegree()),
		DuplicateRate: float64(recvCount)/float64(recvTarget-dontRecvCount+1) - 1, // Duplicate reception rate
		ReceivingRate: float64(recvTarget-dontRecvCount+1) / float64(recvTarget),  // Message delivery rate
	}

	// Serialize metric to JSON
	jsonData, err := json.Marshal(metric)
	jsonData = append(jsonData, '\n') // Add newline for better readability

	if err != nil {
		fmt.Printf("Error serializing network metric: %v\n", err)
		return
	}

	// Write metric to file in thread-safe manner
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
