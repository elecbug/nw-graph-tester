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

type NetworkConfig struct {
	NodeCount    int
	EdgeCount    int
	MaxNodeDelay int64
	MaxLinkDelay int64
	D            int // Gossip sub parameters
	DLow         int // Gossip sub parameters
	DHigh        int // Gossip sub parameters
}

func GenerateRandomNetwork(config NetworkConfig) *Network {
	nodes := make([]node.Node, config.NodeCount)
	delayLog10 := getDelayLog10(max(config.MaxNodeDelay, config.MaxLinkDelay))

	for i := 0; i < config.NodeCount; i++ {
		nodes[i] = node.Node{
			ID:           node.NodeID(i),
			Connections:  make(map[*node.Node]int64),
			NodeDelay:    rand.Int63n(config.MaxNodeDelay),
			RelayMap:     make(map[node.MessageID]time.Time),
			DuplicateMap: make(map[node.MessageID][]node.NodeID),
		}
	}

	network := &Network{
		Nodes:      nodes,
		delayLog10: delayLog10,
	}

	for i := 0; i < config.EdgeCount; i++ {
		if !network.makeRandomConnection(config.MaxLinkDelay) {
			i-- // Retry if connection could not be made
		}
	}

	return network
}

func GenerateGossipSubNetwork(config NetworkConfig) *Network {
	nodes := make([]node.Node, config.NodeCount)
	delayLog10 := getDelayLog10(max(config.MaxNodeDelay, config.MaxLinkDelay))

	for i := 0; i < config.NodeCount; i++ {
		nodes[i] = node.Node{
			ID:           node.NodeID(i),
			Connections:  make(map[*node.Node]int64),
			NodeDelay:    rand.Int63n(config.MaxNodeDelay),
			RelayMap:     make(map[node.MessageID]time.Time),
			DuplicateMap: make(map[node.MessageID][]node.NodeID),
		}
	}

	network := &Network{
		Nodes:      nodes,
		delayLog10: delayLog10,
	}

	for i := 0; i < config.NodeCount; i++ {
		d := 0

		for {
			for j := 0; j < config.DHigh; j++ {
				if network.makeRandomConnection(config.MaxLinkDelay) {
					d++
				}

				if d > config.DHigh {
					break
				}
			}

			if d > config.DLow {
				break
			}
		}
	}

	for re := 0; re < 100; re++ {
		for i := 0; i < config.NodeCount; i++ {
			if len(network.Nodes[i].Connections) > config.DHigh {
				for j := 0; j < len(network.Nodes[i].Connections)-config.D; j++ {
					network.removeConnection(uint64(i), rand.Uint64()%uint64(len(network.Nodes)))
				}
			}

			if len(network.Nodes[i].Connections) < config.DLow {
				for j := 0; j < config.D-len(network.Nodes[i].Connections); j++ {
					if !network.makeRandomConnection(config.MaxLinkDelay) {
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
