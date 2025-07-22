package node

import (
	"fmt"
	"sync"
	"time"
)

type Node struct {
	ID           uint64
	Connections  map[*Node]int64
	NodeDelay    int64
	RelayMap     map[uint64]time.Time
	DuplicateMap map[uint64]uint64 // For tracking duplicates
	Mu           sync.RWMutex
}

func (n *Node) Relay(relayNumber uint64, from *Node) {
	go func() {
		n.Mu.Lock()

		if _, ok := n.RelayMap[relayNumber]; ok {
			n.DuplicateMap[relayNumber] += 1 // Increment duplicate count
			n.Mu.Unlock()
			return
		} else {
			n.RelayMap[relayNumber] = time.Now()
			n.Mu.Unlock()
		}

		time.Sleep(time.Duration(n.NodeDelay) * time.Millisecond)

		for conn, delay := range n.Connections {
			if conn == from {
				continue // Skip excluded node
			}

			go func(conn *Node, delay int64) {
				time.Sleep(time.Duration(delay) * time.Millisecond)

				conn.Relay(relayNumber, n)
			}(conn, delay)
		}
	}()
}

func (n *Node) PrintRelayState() {
	n.Mu.RLock()
	defer n.Mu.RUnlock()

	fmt.Printf("Node ID: %d, RelayMap: %v, DuplicateMap: %v\n", n.ID, n.RelayMap, n.DuplicateMap)
}
