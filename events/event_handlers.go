package events

import (
	"fmt"
	"log"
	"time"

	socketio "github.com/googollee/go-socket.io"
	"github.com/jonipwi/go-chat-client/state"
)

// Message structure for incoming chat messages
type Message struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Sender    string    `json:"sender"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// User structure for user information
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

// Room structure for room information
type Room struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	Members   []string  `json:"members"`
}

// setupEventHandlers configures event listeners for the Socket.IO client
func setupEventHandlers(c *socketio.Client, clientState *state.ClientState) {
	// Connection event
	c.On("connect", func() {
		clientID := c.Id()
		log.Printf("CONNECTION EVENT: Connected successfully to server with client ID: %s", clientID)
		clientState.clientID = clientID
		clientState.SetConnected(true)

		// Call the custom connection handler if set
		if clientState.onConnectHandler != nil {
			clientState.onConnectHandler(clientID)
		}

		// Set username after connecting
		log.Printf("CONNECTION EVENT: Setting username to: %s", clientState.username)
		c.Emit("set_username", []interface{}{clientState.username})
		clientState.TrackMessageSent()

		// Send a test message to global chat
		log.Println("CONNECTION EVENT: Sending test message to global chat")
		c.Emit("global_message", []interface{}{"Hello from Go client!"})
		clientState.TrackMessageSent()
	})

	// Disconnection event
	c.On("disconnect", func() {
		log.Printf("DISCONNECTION EVENT: Disconnected from server. Client ID: %s", clientState.clientID)

		// Add more context about the disconnection
		if clientState.IsConnected() {
			log.Printf("DISCONNECTION CONTEXT: Unexpected disconnection - client was previously connected")
			clientState.AddConnectionError(fmt.Sprintf("Unexpected disconnection from %s", clientState.clientID))
		} else {
			log.Printf("DISCONNECTION CONTEXT: Expected disconnection - client was already marked as disconnected")
		}

		clientState.SetConnected(false)
	})

	// Error event
	c.On("error", func(err error) {
		log.Printf("SOCKET ERROR: Client %s experienced error: %v", clientState.clientID, err)
		clientState.AddConnectionError(fmt.Sprintf("Socket error: %v", err))
	})

	// Chat message events
	c.On("chat_message", func(msg Message) {
		log.Printf("RECEIVED MESSAGE: [%s] %s: %s (ID: %s, Timestamp: %v)",
			msg.Type, msg.Sender, msg.Content, msg.ID, msg.Timestamp)
		clientState.TrackMessageReceived()

		// Forward to console for user to see
		fmt.Printf("📥 [%s] %s: %s\n", msg.Type, msg.Sender, msg.Content)
	})

	// Global message event
	c.On("global_message", func(msg Message) {
		log.Printf("GLOBAL MESSAGE: %s: %s", msg.Sender, msg.Content)
		clientState.TrackMessageReceived()
		fmt.Printf("🌐 [GLOBAL] %s: %s\n", msg.Sender, msg.Content)
	})

	// Group message event
	c.On("group_message", func(msg Message) {
		log.Printf("GROUP MESSAGE: [%s] %s: %s", msg.Type, msg.Sender, msg.Content)
		clientState.TrackMessageReceived()
		fmt.Printf("👥 [GROUP:%s] %s: %s\n", msg.Type, msg.Sender, msg.Content)
	})

	// Guild message event
	c.On("guild_message", func(msg Message) {
		log.Printf("GUILD MESSAGE: [%s] %s: %s", msg.Type, msg.Sender, msg.Content)
		clientState.TrackMessageReceived()
		fmt.Printf("🏰 [GUILD:%s] %s: %s\n", msg.Type, msg.Sender, msg.Content)
	})

	// Private message event
	c.On("private_message", func(msg Message) {
		log.Printf("PRIVATE MESSAGE: From %s: %s", msg.Sender, msg.Content)
		clientState.TrackMessageReceived()
		fmt.Printf("🔒 [PRIVATE] %s: %s\n", msg.Sender, msg.Content)
	})

	// Server heartbeat response
	c.On("server_heartbeat", func(msg string) {
		log.Printf("RECEIVED HEARTBEAT: %s", msg)
		clientState.TrackHeartbeatReceived()
		log.Printf("HEARTBEAT RECEIVED")
	})

	// Test event response
	c.On("test_event", func(data interface{}) {
		log.Printf("TEST EVENT RECEIVED: %v", data)
		clientState.TrackMessageReceived()
	})

	// Username-related events
	c.On("username_updated", func(user User) {
		log.Printf("USERNAME UPDATED: %s -> %s (User ID: %s)",
			clientState.username, user.Username, user.ID)
		clientState.username = user.Username
		clientState.TrackMessageReceived()
	})

	// Username suggestion event
	c.On("username_suggestion", func(suggestion string) {
		log.Printf("USERNAME SUGGESTION: %s", suggestion)
		fmt.Printf("🏷️ Suggested Username: %s\n", suggestion)
	})

	// Room-related events
	c.On("room_created", func(room Room) {
		log.Printf("ROOM CREATED: %s (ID: %s, Type: %s, Members: %v)",
			room.Name, room.ID, room.Type, room.Members)
		clientState.TrackMessageReceived()
		fmt.Printf("🚪 Room Created: %s (ID: %s, Type: %s)\n", room.Name, room.ID, room.Type)
	})

	c.On("user_joined", func(user User) {
		log.Printf("USER JOINED: %s (ID: %s) joined a room", user.Username, user.ID)
		clientState.TrackMessageReceived()
		fmt.Printf("👤 %s joined the room\n", user.Username)
	})

	c.On("user_left", func(user User) {
		log.Printf("USER LEFT: %s (ID: %s) left a room", user.Username, user.ID)
		clientState.TrackMessageReceived()
		fmt.Printf("🚶 %s left the room\n", user.Username)
	})

	// Room listing events
	c.On("rooms_list", func(rooms []Room) {
		log.Printf("ROOMS LIST: Received list of %d rooms", len(rooms))
		fmt.Println("📋 Available Rooms:")
		for i, room := range rooms {
			log.Printf("  Room %d: %s (ID: %s, Type: %s, Members: %d)",
				i+1, room.Name, room.ID, room.Type, len(room.Members))
			fmt.Printf("  %d. %s (ID: %s, Type: %s, Members: %d)\n",
				i+1, room.Name, room.ID, room.Type, len(room.Members))
		}
		clientState.TrackMessageReceived()
	})

	// Error handling events
	c.On("server_error", func(errorMsg string) {
		log.Printf("SERVER ERROR: %s", errorMsg)
		clientState.AddConnectionError(fmt.Sprintf("Server error: %s", errorMsg))
		fmt.Printf("❌ Server Error: %s\n", errorMsg)
		clientState.TrackMessageReceived()
	})

	// Connection status events
	c.On("connection_status", func(status map[string]interface{}) {
		log.Printf("CONNECTION STATUS: %v", status)
		fmt.Println("🔗 Connection Status:")
		for key, value := range status {
			fmt.Printf("  %s: %v\n", key, value)
		}
	})

	// Rate limit warning
	c.On("rate_limit_warning", func(warning string) {
		log.Printf("RATE LIMIT WARNING: %s", warning)
		fmt.Printf("⚠️ Rate Limit Warning: %s\n", warning)
	})

	// Ping/Pong events
	c.On("ping", func() {
		log.Println("PING RECEIVED: Responding with pong")
		c.Emit("pong", nil)
	})

	c.On("pong", func() {
		log.Println("PONG RECEIVED")
		clientState.TrackHeartbeatReceived()
	})

	// Authentication events
	c.On("authentication_required", func() {
		log.Println("AUTHENTICATION REQUIRED")
		fmt.Println("🔐 Authentication is required to continue")
	})

	c.On("authentication_success", func() {
		log.Println("AUTHENTICATION SUCCESSFUL")
		fmt.Println("🔓 Authentication successful")
	})

	c.On("authentication_failed", func(reason string) {
		log.Printf("AUTHENTICATION FAILED: %s", reason)
		fmt.Printf("❌ Authentication failed: %s\n", reason)
		clientState.AddConnectionError(fmt.Sprintf("Authentication failed: %s", reason))
	})
}
