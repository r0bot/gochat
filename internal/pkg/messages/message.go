package messages

type MessageType int

const (
	Broadcast MessageType = iota
	PM
)

type ClientMessage struct {
	ClientId    string
	MessageType MessageType
	Payload     string
	Destination string
}
