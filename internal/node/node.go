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

func (n *Node) Broadcast(messageID MessageID) {
	go func() {
		n.Mu.Lock()

		n.RelayMap[messageID] = time.Now()
		n.DuplicateMap[messageID] = []NodeID{} // Reset duplicates for this relay

		n.Mu.Unlock()

		time.Sleep(time.Duration(n.NodeDelay) * time.Millisecond)

		for conn, delay := range n.Connections {
			if n.check(messageID, conn) {
				continue
			}

			go func(conn *Node, delay int64) {
				time.Sleep(time.Duration(delay) * time.Millisecond)

				conn.relay(messageID, n)
			}(conn, delay)
		}
	}()
}

func (n *Node) relay(messageID MessageID, from *Node) {
	go func() {
		n.Mu.Lock()

		if _, ok := n.RelayMap[messageID]; ok {
			n.DuplicateMap[messageID] = append(n.DuplicateMap[messageID], from.ID) // Track duplicate sender
			n.Mu.Unlock()
			return
		} else {
			n.RelayMap[messageID] = time.Now()
			n.DuplicateMap[messageID] = []NodeID{} // Reset duplicates for this relay
			n.Mu.Unlock()
		}

		time.Sleep(time.Duration(n.NodeDelay) * time.Millisecond)

		for conn, delay := range n.Connections {
			if conn == from {
				continue // Skip excluded node
			}

			if n.check(messageID, conn) {
				continue
			}

			go func(conn *Node, delay int64) {
				time.Sleep(time.Duration(delay) * time.Millisecond)

				conn.relay(messageID, n)
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
