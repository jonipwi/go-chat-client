package commands

import (
	"fmt"
	"strings"

	"../state"
)

// ProcessCommand handles all user input commands
func ProcessCommand(clientState *state.ClientState, input string, host string, port int) {
	// Split the input into parts
	parts := strings.Split(input, " ")
	command := parts[0]

	// Determine the appropriate command handler based on the command
	switch command {
	case "/global":
		handleGlobalMessage(clientState, parts)
	case "/group":
		handleGroupMessage(clientState, parts)
	case "/guild":
		handleGuildMessage(clientState, parts)
	case "/private":
		handlePrivateMessage(clientState, parts)
	case "/create":
		handleCreateRoom(clientState, parts)
	case "/join":
		handleJoinRoom(clientState, parts)
	case "/list":
		handleListRooms(clientState, parts)
	case "/ping":
		handlePing(clientState)
	case "/test":
		handleTestEvent(clientState)
	case "/heartbeat":
		handleManualHeartbeat(clientState)
	case "/stats":
		handleStats(clientState)
	case "/username":
		handleUsernameChange(clientState, parts)
	case "/debug":
		handleDebug(clientState)
	case "/forcereconnect":
		handleForceReconnect(clientState, host, port)
	case "/errors":
		handleConnectionErrors(clientState)
	case "/help":
		PrintCommands()
	default:
		handleDefaultInput(clientState, input)
	}
}

// PrintCommands displays available chat commands
func PrintCommands() {
	fmt.Println("\n==== Chat Commands ====")
	fmt.Println("/global <message> - Send a message to global chat")
	fmt.Println("/group <group_id> <message> - Send a message to a group")
	fmt.Println("/guild <guild_id> <message> - Send a message to a guild")
	fmt.Println("/private <user_id> <message> - Send a private message")
	fmt.Println("/create <type> <name> - Create a new room (type: group or guild)")
	fmt.Println("/join <room_id> - Join a room")
	fmt.Println("/list <type> - List available rooms (type: group or guild)")
	fmt.Println("/ping - Send a ping to test connection")
	fmt.Println("/test - Send a test event")
	fmt.Println("/heartbeat - Send a manual heartbeat")
	fmt.Println("/stats - Show client connection statistics")
	fmt.Println("/username <new_name> - Change your username")
	fmt.Println("/debug - Display connection debugging information")
	fmt.Println("/forcereconnect - Force a reconnection attempt")
	fmt.Println("/errors - Display connection error history")
	fmt.Println("/exit - Disconnect and exit")
	fmt.Println("=======================\n")
}

// Helper function to check if client is connected
func checkClientConnected(clientState *state.ClientState) bool {
	if !clientState.IsConnected() || clientState.Client() == nil {
		fmt.Println("Error: Not connected to server")
		return false
	}
	return true
}
