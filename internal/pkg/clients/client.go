package clients

import (
	"bufio"
	"fmt"
	"net"
)

type Client struct {
	Id     string
	Conn   net.Conn
	Input  chan []byte
	Output chan []byte
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
	for {
		data := make([]byte, 1024)
		data, err := bufio.NewReader(client.Conn).ReadBytes('\n')
		if err != nil {
			fmt.Println("Error while reading client message", err)
			return
		}
		if data != nil {
			client.Output <- data
			client.Send([]byte("Message received."))
		}
	}
}

func (client *Client) Stop() {
	// Close the client channels
	close(client.Output)
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
