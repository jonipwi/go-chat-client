package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/jonipwi/go-chat-client/commands"
	"github.com/jonipwi/go-chat-client/server_connection"
	"github.com/jonipwi/go-chat-client/state"
	"github.com/jonipwi/go-chat-client/utils"
)

func main() {
	// Configure logging
	utils.Logger.Println("STARTUP: Starting Go Socket.IO Chat Client with debug logging...")

	// Server configuration
	host := "127.0.0.1"
	port := 8000
	username := "GoClient"

	// Create client state
	clientState := state.NewClientState(username)

	// Start stats reporting
	go server_connection.ReportStats(clientState)

	// Connect to the server
	utils.Logger.Println("STARTUP: Initiating connection to server...")
	c, err := server_connection.ConnectToServer(host, port, clientState)
	if err != nil {
		utils.Logger.Printf("CONNECTION ERROR: Failed on initial connection to server: %v", err)
		clientState.AddConnectionError(fmt.Sprintf("Initial connection failed: %v", err))
	} else {
		clientState.SetClient(c)
		utils.Logger.Println("STARTUP: Initial connection established")
	}

	// Print available commands
	fmt.Println("\n=== Welcome to Go Chat Client ===")
	fmt.Printf("Server: %s:%d\n", host, port)
	fmt.Printf("Username: %s\n\n", username)
	commands.PrintCommands()

	// Start heartbeat goroutine to keep connection alive
	utils.Logger.Println("STARTUP: Starting heartbeat mechanism...")
	go server_connection.StartHeartbeat(clientState)

	// Start the command loop
	utils.Logger.Println("STARTUP: Starting command loop, ready for user input")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\n> ")
		if !scanner.Scan() {
			utils.Logger.Println("INPUT: Scanner closed")
			break
		}

		input := scanner.Text()
		if input == "" {
			continue
		}

		utils.Logger.Printf("INPUT: User entered: %s", input)

		// Handle exit command
		if input == "/exit" {
			utils.Logger.Println("COMMAND: Exit requested, shutting down client")
			break
		}

		// Process the command
		commands.ProcessCommand(clientState, input, host, port)
	}

	// Attempt to gracefully close the connection if connected
	if clientState.IsConnected() && clientState.Client() != nil {
		utils.Logger.Println("SHUTDOWN: Closing client connection...")
		clientState.CloseConnection()
	}

	utils.Logger.Println("SHUTDOWN: Disconnecting from server...")
	// Final stats report
	utils.Logger.Printf("SHUTDOWN: Final client stats - %s", clientState.GetStats())
}
