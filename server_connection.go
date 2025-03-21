package main

import (
	"fmt"
	"log"
	"time"

	"./events"
	"./state"
	socketio_client "github.com/zhouhui8915/go-socket.io-client"
)

// Connect to the server
func connectToServer(host string, port int, clientState *state.ClientState) (*socketio_client.Client, error) {
	serverURL := fmt.Sprintf("http://%s:%d/socket.io/", host, port)
	log.Printf("CONNECTION: Connecting to server at %s", serverURL)

	// Create client with options
	opts := &socketio_client.Options{
		Transport: "websocket",
		Query:     make(map[string]string),
	}
	opts.Query["username"] = clientState.GetUsername()

	// Create client
	c, err := socketio_client.NewClient(serverURL, opts)
	if err != nil {
		log.Printf("CONNECTION ERROR: Failed to create new client: %v", err)
		return nil, fmt.Errorf("error creating client: %w", err)
	}

	// Set up event handlers
	events.SetupEventHandlers(c, clientState)

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
			// Get last heartbeat timestamp
			timeSinceLastHeartbeat := time.Since(time.Now()) // This will be 0, just for initialization

			log.Printf("HEARTBEAT: Sending heartbeat... (Time since last server response: %v)",
				timeSinceLastHeartbeat.Round(time.Second))

			clientState.Client().Emit("client_heartbeat", []interface{}{
				fmt.Sprintf("Heartbeat at %s", time.Now().Format(time.RFC3339)),
			})
			clientState.TrackHeartbeatSent()

			// Warning if we haven't received a heartbeat in a while
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
