package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/jonipwi/go-chat-client/state"
	"github.com/jonipwi/go-chat-client/utils"
)

// handlePing sends a ping to test connection
func handlePing(clientState *state.ClientState) {
	if !checkClientConnected(clientState) {
		return
	}

	utils.Logger.Println("COMMAND: Executing ping command")

	clientState.Client().Emit("ping", []interface{}{"Manual ping from Go client"})
	utils.Logger.Println("PING SENT: Manual ping")
}

// handleTestEvent sends a test event
func handleTestEvent(clientState *state.ClientState) {
	if !checkClientConnected(clientState) {
		return
	}

	utils.Logger.Println("COMMAND: Executing test event command")

	clientState.Client().Emit("test_event", []interface{}{"Manual test from Go client"})
	utils.Logger.Println("TEST EVENT SENT")
}

// handleManualHeartbeat sends a manual heartbeat
func handleManualHeartbeat(clientState *state.ClientState) {
	if !checkClientConnected(clientState) {
		return
	}

	utils.Logger.Println("COMMAND: Executing manual heartbeat command")

	clientState.Client().Emit("client_heartbeat", []interface{}{"Manual heartbeat from Go client"})
	clientState.TrackHeartbeatSent()
	utils.Logger.Println("MANUAL HEARTBEAT SENT")
}

// handleStats displays client connection statistics
func handleStats(clientState *state.ClientState) {
	utils.Logger.Println("COMMAND: Executing stats command")

	stats := clientState.GetStats()
	fmt.Printf("Client Status: %s\n", stats)
	utils.Logger.Printf("STATS DISPLAYED: %s", stats)
}

// handleUsernameChange changes the client's username
func handleUsernameChange(clientState *state.ClientState, parts []string) {
	if !checkClientConnected(clientState) {
		return
	}

	if len(parts) < 2 {
		utils.Logger.Println("COMMAND ERROR: Invalid username change format")
		fmt.Println("Usage: /username <new_name>")
		return
	}

	// Join all remaining parts as the username
	newName := strings.Join(parts[1:], " ")
	utils.Logger.Printf("COMMAND: Changing username to: %s", newName)

	clientState.Client().Emit("set_username", []interface{}{newName})
	clientState.username = newName
	clientState.TrackMessageSent()
	utils.Logger.Printf("USERNAME CHANGE REQUEST SENT: new_name=%s", newName)
}

// handleDebug displays connection debugging information
func handleDebug(clientState *state.ClientState) {
	utils.Logger.Println("COMMAND: Executing debug information command")

	fmt.Println("\n==== Debug Information ====")
	fmt.Printf("Connection State: %v\n", clientState.IsConnected())
	fmt.Printf("Client ID: %s\n", clientState.clientID)
	fmt.Printf("Last Activity: %s\n", clientState.lastActivity.Format(time.RFC3339))

	// Print connection errors
	fmt.Println("\nConnection Error History:")
	if len(clientState.connectionErrors) == 0 {
		fmt.Println("No connection errors recorded")
	} else {
		for i, err := range clientState.connectionErrors {
			fmt.Printf("%d. %s\n", i+1, err)
		}
	}

	// Print socket information
	if clientState.Client() != nil {
		fmt.Println("\nSocket Information:")
		fmt.Printf("Socket ID: %s\n", clientState.Client().Id())
	} else {
		fmt.Println("\nSocket Information: No active socket client")
	}
	fmt.Println("=======================")
}

// handleForceReconnect forces a reconnection to the server
func handleForceReconnect(clientState *state.ClientState, host string, port int) {
	utils.Logger.Println("COMMAND: Executing forced reconnection")

	if clientState.IsConnected() {
		// Disconnect first
		utils.Logger.Println("FORCED RECONNECT: Disconnecting current connection")
		clientState.SetConnected(false)
		if clientState.Client() != nil {
			clientState.Client().Close()
		}
	}

	utils.Logger.Println("FORCED RECONNECT: Initiating new connection")
	clientState.lastReconnectAttempt = time.Now()

	// You'll need to import the function from the main package
	newClient, err := connectToServer(host, port, clientState)
	if err != nil {
		utils.Logger.Printf("FORCED RECONNECT ERROR: %v", err)
		clientState.AddConnectionError(fmt.Sprintf("Forced reconnect failed: %v", err))
		fmt.Println("Forced reconnection failed:", err)
		return
	}

	clientState.client = newClient
	utils.Logger.Println("FORCED RECONNECT: New client created, waiting for connection events")

	// Wait a bit to see if connection is established
	fmt.Println("Waiting for connection events...")
	time.Sleep(3 * time.Second)

	if clientState.IsConnected() {
		utils.Logger.Println("FORCED RECONNECT: Successfully reconnected")
		fmt.Println("Successfully reconnected to server")
	} else {
		utils.Logger.Println("FORCED RECONNECT: Failed to establish connection after 3 seconds")
		clientState.AddConnectionError("Forced reconnect timeout after 3 seconds")
		fmt.Println("Reconnection attempt in progress - check status with /debug")
	}
}

// handleConnectionErrors displays the connection error history
func handleConnectionErrors(clientState *state.ClientState) {
	utils.Logger.Println("COMMAND: Displaying connection error history")

	fmt.Println("\n==== Connection Error History ====")
	if len(clientState.connectionErrors) == 0 {
		fmt.Println("No connection errors recorded")
	} else {
		for i, err := range clientState.connectionErrors {
			fmt.Printf("%d. %s\n", i+1, err)
		}
	}
	fmt.Println("=================================")
}
