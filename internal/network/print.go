package network

import (
	"fmt"
	"strings"

	"github.com/elecbug/p2p-broadcast-tester/internal/p2p"
)

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

func (n *Network) PrintPropagationGraph(mid p2p.MessageID) {
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
