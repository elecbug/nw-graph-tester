package node

import (
	"fmt"
	"sync"
	"time"
)

type NodeID uint64
type MessageID uint64

type Node struct {
	ID           NodeID
	Connections  map[*Node]int64
	NodeDelay    int64
	RelayMap     map[MessageID]time.Time
	DuplicateMap map[MessageID][]NodeID // For tracking duplicates
	Mu           sync.RWMutex
}

func (n *Node) Relay(relayNumber MessageID, from *Node) {
	go func() {
		n.Mu.Lock()

		if _, ok := n.RelayMap[relayNumber]; ok {
			n.DuplicateMap[relayNumber] = append(n.DuplicateMap[relayNumber], from.ID) // Track duplicate sender
			n.Mu.Unlock()
			return
		} else {
			n.RelayMap[relayNumber] = time.Now()
			n.DuplicateMap[relayNumber] = []NodeID{} // Reset duplicates for this relay
			n.Mu.Unlock()
		}

		time.Sleep(time.Duration(n.NodeDelay) * time.Millisecond)

		for conn, delay := range n.Connections {
			if conn == from {
				continue // Skip excluded node
			}

			if n.check(relayNumber, conn) {
				continue
			}

			go func(conn *Node, delay int64) {
				time.Sleep(time.Duration(delay) * time.Millisecond)

				conn.Relay(relayNumber, n)
			}(conn, delay)
		}
	}()
}

func (n *Node) check(relayNumber MessageID, conn *Node) bool {
	n.Mu.RLock()
	defer n.Mu.RUnlock()

	for _, dupID := range n.DuplicateMap[relayNumber] {
		if dupID == conn.ID {
			return true // Skip if this node has already relayed this message
		}
	}

	return false
}

func (n *Node) PrintRelayState() {
	n.Mu.RLock()
	defer n.Mu.RUnlock()

	fmt.Printf("Node ID: %d, RelayMap: %v, DuplicateMap: %v\n", n.ID, n.RelayMap, n.DuplicateMap)
}
