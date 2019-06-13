package clients

import (
	"fmt"
	"github.com/r0bot/gochat/internal/pkg/messages"
	"strings"
)

type ClientManager struct {
	Clients   map[string]*Client
	Broadcast chan messages.ClientMessage
}

func (manager *ClientManager) Init() {
	//Listen on the disconnect and messages channel
	for message := range manager.Broadcast {
		// Send the message to all clients
		var sb strings.Builder
		// TODO handle the error of string builder
		_, _ = fmt.Fprintf(&sb, "Client %s said: %s", message.ClientId, message.Payload)

		for _, client := range manager.Clients {
			// Send the message to all clients except the one the message is from
			if client.Id != message.ClientId {
				client.Send([]byte(sb.String()))
			}
		}
	}
}

func (manager *ClientManager) AddClient(client *Client) {
	manager.Clients[client.Id] = client
	for message := range client.Output {
		fmt.Printf("Message received from client with id %s. Data: %s", client.Id, string(message.Payload))
		manager.Broadcast <- message
	}
	// If a client Output channel is closed consider it disconnected
	// and remove from the manager
	delete(manager.Clients, client.Id)
}
