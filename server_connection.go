package main

import (
	"fmt"
	"log"
	"time"

	socketio "github.com/googollee/go-socket.io"
	"github.com/jonipwi/go-chat-client/state"
)

// Connect to the server
func connectToServer(host string, port int, clientState *state.ClientState) (*socketio.Client, error) {
	serverURL := fmt.Sprintf("http://%s:%d", host, port)
	log.Printf("CONNECTION: Connecting to server at %s", serverURL)

	// Create client with options
	c, err := socketio.NewClient(serverURL, socketio.NewClientOption())
	if err != nil {
		log.Printf("CONNECTION ERROR: Failed to create new client: %v", err)
		return nil, fmt.Errorf("error creating client: %w", err)
	}

	// Set up event handlers
	setupEventHandlers(c, clientState)

	// Connect the client
	err = c.Connect()
	if err != nil {
		log.Printf("CONNECTION ERROR: Failed to connect: %v", err)
		return nil, fmt.Errorf("error connecting to server: %w", err)
	}

	log.Println("CONNECTION: Client connected and all event handlers set up successfully")
	return c, nil
}

// Start a custom heartbeat to keep the connection alive
func startHeartbeat(clientState *state.ClientState) {
	log.Println("HEARTBEAT: Starting custom heartbeat mechanism")
	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C
		if clientState.IsConnected() && clientState.Client() != nil {
			// Check last heartbeat timestamp
			timeSinceLastHeartbeat := time.Since(clientState.lastHeartbeatReceived)

			log.Printf("HEARTBEAT: Sending heartbeat... (Time since last server response: %v)",
				timeSinceLastHeartbeat.Round(time.Second))

			clientState.Client().Emit("client_heartbeat", []interface{}{
				fmt.Sprintf("Heartbeat at %s", time.Now().Format(time.RFC3339)),
			})
			clientState.TrackHeartbeatSent()

			// Warning if we haven't received a heartbeat in a while
			if !clientState.lastHeartbeatReceived.IsZero() && timeSinceLastHeartbeat > 2*time.Minute {
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

// Report client stats periodically
func reportStats(clientState *state.ClientState) {
	log.Println("STATS: Starting periodic stats reporting")
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		<-ticker.C
		stats := clientState.GetStats()
		log.Printf("CLIENT STATS: %s", stats)
	}
}
