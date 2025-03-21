package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/jonipwi/go-chat-client/state"
)

// ProcessCommand handles all user input commands
func ProcessCommand(clientState *state.ClientState, input string, host string, port int) {
	fmt.Println() // Add a newline before command output
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
	fmt.Println() // Add a newline after command output
}

// PrintCommands displays available chat commands
func PrintCommands() {
	fmt.Println("\n==== Chat Commands ====")
	fmt.Println("/global <message>    - Send a message to global chat")
	fmt.Println("/group <group_id> <message> - Send a message to a group")
	fmt.Println("/guild <guild_id> <message> - Send a message to a guild")
	fmt.Println("/private <user_id> <message> - Send a private message")
	fmt.Println("/create <type> <name>  - Create a new room (type: group or guild)")
	fmt.Println("/join <room_id>     - Join a room")
	fmt.Println("/list <type>        - List available rooms (type: group or guild)")
	fmt.Println("/ping               - Send a ping to test connection")
	fmt.Println("/test               - Send a test event")
	fmt.Println("/heartbeat          - Send a manual heartbeat")
	fmt.Println("/stats              - Show client connection statistics")
	fmt.Println("/username <new_name> - Change your username")
	fmt.Println("/debug              - Display connection debugging information")
	fmt.Println("/forcereconnect     - Force a reconnection attempt")
	fmt.Println("/errors             - Display connection error history")
	fmt.Println("/exit               - Disconnect and exit")
	fmt.Println("=======================\n")
}

// Helper function to check if client is connected
func checkClientConnected(clientState *state.ClientState) bool {
	if !clientState.IsConnected() || clientState.Client() == nil {
		fmt.Println("❌ Error: Not connected to server")
		fmt.Println("   Use /forcereconnect to attempt reconnection")
		return false
	}
	return true
}

// handlePing sends a ping to test the connection
func handlePing(clientState *state.ClientState) {
	fmt.Println("🏓 Testing connection with ping...")
	if !checkClientConnected(clientState) {
		return
	}

	err := clientState.Client().Emit("ping", []interface{}{
		fmt.Sprintf("Ping from %s", clientState.GetUsername()),
	})

	if err != nil {
		fmt.Printf("❌ Error sending ping: %v\n", err)
		return
	}

	fmt.Println("✅ Ping sent successfully!")
}

// handleManualHeartbeat sends a manual heartbeat
func handleManualHeartbeat(clientState *state.ClientState) {
	fmt.Println("💓 Sending manual heartbeat...")
	if !checkClientConnected(clientState) {
		return
	}

	err := clientState.Client().Emit("client_heartbeat", []interface{}{
		fmt.Sprintf("Manual heartbeat from %s at %s",
			clientState.GetUsername(),
			time.Now().Format(time.RFC3339)),
	})

	if err != nil {
		fmt.Printf("❌ Error sending heartbeat: %v\n", err)
		return
	}

	clientState.TrackHeartbeatSent()
	fmt.Println("✅ Manual heartbeat sent successfully!")
}

// handleStats displays client statistics
func handleStats(clientState *state.ClientState) {
	fmt.Println("📊 Client Statistics:")
	stats := clientState.GetStats()
	fmt.Printf("%s\n", stats)
}

// handleUsernameChange changes the client's username
func handleUsernameChange(clientState *state.ClientState, parts []string) {
	if len(parts) < 2 {
		fmt.Println("Error: New username required")
		return
	}

	newUsername := parts[1]
	clientState.SetUsername(newUsername)
	fmt.Printf("Username changed to: %s\n", newUsername)

	// If connected, notify the server
	if clientState.IsConnected() && clientState.Client() != nil {
		err := clientState.Client().Emit("username_change", []interface{}{newUsername})
		if err != nil {
			fmt.Printf("Error notifying server of username change: %v\n", err)
		}
	}
}

// handleGlobalMessage handles the /global command
func handleGlobalMessage(clientState *state.ClientState, args []string) {
	if !checkClientConnected(clientState) {
		return
	}
	if len(args) < 2 {
		fmt.Println("❌ Usage: /global <message>")
		return
	}
	message := strings.Join(args[1:], " ")
	fmt.Printf("🌐 Sending global message: %s\n", message)

	err := clientState.Client().Emit("global_message", []interface{}{message})
	if err != nil {
		fmt.Printf("❌ Error sending global message: %v\n", err)
		return
	}
	clientState.TrackMessageSent()
	fmt.Println("✅ Global message sent successfully!")
}

// handleDebug displays debugging information
func handleDebug(clientState *state.ClientState) {
	fmt.Printf("\nDebug Information:\n")
	fmt.Printf("Connected: %v\n", clientState.IsConnected())
	fmt.Printf("Username: %s\n", clientState.GetUsername())
	fmt.Printf("Client ID: %s\n", clientState.GetClientID())
	fmt.Printf("Current Room: %s\n", clientState.GetCurrentRoom())
	fmt.Printf("Last Activity: %v\n", clientState.GetLastActivity().Format(time.RFC3339))
	fmt.Printf("\nConnection Errors:\n")
	for _, err := range clientState.GetConnectionErrors() {
		fmt.Printf("- %s\n", err)
	}
	fmt.Println()
}

