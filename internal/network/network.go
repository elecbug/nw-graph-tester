package network

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/elecbug/p2p-broadcast-tester/internal/node"
	"github.com/elecbug/p2p-broadcast-tester/internal/p2p"
)

type Network struct {
	Nodes      []node.Node
	delayLog10 int64
}

type NetworkConfig struct {
	NodeCount    int
	EdgeCount    int
	MaxNodeDelay p2p.Delay
	MaxLinkDelay p2p.Delay
	D            int // Gossip sub parameters
	DLow         int // Gossip sub parameters
	DHigh        int // Gossip sub parameters
}

func GenerateRandomNetwork(config NetworkConfig) *Network {
	nodes := make([]node.Node, config.NodeCount)
	delayLog10 := getDelayLog10(max(config.MaxNodeDelay, config.MaxLinkDelay))

	for i := 0; i < config.NodeCount; i++ {
		nodes[i] = *node.NewNode(p2p.NodeID(i), p2p.Delay(rand.Uint64()%uint64(config.MaxNodeDelay)))
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
		delay := p2p.Delay(rand.Uint64() % uint64(config.MaxNodeDelay))
		nodes[i] = *node.NewNode(p2p.NodeID(i), delay)
	}

	network := &Network{
		Nodes:      nodes,
		delayLog10: delayLog10,
	}

	for re := 0; re < config.NodeCount; re++ {
		flag := false

		if re%100 == 0 {
			fmt.Printf("Generating connections for node %d/%d\n", re, config.NodeCount)
		}

		for i := 0; i < config.NodeCount; i++ {
			if len(network.Nodes[i].Connections()) < config.DLow {
				for j := 0; j < config.D-len(network.Nodes[i].Connections()); j++ {
					target := rand.Uint64() % uint64(len(network.Nodes))
					delay := p2p.Delay(rand.Uint64() % uint64(config.MaxLinkDelay))

					if !network.AddBidirectConnection(uint64(i), target, delay) {
						j-- // Retry if connection could not be made
						flag = true
					}
				}
			}

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

func getDelayLog10(delay p2p.Delay) int64 {
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

func (n *Network) makeRandomConnection(maxLinkDelay p2p.Delay) bool {
	if len(n.Nodes) < 2 {
		return false // Not enough nodes to make a connection
	}

	nodeA := rand.Uint64() % uint64(len(n.Nodes))
	nodeB := rand.Uint64() % uint64(len(n.Nodes))

	for nodeB == nodeA {
		nodeB = rand.Uint64() % uint64(len(n.Nodes))
	}

	if _, ok := n.Nodes[nodeA].Connections()[&n.Nodes[nodeB]]; ok {
		return false // Connection already exists
	}

	linkDelay := p2p.Delay(rand.Uint64() % uint64(maxLinkDelay))

	n.AddBidirectConnection(nodeA, nodeB, linkDelay)

	return true
}

func (n *Network) AddBidirectConnection(nodeA uint64, nodeB uint64, linkDelay p2p.Delay) bool {
	if nodeA >= uint64(len(n.Nodes)) || nodeB >= uint64(len(n.Nodes)) {
		return false // Invalid node IDs
	}

	if _, ok1 := n.Nodes[nodeA].Connections()[&n.Nodes[nodeB]]; ok1 {
		return false // Connection already exists
	}
	if _, ok2 := n.Nodes[nodeB].Connections()[&n.Nodes[nodeA]]; ok2 {
		return false // Connection already exists
	}

	n.Nodes[nodeA].Connections()[&n.Nodes[nodeB]] = linkDelay
	n.Nodes[nodeB].Connections()[&n.Nodes[nodeA]] = linkDelay

	return true
}

func (n *Network) AddConnection(nodeA uint64, nodeB uint64, linkDelay p2p.Delay) {
	if nodeA >= uint64(len(n.Nodes)) || nodeB >= uint64(len(n.Nodes)) {
		return // Invalid node IDs
	}

	if _, ok := n.Nodes[nodeA].Connections()[&n.Nodes[nodeB]]; ok {
		return // Connection already exists
	}

	n.Nodes[nodeA].Connections()[&n.Nodes[nodeB]] = linkDelay
}

func (n *Network) RemoveConnection(nodeA uint64, nodeB uint64) {
	if nodeA >= uint64(len(n.Nodes)) || nodeB >= uint64(len(n.Nodes)) {
		return // Invalid node IDs
	}

	delete(n.Nodes[nodeA].Connections(), &n.Nodes[nodeB])
	delete(n.Nodes[nodeB].Connections(), &n.Nodes[nodeA])
}

func (n *Network) Print() {
	for i := range n.Nodes {
		node := &n.Nodes[i]

		connections := "["
		i := 0
		for connNode, delay := range node.Connections() {
			connections += fmt.Sprintf("%d(%dms)", connNode.ID(), delay)
			if i != len(node.Connections())-1 {
				connections += ", "
			}
			i++
		}
		connections += "]"

		fmt.Printf("Node ID: %d, Delay: %d, Connections: %d %s\n", node.ID(), node.Delay(), len(node.Connections()), connections)
	}
}

func (n *Network) PropagationGraph(mid p2p.MessageID) {
	propas := make(map[p2p.NodeID]p2p.NodeID)

	for i := 0; i < len(n.Nodes); i++ {
		recvs := n.Nodes[i].ReceiveRoute(mid)

		if len(recvs) == 0 {
			continue // No propagation for this message
		}

		propas[n.Nodes[i].ID()] = recvs[0]
	}

	fmt.Println("Propagation Graph:")
	for sender, receiver := range propas {
		fmt.Printf("Node %d -> Node %d\n", sender, receiver)
	}
}

func (n *Network) PrintPropagationTree(mid p2p.MessageID) {
	propas := make(map[p2p.NodeID]p2p.NodeID) // child -> parent

	for i := 0; i < len(n.Nodes); i++ {
		recvs := n.Nodes[i].ReceiveRoute(mid)
		if len(recvs) == 0 {
			continue
		}
		propas[n.Nodes[i].ID()] = recvs[0]
	}

	// 1. Build parent -> children map
	children := make(map[p2p.NodeID][]p2p.NodeID)
	var root p2p.NodeID

	isChild := make(map[p2p.NodeID]bool)
	for child, parent := range propas {
		children[parent] = append(children[parent], child)
		isChild[child] = true
	}

	// 2. Find root (who is never a child)
	for node := range propas {
		if !isChild[node] {
			root = node
			break
		}
	}

	// 3. DFS 출력
	var dfs func(node p2p.NodeID, depth int)
	dfs = func(node p2p.NodeID, depth int) {
		fmt.Printf("%sNode %d\n", strings.Repeat("  ", depth), node)
		for _, child := range children[node] {
			dfs(child, depth+1)
		}
	}

	fmt.Println("Propagation Tree:")
	dfs(root, 0)
}
