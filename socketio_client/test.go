package socketio_client

import (
	"fmt"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func startMockServer(t *testing.T) *websocket.Conn {
	// This function will start a mock WebSocket server for testing purposes.
	// It will simulate receiving messages and sending responses.
	server := websocket.NewServer(func(conn *websocket.Conn) {
		// Mock server will echo the "ping" message back to the client
		for {
			var msg map[string]interface{}
			if err := conn.ReadJSON(&msg); err != nil {
				t.Logf("Error reading message from client: %v", err)
				break
			}
			if msg["event"] == "ping" {
				err := conn.WriteJSON(map[string]interface{}{
					"event":   "pong",
					"message": "pong response",
				})
				if err != nil {
					t.Logf("Error sending message to client: %v", err)
					break
				}
			}
		}
	})
	return server
}

// Test NewClient function
func TestNewClient(t *testing.T) {
	// Start a mock server
	server := startMockServer(t)
	defer server.Close()

	serverURL := fmt.Sprintf("ws://%s:%d", "localhost", 8080) // Replace with mock server address
	options := &Options{
		Transport: "websocket",
		Query:     make(map[string]string),
	}

	// Initialize client
	client, err := NewClient(serverURL, options)
	assert.NoError(t, err, "Error creating client")
	assert.True(t, client.IsConnected(), "Client should be connected")
}

// Test SendMessage function
func TestSendMessage(t *testing.T) {
	// Start a mock server
	server := startMockServer(t)
	defer server.Close()

	serverURL := fmt.Sprintf("ws://%s:%d", "localhost", 8080) // Replace with mock server address
	options := &Options{
		Transport: "websocket",
		Query:     make(map[string]string),
	}

	client, err := NewClient(serverURL, options)
	assert.NoError(t, err, "Error creating client")

	// Send a message
	err = client.SendMessage("ping", nil)
	assert.NoError(t, err, "Error sending message")
}

// Test ReceiveMessages function
func TestReceiveMessages(t *testing.T) {
	// Start a mock server
	server := startMockServer(t)
	defer server.Close()

	serverURL := fmt.Sprintf("ws://%s:%d", "localhost", 8080) // Replace with mock server address
	options := &Options{
		Transport: "websocket",
		Query:     make(map[string]string),
	}

	client, err := NewClient(serverURL, options)
	assert.NoError(t, err, "Error creating client")

	// Register the "pong" event handler
	client.RegisterHandler("pong", func(c *Client, message interface{}) {
		t.Log("Received Pong: ", message)
		assert.Equal(t, "pong response", message, "Unexpected response message")
	})

	// Send "ping" message and expect to receive "pong"
	err = client.SendMessage("ping", nil)
	assert.NoError(t, err, "Error sending message")

	// Give the server some time to respond
	time.Sleep(1 * time.Second)
}

// Test Close function
func TestCloseClient(t *testing.T) {
	// Start a mock server
	server := startMockServer(t)
	defer server.Close()

	serverURL := fmt.Sprintf("ws://%s:%d", "localhost", 8080) // Replace with mock server address
	options := &Options{
		Transport: "websocket",
		Query:     make(map[string]string),
	}

	client, err := NewClient(serverURL, options)
	assert.NoError(t, err, "Error creating client")

	// Close the client connection
	err = client.Close()
	assert.NoError(t, err, "Error closing client")

	// Check if the client is disconnected
	assert.False(t, client.IsConnected(), "Client should be disconnected")
}

// Test Heartbeat function
func TestHeartbeat(t *testing.T) {
	// Start a mock server
	server := startMockServer(t)
	defer server.Close()

	serverURL := fmt.Sprintf("ws://%s:%d", "localhost", 8080) // Replace with mock server address
	options := &Options{
		Transport: "websocket",
		Query:     make(map[string]string),
	}

	client, err := NewClient(serverURL, options)
	assert.NoError(t, err, "Error creating client")

	// Send heartbeat message
	err = client.Heartbeat()
	assert.NoError(t, err, "Error sending heartbeat")

	// Give some time for the heartbeat to be sent and acknowledged
	time.Sleep(1 * time.Second)
}
