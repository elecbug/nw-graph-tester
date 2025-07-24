package p2p

type NodeID uint64
type MessageID uint64
type Delay uint64
type BroadcastType int

const (
	BasicPublish BroadcastType = iota
	FloodPublish
)
