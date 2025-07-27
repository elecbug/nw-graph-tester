package p2p

type NodeID uint64
type MessageID uint64
type Delay uint64
type BroadcastType int

const (
	BasicPublish BroadcastType = iota
	WavePublish_1
	WavePublish_2
	WavePublish_3
	WavePublish_4
	WavePublish_5
	WavePublish_6
	WavePublish_7
)

func (bt BroadcastType) String() string {
	switch bt {
	case BasicPublish:
		return "BasicPublish"
	case WavePublish_1:
		return "WavePublish-1"
	case WavePublish_2:
		return "WavePublish-2"
	case WavePublish_3:
		return "WavePublish-3"
	case WavePublish_4:
		return "WavePublish-4"
	case WavePublish_5:
		return "WavePublish-5"
	case WavePublish_6:
		return "WavePublish-6"
	case WavePublish_7:
		return "WavePublish-7"
	default:
		return "Unknown"
	}
}
