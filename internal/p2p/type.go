package p2p

type NodeID uint64
type MessageID uint64
type Delay uint64
type BroadcastType int

const (
	BasicPublish BroadcastType = iota
	WavePublish_10
	WavePublish_20
	WavePublish_30
	WavePublish_40
	WavePublish_50
	WavePublish_60
	WavePublish_70
	WavePublish_80
	WavePublish_90
)

func AllBroadcastTypes() []BroadcastType {
	return []BroadcastType{BasicPublish, WavePublish_10, WavePublish_20, WavePublish_30, WavePublish_40, WavePublish_50, WavePublish_60, WavePublish_70, WavePublish_80, WavePublish_90}
}

func (bt BroadcastType) String() string {
	switch bt {
	case BasicPublish:
		return "BasicPublish"
	case WavePublish_10:
		return "WavePublish-10"
	case WavePublish_20:
		return "WavePublish-20"
	case WavePublish_30:
		return "WavePublish-30"
	case WavePublish_40:
		return "WavePublish-40"
	case WavePublish_50:
		return "WavePublish-50"
	case WavePublish_60:
		return "WavePublish-60"
	case WavePublish_70:
		return "WavePublish-70"
	case WavePublish_80:
		return "WavePublish-80"
	case WavePublish_90:
		return "WavePublish-90"
	default:
		return "Unknown"
	}
}