// handleForceReconnect forces a reconnection attempt
func handleForceReconnect(clientState *state.ClientState, host string, port int) {
	if clientState.IsConnected() && clientState.Client() != nil {
		clientState.CloseConnection()
	}

	fmt.Printf("Attempting to reconnect to %s:%d...\n", host, port)
	serverURL := fmt.Sprintf("http://%s:%d/socket.io/", host, port)
	err := clientState.ConnectToServer(serverURL)
	if err != nil {
		fmt.Printf("Reconnection failed: %v\n", err)
		clientState.AddConnectionError(fmt.Sprintf("Reconnection failed: %v", err))
	} else {
		fmt.Println("Reconnection successful")
	}
}

// handleConnectionErrors displays connection error history
func handleConnectionErrors(clientState *state.ClientState) {
	errors := clientState.GetConnectionErrors()
	if len(errors) == 0 {
		fmt.Println("\nNo connection errors recorded")
		return
	}

	fmt.Printf("\nConnection Error History:\n")
	for _, err := range errors {
		fmt.Printf("- %s\n", err)
	}
	fmt.Println()
}

// handleDefaultInput handles any input that doesn't match a command
func handleDefaultInput(clientState *state.ClientState, input string) {
	if !checkClientConnected(clientState) {
		return
	}

	// Treat as a global message
	err := clientState.Client().Emit("global_message", []interface{}{input})
	if err != nil {
		fmt.Printf("Error sending message: %v\n", err)
		return
	}

	clientState.TrackMessageSent()
	fmt.Println("Message sent successfully")
}

// handleTestEvent sends a test event to the server
func handleTestEvent(clientState *state.ClientState) {
	if !checkClientConnected(clientState) {
		return
	}

	err := clientState.Client().Emit("test_event", []interface{}{
		fmt.Sprintf("Test event from %s", clientState.GetUsername()),
	})

	if err != nil {
		fmt.Printf("Error sending test event: %v\n", err)
		return
	}

	clientState.TrackMessageSent()
	fmt.Println("Test event sent successfully")
}

// handleGroupMessage handles sending messages to a group
func handleGroupMessage(clientState *state.ClientState, args []string) {
	if !checkClientConnected(clientState) {
		return
	}
	if len(args) < 3 {
		fmt.Println("Usage: /group <group_id> <message>")
		return
	}
	groupID := args[1]
	message := strings.Join(args[2:], " ")
	err := clientState.Client().Emit("group_message", []interface{}{groupID, message})
	if err != nil {
		fmt.Printf("Error sending group message: %v\n", err)
		return
	}
	clientState.TrackMessageSent()
	fmt.Println("Group message sent successfully")
}

// handleGuildMessage handles sending messages to a guild
func handleGuildMessage(clientState *state.ClientState, args []string) {
	if !checkClientConnected(clientState) {
		return
	}
	if len(args) < 3 {
		fmt.Println("Usage: /guild <guild_id> <message>")
		return
	}
	guildID := args[1]
	message := strings.Join(args[2:], " ")
	err := clientState.Client().Emit("guild_message", []interface{}{guildID, message})
	if err != nil {
		fmt.Printf("Error sending guild message: %v\n", err)
		return
	}
	clientState.TrackMessageSent()
	fmt.Println("Guild message sent successfully")
}

// handlePrivateMessage handles sending private messages
func handlePrivateMessage(clientState *state.ClientState, args []string) {
	if !checkClientConnected(clientState) {
		return
	}
	if len(args) < 3 {
		fmt.Println("Usage: /private <user_id> <message>")
		return
	}
	userID := args[1]
	message := strings.Join(args[2:], " ")
	err := clientState.Client().Emit("private_message", []interface{}{userID, message})
	if err != nil {
		fmt.Printf("Error sending private message: %v\n", err)
		return
	}
	clientState.TrackMessageSent()
	fmt.Println("Private message sent successfully")
}

// handleCreateRoom handles creating a new room
func handleCreateRoom(clientState *state.ClientState, args []string) {
	if !checkClientConnected(clientState) {
		return
	}
	if len(args) < 3 {
		fmt.Println("Usage: /create <type> <name>")
		return
	}
	roomType := args[1]
	roomName := args[2]
	err := clientState.Client().Emit("create_room", []interface{}{roomType, roomName})
	if err != nil {
		fmt.Printf("Error creating room: %v\n", err)
		return
	}
	fmt.Printf("Room creation request sent for %s: %s\n", roomType, roomName)
}

// handleJoinRoom handles joining a room
func handleJoinRoom(clientState *state.ClientState, args []string) {
	if !checkClientConnected(clientState) {
		return
	}
	if len(args) < 2 {
		fmt.Println("Usage: /join <room_id>")
		return
	}
	roomID := args[1]
	err := clientState.Client().Emit("join_room", []interface{}{roomID})
	if err != nil {
		fmt.Printf("Error joining room: %v\n", err)
		return
	}
	clientState.SetCurrentRoom(roomID)
	fmt.Printf("Joined room: %s\n", roomID)
}

// handleListRooms handles listing available rooms
func handleListRooms(clientState *state.ClientState, args []string) {
	if !checkClientConnected(clientState) {
		return
	}
	if len(args) < 2 {
		fmt.Println("Usage: /list <type>")
		return
	}
	roomType := args[1]
	err := clientState.Client().Emit("list_rooms", []interface{}{roomType})
	if err != nil {
		fmt.Printf("Error requesting room list: %v\n", err)
		return
	}
	fmt.Printf("Room list request sent for type: %s\n", roomType)
}
