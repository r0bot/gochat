package main

import (
	"fmt"
	guuid "github.com/google/uuid"
	"github.com/r0bot/gochat/internal/pkg/clients"
	"github.com/r0bot/gochat/internal/pkg/messages"
	"net"
)

type Configuration struct {
	Address string
	Network string
}

func getConfig() Configuration {
	// TODO load the configuration from the file
	return Configuration{"127.0.0.1:3310", "tcp"}
}

func main() {
	// Load config
	config := getConfig()
	//Init server
	listener, err := net.Listen(config.Network, config.Address)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Defer to close the server if the program exits
	defer listener.Close()

	//Create client manager
	clientManager := clients.ClientManager{make(map[string]*clients.Client), make(chan messages.ClientMessage, 100)}
	// Init the manager in a routine
	go clientManager.Init()
	fmt.Printf("Server listening on adress %s \n", config.Address)
	// Loop forever to accept connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		client := clients.Client{guuid.New().String(), conn, make(chan messages.ClientMessage)}

		// Spawn a routine to handle every client concurrently
		go client.Init()
		// Add the client to the manager
		go clientManager.AddClient(&client)
	}
}
