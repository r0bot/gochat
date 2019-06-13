package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/r0bot/gochat/internal/pkg/messages"
	"net"
	"os"
	"strings"
)

type client struct {
	connection net.Conn
}

type userInput struct {
	InputType string
	Payload   string
}

func (client *client) sendMessage(usrInput userInput) {
	clientMessage := messages.ClientMessage{"", "message", usrInput.Payload}
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

func parseInput(input string) userInput {
	usrInput := userInput{}
	if strings.HasPrefix(input, "\\") {
		usrInput.InputType = "command"
		// Trim the slash and white spaces
		input = strings.TrimPrefix(input, "\\")
		usrInput.Payload = strings.TrimSpace(input)
	} else {
		usrInput.InputType = "message"
		usrInput.Payload = input
	}
	return usrInput
}

func executeCommand(usrInput userInput, client *client) {
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
	for {
		fmt.Println("Type a message or a command (using \\)")
		input, err := reader.ReadString('\n')
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
		}
		usrInput := parseInput(input)
		switch usrInput.InputType {
		case "command":
			executeCommand(usrInput, &client)
		case "message":
			client.sendMessage(usrInput)
		default:
			fmt.Println("Invalid input.")
		}
	}
}
