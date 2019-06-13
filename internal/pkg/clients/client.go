package clients

import (
	"encoding/json"
	"fmt"
	"github.com/r0bot/gochat/internal/pkg/messages"
	"net"
)

type Client struct {
	Id     string
	Conn   net.Conn
	Output chan messages.ClientMessage
}

func (client *Client) Send(data []byte) {
	// Listen on the input channel and if so send over the connection
	_, err := client.Conn.Write(data)
	if err != nil {
		fmt.Printf("Error while sending message to client %s. Error %s \n", client.Id, err)
		client.Stop()
	}
}

func (client *Client) listen() {
	defer client.Stop()
	// Listen on the connection for a message and if received send to the client output
	d := json.NewDecoder(client.Conn)
	for {
		var message messages.ClientMessage
		err := d.Decode(&message)
		if err != nil {
			fmt.Println("Error while reading client message", err)
			return
		}
		// Attach teh current clientId to the message
		message.ClientId = client.Id
		client.Output <- message
		client.Send([]byte("Message received."))
	}
}

func (client *Client) Stop() {
	fmt.Printf("Client with id %s disconnecting \n", client.Id)
	// Close the client channels
	_, opened := <-client.Output
	if opened {
		close(client.Output)
	}
	// Close the connection
	err := client.Conn.Close()
	if err != nil {
		fmt.Println("Error while closing client connection", err)
		return
	}
}

func (client *Client) Init() {
	// TODO establish connection handshake and authentication
	fmt.Printf("Client with id %s connected \n", client.Id)

	//Spawn routine to listen on the connection
	go client.listen()

}
