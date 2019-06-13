package input

import "strings"

type MessageType int

const (
	Message MessageType = iota
	Command
)

type UserInput struct {
	InputType MessageType
	Payload   string
}

func ParseInput(input string) UserInput {
	usrInput := UserInput{}
	if strings.HasPrefix(input, "\\") {
		usrInput.InputType = Command
		// Trim the slash and white spaces
		input = strings.TrimPrefix(input, "\\")
		usrInput.Payload = strings.TrimSpace(input)
	} else {
		usrInput.InputType = Message
		usrInput.Payload = input
	}
	return usrInput
}
