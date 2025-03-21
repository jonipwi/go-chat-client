package events

import (
	"fmt"
	"time"

	"github.com/jonipwi/go-chat-client/state"
	"github.com/jonipwi/go-chat-client/utils"

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
	client.On("error", func(args ...interface{}) {
		errMsg := "Unknown error"
		if len(args) > 0 {
			if err, ok := args[0].(error); ok {
				errMsg = err.Error()
			} else if str, ok := args[0].(string); ok {
				errMsg = str
			}
		}
		utils.Logger.Printf("ERROR: Socket.IO error: %s", errMsg)
		clientState.AddConnectionError(fmt.Sprintf("Socket.IO error: %s", errMsg))
	})

	client.On("connect", func(args ...interface{}) {
		utils.Logger.Println("EVENT: Connected to server")
		clientState.SetConnected(true)
		if len(args) > 0 {
			if id, ok := args[0].(string); ok {
				clientState.SetClientID(id)
				utils.Logger.Printf("EVENT: Received client ID: %s", id)
			}
		}
	})

	client.On("disconnect", func() {
		utils.Logger.Println("EVENT: Disconnected from server")
		clientState.SetConnected(false)
		clientState.SetCurrentRoom("")
	})

	client.On("message", func(msg string) {
		utils.Logger.Printf("EVENT: Received message: %s", msg)
		clientState.TrackMessageReceived()
	})

	client.On("chat message", func(msg string) {
		utils.Logger.Printf("EVENT: Received chat message: %s", msg)
		clientState.TrackMessageReceived()
	})

	client.On("user joined", func(username string) {
		utils.Logger.Printf("EVENT: User joined: %s", username)
	})

	client.On("user left", func(username string) {
		utils.Logger.Printf("EVENT: User left: %s", username)
	})

	client.On("typing", func(username string) {
		utils.Logger.Printf("EVENT: User %s is typing...", username)
	})

	client.On("stop typing", func(username string) {
		utils.Logger.Printf("EVENT: User %s stopped typing", username)
	})

	client.On("user list", func(users []string) {
		utils.Logger.Printf("EVENT: Current users: %v", users)
	})

	client.On("private message", func(from string, msg string) {
		utils.Logger.Printf("EVENT: Private message from %s: %s", from, msg)
		clientState.TrackMessageReceived()
	})

	client.On("room joined", func(room string) {
		utils.Logger.Printf("EVENT: Joined room: %s", room)
		clientState.SetCurrentRoom(room)
	})

	client.On("room left", func(room string) {
		utils.Logger.Printf("EVENT: Left room: %s", room)
		if clientState.GetCurrentRoom() == room {
			clientState.SetCurrentRoom("")
		}
	})

	client.On("heartbeat", func(args ...interface{}) {
		utils.Logger.Println("EVENT: Received server heartbeat")
		clientState.TrackHeartbeatReceived()
	})
}
