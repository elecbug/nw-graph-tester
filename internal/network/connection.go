package network

import (
	"math/rand"

	"github.com/elecbug/p2p-broadcast-tester/internal/p2p"
)

// makeRandomConnection creates a random bidirectional connection between two random nodes
func (n *Network) makeRandomConnection(link p2p.Delay) bool {
	if len(n.Nodes) < 2 {
		return false // Not enough nodes to make a connection
	}

	// Select two different random nodes
	nodeA := rand.Uint64() % uint64(len(n.Nodes))
	nodeB := rand.Uint64() % uint64(len(n.Nodes))

	// Ensure nodeB is different from nodeA
	for nodeB == nodeA {
		nodeB = rand.Uint64() % uint64(len(n.Nodes))
	}

	// Check if connection already exists
	if _, ok := n.Nodes[nodeA].Connections()[&n.Nodes[nodeB]]; ok {
		return false // Connection already exists
	}

	// Generate random link delay within the specified range
	linkDelay := p2p.Delay(rand.Uint64() % uint64(link))

	n.AddBidirectConnection(nodeA, nodeB, linkDelay)

	return true
}

// AddBidirectConnection creates a bidirectional connection between two specified nodes
// Returns true if successful, false if connection already exists or invalid node IDs
func (n *Network) AddBidirectConnection(nodeA uint64, nodeB uint64, link p2p.Delay) bool {
	// Validate node IDs
	if nodeA >= uint64(len(n.Nodes)) || nodeB >= uint64(len(n.Nodes)) {
		return false // Invalid node IDs
	}

	// Check if connection already exists in either direction
	if _, ok1 := n.Nodes[nodeA].Connections()[&n.Nodes[nodeB]]; ok1 {
		return false // Connection already exists
	}
	if _, ok2 := n.Nodes[nodeB].Connections()[&n.Nodes[nodeA]]; ok2 {
		return false // Connection already exists
	}

	// Create bidirectional connection with same delay for both directions
	n.Nodes[nodeA].Connections()[&n.Nodes[nodeB]] = link
	n.Nodes[nodeB].Connections()[&n.Nodes[nodeA]] = link

	return true
}

// AddDirectConnection creates a unidirectional connection from nodeA to nodeB
// Used for creating directed graphs or asymmetric network topologies
func (n *Network) AddDirectConnection(nodeA uint64, nodeB uint64, link p2p.Delay) {
	// Validate node IDs
	if nodeA >= uint64(len(n.Nodes)) || nodeB >= uint64(len(n.Nodes)) {
		return // Invalid node IDs
	}

	// Check if connection already exists
	if _, ok := n.Nodes[nodeA].Connections()[&n.Nodes[nodeB]]; ok {
		return // Connection already exists
	}

	// Create unidirectional connection from nodeA to nodeB
	n.Nodes[nodeA].Connections()[&n.Nodes[nodeB]] = link
}

// RemoveConnection removes bidirectional connection between two specified nodes
// Deletes the connection in both directions to maintain network consistency
func (n *Network) RemoveConnection(nodeA uint64, nodeB uint64) {
	// Validate node IDs
	if nodeA >= uint64(len(n.Nodes)) || nodeB >= uint64(len(n.Nodes)) {
		return // Invalid node IDs
	}

	// Remove connection in both directions
	delete(n.Nodes[nodeA].Connections(), &n.Nodes[nodeB])
	delete(n.Nodes[nodeB].Connections(), &n.Nodes[nodeA])
}

// AvgDegree calculates the average degree (number of connections) across all nodes
// Returns 0 if the network has no nodes
func (n *Network) AvgDegree() float64 {
	if len(n.Nodes) == 0 {
		return 0
	}

	// Sum up the degree of all nodes
	totalDegree := 0
	for i := range n.Nodes {
		totalDegree += len(n.Nodes[i].Connections())
	}

	// Calculate and return average degree
	return float64(totalDegree) / float64(len(n.Nodes))
}

// delay generates a random delay value within the specified range [min, max]
// Ensures min <= max by swapping if necessary
func delay(min, max p2p.Delay) p2p.Delay {
	if min > max {
		min, max = max, min // Ensure min is less than or equal to max
	}

	// Generate random value in range [min, max] inclusive
	return p2p.Delay(rand.Uint64()%(uint64(max)-uint64(min)+1) + uint64(min))
}
