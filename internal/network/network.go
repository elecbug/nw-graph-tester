package network

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/elecbug/nw-graph-tester/internal/node"
)

type Network struct {
	Nodes      []node.Node
	delayLog10 int64
}

func GenerateRandomNetwork(nodeCount int, edgeCount int, maxNodeDelay int64, maxLinkDelay int64) *Network {
	nodes := make([]node.Node, nodeCount)
	delayLog10 := getDelayLog10(max(maxNodeDelay, maxLinkDelay))

	for i := 0; i < nodeCount; i++ {
		nodes[i] = node.Node{
			ID:           node.NodeID(i),
			Connections:  make(map[*node.Node]int64),
			NodeDelay:    rand.Int63n(maxNodeDelay),
			RelayMap:     make(map[node.MessageID]time.Time),
			DuplicateMap: make(map[node.MessageID][]node.NodeID),
		}
	}

	network := &Network{
		Nodes:      nodes,
		delayLog10: delayLog10,
	}

	for i := 0; i < edgeCount; i++ {
		if !network.makeRandomConnection(maxLinkDelay) {
			i-- // Retry if connection could not be made
		}
	}

	return network
}

func GenerateGossipSubNetwork(nodeCount int, d, dLow, dHigh int, maxNodeDelay int64, maxLinkDelay int64) *Network {
	nodes := make([]node.Node, nodeCount)
	delayLog10 := getDelayLog10(max(maxNodeDelay, maxLinkDelay))

	for i := 0; i < nodeCount; i++ {
		nodes[i] = node.Node{
			ID:           node.NodeID(i),
			Connections:  make(map[*node.Node]int64),
			NodeDelay:    rand.Int63n(maxNodeDelay),
			RelayMap:     make(map[node.MessageID]time.Time),
			DuplicateMap: make(map[node.MessageID][]node.NodeID),
		}
	}

	network := &Network{
		Nodes:      nodes,
		delayLog10: delayLog10,
	}

	for i := 0; i < nodeCount; i++ {
		d := 0

		for {
			for j := 0; j < dHigh; j++ {
				if network.makeRandomConnection(maxLinkDelay) {
					d++
				}

				if d > dHigh {
					break
				}
			}

			if d > dLow {
				break
			}
		}
	}

	for re := 0; re < 100; re++ {
		for i := 0; i < nodeCount; i++ {
			if len(network.Nodes[i].Connections) > dHigh {
				for j := 0; j < len(network.Nodes[i].Connections)-d; j++ {
					network.removeConnection(uint64(i), rand.Uint64()%uint64(len(network.Nodes)))
				}
			}

			if len(network.Nodes[i].Connections) < dLow {
				for j := 0; j < d-len(network.Nodes[i].Connections); j++ {
					if !network.makeRandomConnection(maxLinkDelay) {
						j-- // Retry if connection could not be made
					}
				}
			}
		}
	}

	return network
}

func getDelayLog10(delay int64) int64 {
	if delay <= 0 {
		return 0
	}

	log10 := int64(0)
	for delay > 9 {
		delay /= 10
		log10++
	}

	return log10
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

	if _, ok := n.Nodes[nodeA].Connections[&n.Nodes[nodeB]]; ok {
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

	n.Nodes[nodeA].Connections[&n.Nodes[nodeB]] = linkDelay
	n.Nodes[nodeB].Connections[&n.Nodes[nodeA]] = linkDelay
}

func (n *Network) removeConnection(nodeA uint64, nodeB uint64) {
	if nodeA >= uint64(len(n.Nodes)) || nodeB >= uint64(len(n.Nodes)) {
		return // Invalid node IDs
	}

	delete(n.Nodes[nodeA].Connections, &n.Nodes[nodeB])
	delete(n.Nodes[nodeB].Connections, &n.Nodes[nodeA])
}

func (n *Network) Print() {
	for i := range n.Nodes {
		node := &n.Nodes[i]

		connections := "["
		i := 0
		for connNode, delay := range node.Connections {
			connections += fmt.Sprintf("%d(%dms)", connNode.ID, delay)
			if i != len(node.Connections)-1 {
				connections += ", "
			}
			i++
		}
		connections += "]"

		fmt.Printf("Node ID: %d, Delay: %d, Connections: %d %s\n", node.ID, node.NodeDelay, len(node.Connections), connections)
	}
}
