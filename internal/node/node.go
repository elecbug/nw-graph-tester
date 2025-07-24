package node

import (
	"fmt"
	"sync"
	"time"

	"github.com/elecbug/p2p-broadcast-tester/internal/p2p"
)

type Node struct {
	id          p2p.NodeID
	delay       p2p.Delay
	connections map[*Node]p2p.Delay
	relayMap    map[p2p.MessageID]time.Time
	receiveMap  map[p2p.MessageID][]p2p.NodeID // For tracking duplicates
	mu          sync.RWMutex
}

func (n *Node) ID() p2p.NodeID {
	return n.id
}

func (n *Node) Delay() p2p.Delay {
	return n.delay
}

func (n *Node) Connections() map[*Node]p2p.Delay {
	return n.connections
}

func (n *Node) ReceiveRoute(messageID p2p.MessageID) []p2p.NodeID {
	n.mu.RLock()
	defer n.mu.RUnlock()

	return n.receiveMap[messageID]
}

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

func (n *Node) Broadcast(messageID p2p.MessageID, broadcastType p2p.BroadcastType) {

	switch broadcastType {
	case p2p.BasicPublish:
		go func() {
			n.mu.Lock()

			n.relayMap[messageID] = time.Now()
			n.receiveMap[messageID] = []p2p.NodeID{} // Reset duplicates for this relay

			n.mu.Unlock()

			time.Sleep(time.Duration(n.delay) * time.Millisecond)

			for conn, delay := range n.connections {
				if n.checkReceiving(messageID, conn) {
					continue
				}

				go func(conn *Node, delay p2p.Delay) {
					time.Sleep(time.Duration(delay) * time.Millisecond)

					conn.relay(messageID, n)
				}(conn, delay)
			}
		}()
	}
}

func (n *Node) relay(messageID p2p.MessageID, from *Node) {
	go func() {
		n.mu.Lock()

		if _, ok := n.relayMap[messageID]; ok {
			n.receiveMap[messageID] = append(n.receiveMap[messageID], from.id) // Track duplicate sender
			n.mu.Unlock()
			return
		} else {
			n.relayMap[messageID] = time.Now()
			n.receiveMap[messageID] = []p2p.NodeID{from.id} // Reset duplicates for this relay
			n.mu.Unlock()
		}

		time.Sleep(time.Duration(n.delay) * time.Millisecond)

		for conn, delay := range n.connections {
			if conn == from {
				continue // Skip excluded node
			}

			if n.checkReceiving(messageID, conn) {
				continue
			}

			go func(conn *Node, delay p2p.Delay) {
				time.Sleep(time.Duration(delay) * time.Millisecond)

				conn.relay(messageID, n)
			}(conn, delay)
		}
	}()
}

func (n *Node) checkReceiving(relayNumber p2p.MessageID, conn *Node) bool {
	n.mu.RLock()
	defer n.mu.RUnlock()

	for _, dupID := range n.receiveMap[relayNumber] {
		if dupID == conn.id {
			return true // Skip if this node has already relayed this message
		}
	}

	return false
}

func (n *Node) PrintRelayState() {
	n.mu.RLock()
	defer n.mu.RUnlock()

	for msgID, relayTime := range n.relayMap {
		recvs := "["
		for i, recvID := range n.receiveMap[msgID] {
			recvs += fmt.Sprintf("%d", recvID)
			if i < len(n.receiveMap[msgID])-1 {
				recvs += ", "
			}
		}
		recvs += "]"

		fmt.Printf("Node ID: %d, Message ID: %d, Relay Time: %s, Receive Route: %v\n", n.id, msgID, relayTime.Format(time.RFC3339), recvs)
	}
}
