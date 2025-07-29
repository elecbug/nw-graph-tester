package network

import (
	"math/rand"

	"github.com/elecbug/p2p-broadcast-tester/internal/node"
	"github.com/elecbug/p2p-broadcast-tester/internal/p2p"
)

// Network represents a P2P network containing multiple nodes
type Network struct {
	Nodes []node.Node // List of all nodes in the network
}

// NetworkConfig contains configuration parameters for network generation
type NetworkConfig struct {
	NodeCount    int       // Total number of nodes in the network
	MinNodeDelay p2p.Delay // Minimum processing delay for nodes
	MaxNodeDelay p2p.Delay // Maximum processing delay for nodes
	MinLinkDelay p2p.Delay // Minimum transmission delay for links
	MaxLinkDelay p2p.Delay // Maximum transmission delay for links
	EdgeCount    int       // Total number of edges to create (for random network)
	D            int       // Target degree for each node (for degree-limited network)
	DLow         int       // Minimum allowed degree for nodes
	DHigh        int       // Maximum allowed degree for nodes
}

// GenerateRandomNetwork creates a network with randomly distributed connections
func GenerateRandomNetwork(config NetworkConfig) *Network {
	// Create nodes with random delays within specified range
	nodes := make([]node.Node, config.NodeCount)

	for i := 0; i < config.NodeCount; i++ {
		nodes[i] = *node.NewNode(p2p.NodeID(i), delay(config.MinNodeDelay, config.MaxNodeDelay))
	}

	network := &Network{
		Nodes: nodes,
	}

	// Create random connections between nodes
	for i := 0; i < config.EdgeCount; i++ {
		if !network.makeRandomConnection(delay(config.MinLinkDelay, config.MaxLinkDelay)) {
			i-- // Retry if connection could not be made
		}
	}

	return network
}

// GenerateLimitDegreeNetwork creates a network where each node's degree is controlled
// to be within specified bounds (DLow <= degree <= DHigh)
func GenerateLimitDegreeNetwork(config NetworkConfig) *Network {
	// Create nodes with random delays within specified range
	nodes := make([]node.Node, config.NodeCount)

	for i := 0; i < config.NodeCount; i++ {
		nodes[i] = *node.NewNode(p2p.NodeID(i), delay(config.MinNodeDelay, config.MaxNodeDelay))
	}

	network := &Network{
		Nodes: nodes,
	}

	// Iteratively adjust node degrees to meet constraints
	for re := 0; re < config.NodeCount; re++ {
		flag := false

		// Progress indicator (commented out for performance)
		// if re%100 == 0 {
		// 	fmt.Printf("Generating connections for node %d/%d\n", re, config.NodeCount)
		// }

		// Check each node's degree and adjust if necessary
		for i := 0; i < config.NodeCount; i++ {
			// Add connections if degree is below minimum threshold
			if len(network.Nodes[i].Connections()) < config.DLow {
				for j := 0; j < config.D-len(network.Nodes[i].Connections()); j++ {
					target := rand.Uint64() % uint64(len(network.Nodes))

					if !network.AddBidirectConnection(uint64(i), target, delay(config.MinLinkDelay, config.MaxLinkDelay)) {
						j-- // Retry if connection could not be made
						flag = true
					}
				}
			}

			// Remove connections if degree is above maximum threshold
			if len(network.Nodes[i].Connections()) > config.DHigh {
				for j := 0; j < len(network.Nodes[i].Connections())-config.D; j++ {
					target := rand.Uint64() % uint64(len(network.Nodes))

					network.RemoveConnection(uint64(i), target)
					flag = true
				}
			}
		}

		if !flag {
			break // No changes made, exit early
		}
	}

	return network
}
