package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/r0bot/gochat/internal/pkg/input"
	"github.com/r0bot/gochat/internal/pkg/messages"
	"net"
	"os"
)

type client struct {
	connection net.Conn
}

func (client *client) sendMessage(usrInput input.UserInput) {
	clientMessage := messages.ClientMessage{"", messages.Broadcast, usrInput.Payload}
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
		client.sendMessage(usrInput)
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
	fmt.Println(usrInput.Payload)
	if usrInput.Payload == "exit" {
		fmt.Println("Exiting.")
		os.Exit(0)
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
			client.sendMessage(usrInput)
		default:
			fmt.Println("Invalid input.")
		}
	}
}
