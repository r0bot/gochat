package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/r0bot/gochat/internal/pkg/input"
	"github.com/r0bot/gochat/internal/pkg/messages"
	"net"
	"os"
	"strings"
)

type client struct {
	connection net.Conn
}

func (client *client) sendMessage(clientMessage messages.ClientMessage) {
	if client.connection != nil {
		encoder := json.NewEncoder(client.connection)
		err := encoder.Encode(&clientMessage)
		if err != nil {
			fmt.Println("Could not send message to server!")
		}
	} else {
		// Attempt to reconnect to the server
		fmt.Println("Reconnecting..")
		client.connectToServer()
		// Call the function again
		client.sendMessage(clientMessage)
	}
}

func (client *client) connectToServer() {
	// TODO use configuration rather than hardcoded values
	connection, err := net.Dial("tcp", "127.0.0.1:3310")
	if err != nil {
		fmt.Println(err)
	}
	client.connection = connection
	fmt.Println("Connected to server.")
	//Listen for messages from the server spawn a routine to handle that
	go client.listen()
}

func (client *client) listen() {
	for {
		message := make([]byte, 4096)
		length, err := client.connection.Read(message)
		if err != nil {
			_ = client.connection.Close()
			client.connection = nil
			break
		}
		if length > 0 {
			fmt.Println("Message from server: " + string(message))
		}
	}
}

func executeCommand(usrInput input.UserInput, client *client) {
	commandComponents := strings.Fields(usrInput.Payload)
	//Check the first field of the string as this should be the command name
	//If exit quit the process
	if commandComponents[0] == "exit" {
		fmt.Println("Exiting.")
		os.Exit(0)
	}
	// If pm create the message and send it
	if commandComponents[0] == "pm" {
		// the destination of the message (ClientId) should be the second part of the command
		// and the payload teh third
		clientMessage := messages.ClientMessage{
			MessageType: messages.PM,
			Payload:     commandComponents[2],
			Destination: commandComponents[1],
		}
		client.sendMessage(clientMessage)
		return
	}

	fmt.Println("Unrecognised command.")
}

func main() {
	client := client{}
	client.connectToServer()
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Type a message or a command (using \\)")
	for {
		stdInput, err := reader.ReadString('\n')
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
		}
		usrInput := input.ParseInput(stdInput)

		switch usrInput.InputType {
		case input.Command:
			executeCommand(usrInput, &client)
		case input.Message:
			clientMessage := messages.ClientMessage{MessageType: messages.Broadcast, Payload: usrInput.Payload}
			client.sendMessage(clientMessage)
		default:
			fmt.Println("Invalid input.")
		}
	}
}
