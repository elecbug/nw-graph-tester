package p2p

import "fmt"

// NodeID represents a unique identifier for a node in the P2P network
type NodeID uint64

// MessageID represents a unique identifier for a message in the network
type MessageID uint64

// Delay represents the network delay in milliseconds
type Delay uint64

// BroadcastType defines the type and configuration of broadcast method
type BroadcastType struct {
	Type  string // The broadcast algorithm type (BasicPublish or WavePublish)
	Level int    // The level parameter for WavePublish (unused for BasicPublish)
}

// Constants for different broadcast algorithm types
const (
	BasicPublish = "BasicPublish" // Simple flooding-based broadcast
	WavePublish  = "WavePublish"  // Wave-based broadcast with level control
)

// String returns a human-readable string representation of the broadcast type
func (bt BroadcastType) String() string {
	switch bt.Type {
	case BasicPublish:
		return "BasicPublish"
	case WavePublish:
		// Include the level parameter in the string for WavePublish
		return fmt.Sprintf("WavePublish-%d", bt.Level)
	default:
		return "Unknown"
	}
}
