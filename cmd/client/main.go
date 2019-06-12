package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

type client struct {
	connection net.Conn
}

func (client *client) sendMessage(message string) {
	if client.connection != nil {
		output := []byte(message)
		_, err := client.connection.Write(output)
		if err != nil {
			fmt.Println("Could not send message to server!")
		}
	} else {
		// Attempt to reconnect to the server
		fmt.Println("Reconnecting..")
		client.connectToServer()
		// Call the function again
		client.sendMessage(message)
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
		message := make([]byte, 1024)
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

func main() {
	client := client{}
	client.connectToServer()
	reader := bufio.NewReader(os.Stdin)
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
		}
		client.sendMessage(input)
	}
}
