package clients

import (
	"testing"
	"time"
)

func TestAddClient(t *testing.T) {
	clientId := "testClientId"
	manager := ClientManager{make(map[string]*Client), nil}
	client := Client{clientId, nil, make(chan []byte), make(chan []byte)}
	go manager.AddClient(&client)
	//Sleep so the routine can execute
	time.Sleep(time.Millisecond)
	if manager.Clients["testClientId"] == nil {
		t.Errorf("Clients Manager should contain client with id: %s.", clientId)
	}

	// When the client output channel is closed the manager should remove the client from the map
	close(client.Output)
	//Sleep so the routine can remove the client after the channel close
	time.Sleep(time.Millisecond)
	if manager.Clients["testClientId"] != nil {
		t.Errorf("Clients Manager should not contain client with id: %s.", clientId)
	}
}
