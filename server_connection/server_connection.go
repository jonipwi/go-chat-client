// server_connection.go
package server_connection

import (
	"fmt"
	"log"
	"time"

	"github.com/jonipwi/go-chat-client/state"
	socketio_client "github.com/zhouhui8915/go-socket.io-client"
)

// ConnectToServer connects to the server
func ConnectToServer(host string, port int, clientState *state.ClientState) (*socketio_client.Client, error) {
	serverURL := fmt.Sprintf("http://%s:%d/socket.io/", host, port)
	log.Printf("CONNECTION: Connecting to server at %s", serverURL)

	opts := &socketio_client.Options{
		Transport: "websocket",
		Query:     make(map[string]string),
	}
	opts.Query["username"] = clientState.GetUsername()

	var c *socketio_client.Client
	var err error
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		c, err = socketio_client.NewClient(serverURL, opts)
		if err == nil {
			break
		}
		log.Printf("CONNECTION ERROR: Attempt %d/%d failed: %v", i+1, maxRetries, err)
		if i < maxRetries-1 {
			time.Sleep(time.Second * 2)
		}
	}

	if err != nil {
		log.Printf("CONNECTION ERROR: Failed to create client after %d attempts: %v", maxRetries, err)
		return nil, fmt.Errorf("error creating client: %w", err)
	}

	// Setup event handlers
	c.On("connect", func(msg string) {
		log.Printf("CONNECTION: Connected with client ID: %s", msg)
		clientState.SetClientID(msg)
		clientState.SetConnected(true)
	})

	c.On("disconnect", func() {
		log.Println("CONNECTION: Disconnected from server")
		clientState.SetConnected(false)
	})

	c.On("chat message", func(msg string) {
		log.Printf("CHAT: %s", msg)
		clientState.TrackMessageReceived()
	})

	c.On("message", func(msg string) {
		log.Printf("SERVER MESSAGE: %s", msg)
		clientState.TrackMessageReceived()
	})

	c.On("heartbeat", func(data map[string]interface{}) {
		log.Printf("HEARTBEAT: Received server heartbeat: %v", data)
		clientState.TrackHeartbeatReceived()
	})

	c.On("room joined", func(roomID string) {
		log.Printf("ROOM: Joined room %s", roomID)
		clientState.SetCurrentRoom(roomID)
	})

	c.On("private message", func(sender string, msg string) {
		log.Printf("PRIVATE: From %s: %s", sender, msg)
		clientState.TrackMessageReceived()
	})

	c.On("user joined", func(username string) {
		log.Printf("ROOM: User %s joined", username)
	})

	c.On("user left", func(username string) {
		log.Printf("ROOM: User %s left", username)
	})

	c.On("typing", func(username string) {
		log.Printf("TYPING: %s is typing...", username)
	})

	c.On("stop typing", func(username string) {
		log.Printf("TYPING: %s stopped typing", username)
	})

	log.Println("CONNECTION: Client connected successfully")
	return c, nil
}

// StartHeartbeat starts a custom heartbeat mechanism
func StartHeartbeat(clientState *state.ClientState) {
	log.Println("HEARTBEAT: Starting custom heartbeat mechanism")
	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C
		if clientState.IsConnected() && clientState.Client() != nil {
			lastHeartbeat := clientState.GetLastActivity()
			timeSinceLastHeartbeat := time.Since(lastHeartbeat)

			// Check if ClientID is set before sending heartbeat
			clientID := clientState.GetClientID()
			if clientID == "" {
				log.Println("HEARTBEAT: Skipping heartbeat - no client ID set yet")
				continue
			}

			log.Printf("HEARTBEAT: Sending heartbeat... (Time since last server response: %v)",
				timeSinceLastHeartbeat.Round(time.Second))

			// Fixed: Send single string parameter instead of array
			heartbeatMsg := fmt.Sprintf("Heartbeat from %s at %s",
				clientID,
				time.Now().Format(time.RFC3339))

			err := clientState.Client().Emit("client_heartbeat", heartbeatMsg)

			if err != nil {
				log.Printf("HEARTBEAT ERROR: Failed to send heartbeat: %v", err)
				clientState.AddConnectionError(fmt.Sprintf("Heartbeat send failed: %v", err))
				continue
			}

			clientState.TrackHeartbeatSent()

			if timeSinceLastHeartbeat > 2*time.Minute {
				log.Printf("HEARTBEAT WARNING: No server response in %v!",
					timeSinceLastHeartbeat.Round(time.Second))
				clientState.AddConnectionError(fmt.Sprintf("No heartbeat response in %v",
					timeSinceLastHeartbeat.Round(time.Second)))
			}
		} else {
			log.Println("HEARTBEAT: Skipping heartbeat - not connected")
		}
	}
}

// ReportStats reports client stats periodically
func ReportStats(clientState *state.ClientState) {
	log.Println("STATS: Starting periodic stats reporting")
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		<-ticker.C
		stats := clientState.GetStats()
		log.Printf("CLIENT STATS: %s", stats)
	}
}
