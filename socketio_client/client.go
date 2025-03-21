package socketio_client

import (
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// Client represents a socket.io client
type Client struct {
	conn               *websocket.Conn
	serverURL          string
	connected          bool
	lastPingTime       time.Time
	connectionHandlers map[string]func(*Client, interface{})
}

// NewClient creates a new socket.io client instance
func NewClient(serverURL string, options *Options) (*Client, error) {
	// Create the WebSocket connection
	conn, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error establishing websocket connection: %w", err)
	}

	client := &Client{
		conn:               conn,
		serverURL:          serverURL,
		connected:          true,
		connectionHandlers: make(map[string]func(*Client, interface{})),
	}

	// Start receiving messages in the background
	go client.receiveMessages()

	return client, nil
}

// Close closes the socket.io client connection
func (c *Client) Close() error {
	if c.conn == nil {
		return fmt.Errorf("no active connection to close")
	}

	err := c.conn.Close()
	if err != nil {
		return fmt.Errorf("error closing connection: %w", err)
	}

	c.connected = false
	return nil
}

// IsConnected checks if the client is connected to the server
func (c *Client) IsConnected() bool {
	return c.connected
}

// SendMessage sends a message to the server
func (c *Client) SendMessage(event string, message interface{}) error {
	if !c.IsConnected() {
		return fmt.Errorf("cannot send message, client is not connected")
	}

	// Send message to server
	err := c.conn.WriteJSON(map[string]interface{}{
		"event":   event,
		"message": message,
	})
	if err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}
	return nil
}

// ReceiveMessages listens for incoming messages from the server
func (c *Client) receiveMessages() {
	for {
		var msg map[string]interface{}
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		// Handle the received message
		if event, ok := msg["event"].(string); ok {
			if handler, exists := c.connectionHandlers[event]; exists {
				handler(c, msg["message"])
			} else {
				log.Printf("No handler for event: %s", event)
			}
		}
	}
}

// RegisterHandler registers a handler for specific events
func (c *Client) RegisterHandler(event string, handler func(*Client, interface{})) {
	c.connectionHandlers[event] = handler
}

// Heartbeat sends a heartbeat message to the server
func (c *Client) Heartbeat() error {
	return c.SendMessage("heartbeat", nil)
}

// Ping sends a ping message to the server to test the connection
func (c *Client) Ping() error {
	c.lastPingTime = time.Now()
	return c.SendMessage("ping", nil)
}

// GetLastPingTime returns the time of the last ping
func (c *Client) GetLastPingTime() time.Time {
	return c.lastPingTime
}

// Options holds configuration options for the socket.io client
type Options struct {
	Transport string
	Query     map[string]string
}
