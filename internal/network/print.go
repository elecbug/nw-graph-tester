package network

import (
	"fmt"
	"strings"

	"github.com/elecbug/p2p-broadcast-tester/internal/p2p"
)

// Print displays detailed information about all nodes in the network
// Shows node ID, delay, connection count, and list of connected nodes with their delays
func (n *Network) Print() {
	for i := range n.Nodes {
		node := &n.Nodes[i]

		// Build connection string showing connected nodes and their delays
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

// PrintPropagationRoute displays the propagation route for a specific message
// Shows which node received the message from which other node (first sender only)
func (n *Network) PrintPropagationRoute(mid p2p.MessageID) {
	// Map to store propagation relationships (receiver -> sender)
	propas := make(map[p2p.NodeID]p2p.NodeID)

	// Extract propagation information from each node
	for i := 0; i < len(n.Nodes); i++ {
		recvs := n.Nodes[i].ReceiveRoute(mid)

		if len(recvs) == 0 {
			continue // No propagation for this message
		}

		// Store the first sender (primary propagation path)
		propas[n.Nodes[i].ID()] = recvs[0]
	}

	// Display the propagation route
	fmt.Println("Propagation Route:")
	for sender, receiver := range propas {
		fmt.Printf("Node %d -> Node %d\n", sender, receiver)
	}
}

// PrintPropagationTree displays the message propagation as a tree structure
// Shows the hierarchical relationship of how the message spread through the network
func (n *Network) PrintPropagationTree(mid p2p.MessageID) {
	propas := make(map[p2p.NodeID]p2p.NodeID) // child -> parent mapping

	// Extract propagation relationships from all nodes
	for i := 0; i < len(n.Nodes); i++ {
		recvs := n.Nodes[i].ReceiveRoute(mid)
		if len(recvs) == 0 {
			continue
		}
		// Map each node to its parent (first sender)
		propas[n.Nodes[i].ID()] = recvs[0]
	}

	// 1. Build parent -> children map for tree structure
	children := make(map[p2p.NodeID][]p2p.NodeID)
	var root p2p.NodeID

	// Track which nodes are children (have parents)
	isChild := make(map[p2p.NodeID]bool)
	for child, parent := range propas {
		children[parent] = append(children[parent], child)
		isChild[child] = true
	}

	// 2. Find root node (the one that is never a child)
	for node := range propas {
		if !isChild[node] {
			root = node
			break
		}
	}

	// 3. DFS traversal to print tree structure with indentation
	var dfs func(node p2p.NodeID, depth int)
	dfs = func(node p2p.NodeID, depth int) {
		// Print current node with appropriate indentation
		fmt.Printf("%sNode %d\n", strings.Repeat("  ", depth), node)
		// Recursively print all children
		for _, child := range children[node] {
			dfs(child, depth+1)
		}
	}

	fmt.Println("Propagation Tree:")
	dfs(root, 0)
}
