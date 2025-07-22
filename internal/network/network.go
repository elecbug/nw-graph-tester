package network

import (
	"fmt"
	"math/rand"

	"github.com/elecbug/nw-graph-tester/internal/node"
)

type Network struct {
	Nodes []node.Node
}

func GenrateRandomNetwork(nodeCount int, edgeCount int, maxNodeDelay int64, maxLinkDelay int64) *Network {
	nodes := make([]node.Node, nodeCount)

	for i := 0; i < nodeCount; i++ {
		nodes[i] = node.Node{
			ID:          uint64(i),
			Connections: make(map[uint64]int64),
			NodeDelay:   rand.Int63n(maxNodeDelay),
		}
	}

	network := &Network{
		Nodes: nodes,
	}

	for i := 0; i < edgeCount; i++ {
		if !network.makeRandomConnection(maxLinkDelay) {
			i-- // Retry if connection could not be made
		}
	}

	return network
}

func (n *Network) makeRandomConnection(maxLinkDelay int64) bool {
	if len(n.Nodes) < 2 {
		return false // Not enough nodes to make a connection
	}

	nodeA := rand.Uint64() % uint64(len(n.Nodes))
	nodeB := rand.Uint64() % uint64(len(n.Nodes))

	for nodeB == nodeA {
		nodeB = rand.Uint64() % uint64(len(n.Nodes))
	}

	if _, ok := n.Nodes[nodeA].Connections[nodeB]; ok {
		return false // Connection already exists
	}

	linkDelay := rand.Int63n(maxLinkDelay)

	n.addConnection(nodeA, nodeB, linkDelay)

	return true
}

func (n *Network) addConnection(nodeA uint64, nodeB uint64, linkDelay int64) {
	if nodeA >= uint64(len(n.Nodes)) || nodeB >= uint64(len(n.Nodes)) {
		return // Invalid node IDs
	}

	n.Nodes[nodeA].Connections[nodeB] = linkDelay
	n.Nodes[nodeB].Connections[nodeA] = linkDelay
}

func (n *Network) Print() {
	for _, node := range n.Nodes {
		connections := "["
		i := 0
		for connID, delay := range node.Connections {
			connections += fmt.Sprintf("%2d(%2dms)", connID, delay)
			if i != len(node.Connections)-1 {
				connections += ", "
			}
			i++
		}
		connections += "]"
		fmt.Printf("Node ID: %2d, Delay: %2d, Connections: %d %s\n", node.ID, node.NodeDelay, len(node.Connections), connections)
	}
}
