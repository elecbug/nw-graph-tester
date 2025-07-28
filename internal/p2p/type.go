package p2p

import "fmt"

type NodeID uint64
type MessageID uint64
type Delay uint64
type BroadcastType struct {
	Type  string
	Level int
}

const (
	BasicPublish = "BasicPublish"
	WavePublish  = "WavePublish"
)

func (bt BroadcastType) String() string {
	switch bt.Type {
	case BasicPublish:
		return "BasicPublish"
	case WavePublish:
		return fmt.Sprintf("WavePublish-%d", bt.Level)
	default:
		return "Unknown"
	}
}
