package node

import (
	"sync"
	"time"

	"github.com/elecbug/p2p-broadcast-tester/internal/p2p"
)

// Node represents a single node in the P2P network
type Node struct {
	id          p2p.NodeID                     // Unique identifier for this node
	delay       p2p.Delay                      // Network delay for this node
	relayMap    map[p2p.MessageID]time.Time    // For tracking relay times
	receiveMap  map[p2p.MessageID][]p2p.NodeID // For tracking duplicates
	connections map[*Node]p2p.Delay            // Map of connected nodes and their delays
	mu          sync.RWMutex                   // Mutex for thread-safe access
}

// NewNode creates a new node with the given ID and delay
func NewNode(id p2p.NodeID, delay p2p.Delay) *Node {
	return &Node{
		id:          id,
		connections: make(map[*Node]p2p.Delay),
		delay:       delay,
		relayMap:    make(map[p2p.MessageID]time.Time),
		receiveMap:  make(map[p2p.MessageID][]p2p.NodeID),
		mu:          sync.RWMutex{},
	}
}

// ID returns the unique identifier of this node
func (n *Node) ID() p2p.NodeID {
	return n.id
}

// Delay returns the network delay of this node
func (n *Node) Delay() p2p.Delay {
	return n.delay
}

// Connections returns the map of connected nodes and their delays
func (n *Node) Connections() map[*Node]p2p.Delay {
	return n.connections
}

// RelayTime returns the relay time for a specific message ID and whether it exists
func (n *Node) RelayTime(messageID p2p.MessageID) (time.Time, bool) {
	n.mu.RLock()
	defer n.mu.RUnlock()

	relayTime, exists := n.relayMap[messageID]
	return relayTime, exists
}

// ReceiveRoute returns the list of node IDs from which a message was received
func (n *Node) ReceiveRoute(messageID p2p.MessageID) []p2p.NodeID {
	n.mu.RLock()
	defer n.mu.RUnlock()

	return n.receiveMap[messageID]
}
