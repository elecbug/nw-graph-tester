package node

type Node struct {
	ID          uint64
	Connections map[uint64]int64
	NodeDelay   int64
}
