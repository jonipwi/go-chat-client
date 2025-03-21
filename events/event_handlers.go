package events

import (
	"fmt"
	"time"

	"../state"
	socketio_client "github.com/zhouhui8915/go-socket.io-client"
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

// SetupEventHandlers configures event listeners for the Socket.IO client
func SetupEventHandlers(client *socketio_client.Client, clientState *state.ClientState) {
	client.On("error", func() {
		fmt.Println("Error occurred")
	})

	client.On("connect", func(args ...interface{}) {
		fmt.Println("Connected to server")
		if len(args) > 0 {
			if id, ok := args[0].(string); ok {
				clientState.SetClientID(id)
			}
		}
	})

	client.On("disconnect", func() {
		fmt.Println("Disconnected from server")
	})

	client.On("message", func(msg string) {
		fmt.Printf("Received message: %s\n", msg)
	})

	client.On("chat message", func(msg string) {
		fmt.Printf("Received chat message: %s\n", msg)
	})

	client.On("user joined", func(username string) {
		fmt.Printf("User joined: %s\n", username)
	})

	client.On("user left", func(username string) {
		fmt.Printf("User left: %s\n", username)
	})

	client.On("typing", func(username string) {
		fmt.Printf("User %s is typing...\n", username)
	})

	client.On("stop typing", func(username string) {
		fmt.Printf("User %s stopped typing\n", username)
	})

	client.On("user list", func(users []string) {
		fmt.Printf("Current users: %v\n", users)
	})

	client.On("private message", func(from string, msg string) {
		fmt.Printf("Private message from %s: %s\n", from, msg)
	})

	client.On("room joined", func(room string) {
		fmt.Printf("Joined room: %s\n", room)
		clientState.SetCurrentRoom(room)
	})

	client.On("room left", func(room string) {
		fmt.Printf("Left room: %s\n", room)
		if clientState.GetCurrentRoom() == room {
			clientState.SetCurrentRoom("")
		}
	})
}
