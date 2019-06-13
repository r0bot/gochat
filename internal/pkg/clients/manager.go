package clients

import (
	"fmt"
	"github.com/r0bot/gochat/internal/pkg/messages"
	"strings"
)

type ClientManager struct {
	Clients   map[string]*Client
	Broadcast chan messages.ClientMessage
	PMs       chan messages.ClientMessage
}

func (manager *ClientManager) Init() {
	for {
		select {
		case message := <-manager.Broadcast:
			var sb strings.Builder
			// TODO handle the error of string builder
			_, _ = fmt.Fprintf(&sb, "Client %s said: %s", message.ClientId, message.Payload)

			for _, client := range manager.Clients {
				// Send the message to all clients except the one the message is from
				if client.Id != message.ClientId {
					client.Send([]byte(sb.String()))
				}
			}
		case message := <-manager.PMs:
			var sb strings.Builder
			// TODO handle the error of string builder
			_, _ = fmt.Fprintf(&sb, "Client %s sent you message: %s", message.ClientId, message.Payload)

			for _, client := range manager.Clients {
				// Send the message only to the clientId specified in the destination
				if client.Id == message.Destination {
					client.Send([]byte(sb.String()))
					break
				}
			}
		}
	}
}

func (manager *ClientManager) AddClient(client *Client) {
	manager.Clients[client.Id] = client
	for message := range client.Output {

		switch message.MessageType {
		case messages.Broadcast:
			fmt.Printf("Broadcast received from client with id %s. Data: %s", client.Id, string(message.Payload))
			manager.Broadcast <- message
		case messages.PM:
			fmt.Printf("PM received from client with id %s for client with id %s. Data: %s",
				client.Id, message.Destination, string(message.Payload),
			)
			manager.PMs <- message
		}

	}
	// If a client Output channel is closed consider it disconnected
	// and remove from the manager
	delete(manager.Clients, client.Id)
}
