package state

import (
	"fmt"
	"log"
	"time"

	socketio_client "github.com/zhouhui8915/go-socket.io-client"
)

// ClientState keeps track of the client state
type ClientState struct {
	connected             bool
	client                *socketio_client.Client
	username              string
	clientID              string
	lastHeartbeatSent     time.Time
	lastHeartbeatReceived time.Time
	lastActivity          time.Time
	connectionStarted     time.Time
	messagesSent          int
	messagesReceived      int
	heartbeatsSent        int
	heartbeatsReceived    int
	connectionErrors      []string
	lastReconnectAttempt  time.Time
	currentRoom           string
}

// NewClientState creates a new ClientState instance
func NewClientState(username string) *ClientState {
	return &ClientState{
		connected:        false,
		username:         username,
		lastActivity:     time.Now(),
		connectionErrors: make([]string, 0, 10),
	}
}

// IsConnected returns the current connection status
func (cs *ClientState) IsConnected() bool {
	return cs.connected
}

// Client returns the current socket.io client
func (cs *ClientState) Client() *socketio_client.Client {
	return cs.client
}

// SetClient updates the socket.io client
func (cs *ClientState) SetClient(client *socketio_client.Client) {
	cs.client = client
}

// SetConnected updates the connection status
func (cs *ClientState) SetConnected(connected bool) {
	wasConnected := cs.connected
	cs.connected = connected

	if connected && !wasConnected {
		cs.connectionStarted = time.Now()
		log.Printf("CONNECTION STATE: Connected at %v\n", cs.connectionStarted.Format(time.RFC3339))
	} else if !connected && wasConnected {
		duration := time.Since(cs.connectionStarted).Round(time.Second)
		log.Printf("CONNECTION STATE: Disconnected after %v\n", duration)
	}
}

// UpdateActivity updates the last activity timestamp
func (cs *ClientState) UpdateActivity() {
	cs.lastActivity = time.Now()
}

// GetLastActivity returns the last activity timestamp
func (cs *ClientState) GetLastActivity() time.Time {
	return cs.lastActivity
}

// SetLastReconnectAttempt updates the last reconnect attempt timestamp
func (cs *ClientState) SetLastReconnectAttempt(t time.Time) {
	cs.lastReconnectAttempt = t
}

// AddConnectionError adds a new connection error to the history
func (cs *ClientState) AddConnectionError(err string) {
	if len(cs.connectionErrors) >= 10 {
		cs.connectionErrors = cs.connectionErrors[1:]
	}

	cs.connectionErrors = append(cs.connectionErrors, fmt.Sprintf("[%s] %s",
		time.Now().Format("15:04:05"), err))
}

// GetStats returns a formatted string of connection statistics
func (cs *ClientState) GetStats() string {
	var connStatus string
	var connDuration time.Duration

	if cs.connected {
		connStatus = "Connected"
		connDuration = time.Since(cs.connectionStarted).Round(time.Second)
	} else {
		connStatus = "Disconnected"
		if !cs.connectionStarted.IsZero() {
			connDuration = cs.lastActivity.Sub(cs.connectionStarted).Round(time.Second)
		}
	}

	timeSinceLastHeartbeatSent := "Never"
	if !cs.lastHeartbeatSent.IsZero() {
		timeSinceLastHeartbeatSent = time.Since(cs.lastHeartbeatSent).Round(time.Second).String()
	}

	timeSinceLastHeartbeatReceived := "Never"
	if !cs.lastHeartbeatReceived.IsZero() {
		timeSinceLastHeartbeatReceived = time.Since(cs.lastHeartbeatReceived).Round(time.Second).String()
	}

	// Add reconnection info
	var reconnInfo string
	if !cs.lastReconnectAttempt.IsZero() {
		reconnInfo = fmt.Sprintf(", Last reconnect attempt: %s ago",
			time.Since(cs.lastReconnectAttempt).Round(time.Second))
	}

	return fmt.Sprintf("Status: %s, Duration: %v, Client ID: %s, Username: %s, "+
		"Messages Sent: %d, Messages Received: %d, Heartbeats Sent: %d, Heartbeats Received: %d, "+
		"Time Since Last Heartbeat Sent: %s, Time Since Last Heartbeat Received: %s%s",
		connStatus, connDuration, cs.clientID, cs.username,
		cs.messagesSent, cs.messagesReceived, cs.heartbeatsSent, cs.heartbeatsReceived,
		timeSinceLastHeartbeatSent, timeSinceLastHeartbeatReceived, reconnInfo)
}

// TrackMessageSent increments the messages sent counter
func (cs *ClientState) TrackMessageSent() {
	cs.messagesSent++
	cs.lastActivity = time.Now()
}

// TrackMessageReceived increments the messages received counter
func (cs *ClientState) TrackMessageReceived() {
	cs.messagesReceived++
	cs.lastActivity = time.Now()
}

// TrackHeartbeatSent increments the heartbeats sent counter
func (cs *ClientState) TrackHeartbeatSent() {
	cs.heartbeatsSent++
	cs.lastHeartbeatSent = time.Now()
	cs.lastActivity = cs.lastHeartbeatSent
}

// TrackHeartbeatReceived increments the heartbeats received counter
func (cs *ClientState) TrackHeartbeatReceived() {
	cs.heartbeatsReceived++
	cs.lastHeartbeatReceived = time.Now()
	cs.lastActivity = cs.lastHeartbeatReceived
}

// GetUsername returns the current username
func (cs *ClientState) GetUsername() string {
	return cs.username
}

// SetUsername updates the username
func (cs *ClientState) SetUsername(username string) {
	cs.username = username
}

// GetConnectionErrors returns the list of connection errors
func (cs *ClientState) GetConnectionErrors() []string {
	return append([]string{}, cs.connectionErrors...)
}

// GetClientID returns the current client ID
func (cs *ClientState) GetClientID() string {
	return cs.clientID
}

// SetClientID updates the client ID
func (cs *ClientState) SetClientID(clientID string) {
	cs.clientID = clientID
}

// GetCurrentRoom returns the current room
func (cs *ClientState) GetCurrentRoom() string {
	return cs.currentRoom
}

// SetCurrentRoom updates the current room
func (cs *ClientState) SetCurrentRoom(room string) {
	cs.currentRoom = room
}

// ConnectToServer establishes a connection to the WebSocket server
func (cs *ClientState) ConnectToServer(serverURL string) error {
	opts := &socketio_client.Options{
		Transport: "websocket",
		Query:     make(map[string]string),
	}
	opts.Query["username"] = cs.username

	client, err := socketio_client.NewClient(serverURL, opts)
	if err != nil {
		return err
	}

	cs.client = client
	cs.connected = true

	return nil
}

// CloseConnection closes the client connection and updates the state
func (cs *ClientState) CloseConnection() {
	if cs.client != nil {
		cs.client.Emit("disconnect", []interface{}{})
		cs.client = nil
		cs.connected = false
		cs.lastActivity = time.Now()
	}
}
