package sender

type Packet struct {
	Type string
	Command string
	Data string // JSON string
	RawData bool
	Timeout int32
	Timestamp int64
	ConnectionId string
}
