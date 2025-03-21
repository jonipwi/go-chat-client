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

	clientState.SetConnected(true)
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

			log.Printf("HEARTBEAT: Sending heartbeat... (Time since last server response: %v)",
				timeSinceLastHeartbeat.Round(time.Second))

			err := clientState.Client().Emit("client_heartbeat", []interface{}{
				fmt.Sprintf("Heartbeat from %s at %s",
					clientState.GetClientID(),
					time.Now().Format(time.RFC3339)),
			})

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
		if clientState.IsConnected() {
			stats := clientState.GetStats()
			log.Printf("CLIENT STATS: %s", stats)
		}
	}
}
