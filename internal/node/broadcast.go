package node

import (
	"math/rand"
	"sync"
	"time"

	"github.com/elecbug/p2p-broadcast-tester/internal/p2p"
)

// Broadcast initiates a message broadcast using the specified broadcast type
func (n *Node) Broadcast(messageID p2p.MessageID, broadcastType p2p.BroadcastType, wg *sync.WaitGroup) {
	defer wg.Done()
	wg.Add(1)

	switch broadcastType.Type {
	case p2p.BasicPublish:
		// Basic flooding-based broadcast: send to all connected nodes
		n.mu.Lock()

		n.relayMap[messageID] = time.Now()
		n.receiveMap[messageID] = []p2p.NodeID{} // Reset duplicates for this relay

		n.mu.Unlock()

		// Simulate node processing delay
		time.Sleep(time.Duration(n.delay) * time.Millisecond)

		// Send message to all connected nodes
		for conn, delay := range n.connections {
			if n.checkReceiving(messageID, conn) {
				continue
			}

			wg.Add(1)

			go func(conn *Node, delay p2p.Delay) {
				defer wg.Done()

				// Simulate network transmission delay
				time.Sleep(time.Duration(delay) * time.Millisecond)

				conn.relayBasic(messageID, n, wg)
			}(conn, delay)
		}
	case p2p.WavePublish:
		// Wave-based broadcast with controlled propagation using level parameter
		coef := float64(broadcastType.Level) / 100.0 // Convert level to coefficient (0.0 to 1.0)

		n.mu.Lock()

		n.relayMap[messageID] = time.Now()
		n.receiveMap[messageID] = []p2p.NodeID{} // Reset duplicates for this relay

		n.mu.Unlock()

		// Simulate node processing delay
		time.Sleep(time.Duration(n.delay) * time.Millisecond)

		// Send message to all connected nodes with wave propagation
		for conn, delay := range n.connections {
			if n.checkReceiving(messageID, conn) {
				continue
			}

			wg.Add(1)

			go func(conn *Node, delay p2p.Delay) {
				defer wg.Done()

				// Simulate network transmission delay
				time.Sleep(time.Duration(delay) * time.Millisecond)

				conn.relayWave(messageID, n, 0, coef, wg)
			}(conn, delay)
		}
	}
}

// relayBasic handles message relay using basic flooding algorithm
func (n *Node) relayBasic(messageID p2p.MessageID, from *Node, wg *sync.WaitGroup) {
	n.mu.Lock()

	// Check if message has already been processed by this node
	if _, ok := n.relayMap[messageID]; ok {
		n.receiveMap[messageID] = append(n.receiveMap[messageID], from.id) // Track duplicate sender
		n.mu.Unlock()
		return
	} else {
		// First time receiving this message
		n.relayMap[messageID] = time.Now()
		n.receiveMap[messageID] = []p2p.NodeID{from.id} // Reset duplicates for this relay
		n.mu.Unlock()
	}

	// Simulate node processing delay
	time.Sleep(time.Duration(n.delay) * time.Millisecond)

	// Forward message to all connected nodes except the sender
	for conn, delay := range n.connections {
		if conn == from {
			continue // Skip excluded node
		}

		if n.checkReceiving(messageID, conn) {
			continue
		}

		wg.Add(1)

		go func(conn *Node, delay p2p.Delay) {
			defer wg.Done()

			// Simulate network transmission delay
			time.Sleep(time.Duration(delay) * time.Millisecond)

			conn.relayBasic(messageID, n, wg)
		}(conn, delay)
	}
}

// relayWave handles message relay using wave-based algorithm with hop-based selective forwarding
func (n *Node) relayWave(messageID p2p.MessageID, from *Node, hop int, coef float64, wg *sync.WaitGroup) {
	n.mu.Lock()

	// Check if message has already been processed by this node
	if _, ok := n.relayMap[messageID]; ok {
		n.receiveMap[messageID] = append(n.receiveMap[messageID], from.id) // Track duplicate sender
		n.mu.Unlock()
		return
	} else {
		// First time receiving this message
		n.relayMap[messageID] = time.Now()
		n.receiveMap[messageID] = []p2p.NodeID{from.id} // Reset duplicates for this relay
		n.mu.Unlock()
	}

	// Simulate node processing delay
	time.Sleep(time.Duration(n.delay) * time.Millisecond)

	if hop%2 == 0 {
		// Even hop: forward to all connected nodes (full propagation)
		for conn, delay := range n.connections {
			if conn == from {
				continue // Skip excluded node
			}

			if n.checkReceiving(messageID, conn) {
				continue
			}

			wg.Add(1)

			go func(conn *Node, delay p2p.Delay) {
				defer wg.Done()

				// Simulate network transmission delay
				time.Sleep(time.Duration(delay) * time.Millisecond)

				conn.relayWave(messageID, n, hop+1, coef, wg)
			}(conn, delay)
		}
	} else {
		// Odd hop: forward to limited number of nodes based on coefficient
		randN := rand.Intn(len(n.connections))
		i := 0
		send := 0

		// Calculate maximum number of nodes to send to (at least 1)
		maxSend := max(int(coef*float64(len(n.connections))), 1)

		// Create a copy of connections to modify during iteration
		copiedConnections := make(map[*Node]p2p.Delay, len(n.connections))
		for conn, delay := range n.connections {
			copiedConnections[conn] = delay
		}

		// Randomly select and send to maxSend number of nodes
		for send < maxSend {
			flag := false

			for conn, delay := range copiedConnections {
				if conn == from {
					continue // Skip excluded node
				}

				if n.checkReceiving(messageID, conn) {
					continue
				}

				flag = true

				// Send to randomly selected node
				if i == randN {
					wg.Add(1)

					go func(conn *Node, delay p2p.Delay) {
						defer wg.Done()

						// Simulate network transmission delay
						time.Sleep(time.Duration(delay) * time.Millisecond)

						conn.relayWave(messageID, n, hop+1, coef, wg)
					}(conn, delay)

					// Remove selected node from available connections
					delete(copiedConnections, conn)

					i = 0
					send++

					break
				} else {
					i++
				}
			}

			if !flag {
				break // No more connections to send to
			}
		}
	}
}

// checkReceiving checks if a node has already received a message from this connection
// to prevent duplicate transmissions and optimize network efficiency
func (n *Node) checkReceiving(relayNumber p2p.MessageID, conn *Node) bool {
	n.mu.RLock()
	defer n.mu.RUnlock()

	// Check if the connection has already relayed this message
	for _, dupID := range n.receiveMap[relayNumber] {
		if dupID == conn.id {
			return true // Skip if this node has already relayed this message
		}
	}

	return false
}
