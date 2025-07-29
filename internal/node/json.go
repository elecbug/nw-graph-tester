package node

import (
	"encoding/json"
	"log"
	"time"

	"github.com/elecbug/p2p-broadcast-tester/internal/p2p"
)

// ToJson converts the node to a JSON string representation
func (n *Node) ToJson() (string, error) {
	// Create a serializable version of the node with node pointers converted to IDs
	nodeM := struct {
		ID          p2p.NodeID                     `json:"id"`
		Delay       p2p.Delay                      `json:"delay"`
		Connections map[p2p.NodeID]p2p.Delay       `json:"connections"`
		RelayMap    map[p2p.MessageID]time.Time    `json:"relay_map"`
		ReceiveMap  map[p2p.MessageID][]p2p.NodeID `json:"receive_map"`
	}{
		ID:          n.id,
		Delay:       n.delay,
		Connections: connectionConv(n.connections), // Convert node pointers to node IDs
		RelayMap:    n.relayMap,
		ReceiveMap:  n.receiveMap,
	}

	// Marshal the struct to JSON
	data, err := json.Marshal(nodeM)
	if err != nil {
		log.Printf("Error serializing node to JSON: %v", err)
		return "", err
	}

	return string(data), nil
}

// connectionConv converts a map of node pointers to a map of node IDs
// This is necessary because node pointers cannot be directly serialized to JSON
func connectionConv(connections map[*Node]p2p.Delay) map[p2p.NodeID]p2p.Delay {
	result := make(map[p2p.NodeID]p2p.Delay)

	// Convert each node pointer to its ID while preserving the delay value
	for node, delay := range connections {
		result[node.id] = delay
	}

	return result
}
