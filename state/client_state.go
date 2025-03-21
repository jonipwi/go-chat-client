package state

import (
	"fmt"
	"log"
	"sync"
	"time"

	socketio_client "github.com/zhouhui8915/go-socket.io-client"
)

// ClientState keeps track of the client state
type ClientState struct {
	connected             bool
	client                *socketio_client.Client
	username              string
	clientID              string
	mutex                 sync.Mutex
	lastHeartbeatSent     time.Time
	lastHeartbeatReceived time.Time
	lastActivity          time.Time
	connectionStarted     time.Time
	messagesSent          int
	messagesReceived      int
	heartbeatsSent        int
	heartbeatsReceived    int
	onConnectHandler      func(string)
	connectionErrors      []string
	lastReconnectAttempt  time.Time
	CurrentRoom           string
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

// Getters and Setters

// IsConnected returns the current connection status
func (s *ClientState) IsConnected() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.connected
}

// Client returns the current socket.io client
func (s *ClientState) Client() *socketio_client.Client {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.client
}

// SetClient updates the socket.io client
func (s *ClientState) SetClient(client *socketio_client.Client) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.client = client
}

// SetConnected updates the connection status
func (s *ClientState) SetConnected(connected bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	wasConnected := s.connected
	s.connected = connected

	if connected && !wasConnected {
		s.connectionStarted = time.Now()
		fmt.Printf("CONNECTION STATE: Connected at %v\n", s.connectionStarted.Format(time.RFC3339))
	} else if !connected && wasConnected {
		duration := time.Since(s.connectionStarted).Round(time.Second)
		fmt.Printf("CONNECTION STATE: Disconnected after %v\n", duration)
	}
}

// UpdateActivity updates the last activity timestamp
func (s *ClientState) UpdateActivity() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.lastActivity = time.Now()
}

// GetLastActivity returns the last activity timestamp
func (s *ClientState) GetLastActivity() time.Time {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.lastActivity
}

// SetLastReconnectAttempt updates the last reconnect attempt timestamp
func (s *ClientState) SetLastReconnectAttempt(t time.Time) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.lastReconnectAttempt = t
}

// AddConnectionError adds a new connection error to the history
func (s *ClientState) AddConnectionError(err string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Keep only the last 10 errors
	if len(s.connectionErrors) >= 10 {
		s.connectionErrors = s.connectionErrors[1:]
	}

	s.connectionErrors = append(s.connectionErrors, fmt.Sprintf("[%s] %s",
		time.Now().Format("15:04:05"), err))
}

// GetStats returns a formatted string of connection statistics
func (s *ClientState) GetStats() string {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var connStatus string
	var connDuration time.Duration

	if s.connected {
		connStatus = "Connected"
		connDuration = time.Since(s.connectionStarted).Round(time.Second)
	} else {
		connStatus = "Disconnected"
		if !s.connectionStarted.IsZero() {
			connDuration = s.lastActivity.Sub(s.connectionStarted).Round(time.Second)
		}
	}

	timeSinceLastHeartbeatSent := "Never"
	if !s.lastHeartbeatSent.IsZero() {
		timeSinceLastHeartbeatSent = time.Since(s.lastHeartbeatSent).Round(time.Second).String()
	}

	timeSinceLastHeartbeatReceived := "Never"
	if !s.lastHeartbeatReceived.IsZero() {
		timeSinceLastHeartbeatReceived = time.Since(s.lastHeartbeatReceived).Round(time.Second).String()
	}

	// Add reconnection info
	var reconnInfo string
	if !s.lastReconnectAttempt.IsZero() {
		reconnInfo = fmt.Sprintf(", Last reconnect attempt: %s ago",
			time.Since(s.lastReconnectAttempt).Round(time.Second))
	}

	return fmt.Sprintf("Status: %s, Duration: %v, Client ID: %s, Username: %s, "+
		"Messages Sent: %d, Messages Received: %d, Heartbeats Sent: %d, Heartbeats Received: %d, "+
		"Time Since Last Heartbeat Sent: %s, Time Since Last Heartbeat Received: %s%s",
		connStatus, connDuration, s.clientID, s.username,
		s.messagesSent, s.messagesReceived, s.heartbeatsSent, s.heartbeatsReceived,
		timeSinceLastHeartbeatSent, timeSinceLastHeartbeatReceived, reconnInfo)
}

// Message and Heartbeat Tracking Methods

// TrackMessageSent increments the messages sent counter
func (s *ClientState) TrackMessageSent() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.messagesSent++
	s.lastActivity = time.Now()
}

// TrackMessageReceived increments the messages received counter
func (s *ClientState) TrackMessageReceived() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.messagesReceived++
	s.lastActivity = time.Now()
}

// TrackHeartbeatSent increments the heartbeats sent counter
func (s *ClientState) TrackHeartbeatSent() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.heartbeatsSent++
	s.lastHeartbeatSent = time.Now()
	s.lastActivity = s.lastHeartbeatSent
}

// TrackHeartbeatReceived increments the heartbeats received counter
func (s *ClientState) TrackHeartbeatReceived() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.heartbeatsReceived++
	s.lastHeartbeatReceived = time.Now()
	s.lastActivity = s.lastHeartbeatReceived
}

// Additional Utility Methods

// GetUsername returns the current username
func (s *ClientState) GetUsername() string {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.username
}

// SetUsername updates the username
func (s *ClientState) SetUsername(username string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.username = username
}

// GetConnectionErrors returns the list of connection errors
func (s *ClientState) GetConnectionErrors() []string {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return append([]string{}, s.connectionErrors...)
}

// GetClientID returns the current client ID
func (s *ClientState) GetClientID() string {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.clientID
}

// SetClientID updates the client ID
func (s *ClientState) SetClientID(clientID string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.clientID = clientID
}

// GetOnConnectHandler returns the current onConnect handler
func (s *ClientState) GetOnConnectHandler() func(string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.onConnectHandler
}

// SetOnConnectHandler updates the onConnect handler
func (s *ClientState) SetOnConnectHandler(handler func(string)) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.onConnectHandler = handler
}

// GetCurrentRoom returns the current room
func (s *ClientState) GetCurrentRoom() string {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.CurrentRoom
}

// SetCurrentRoom updates the current room
func (s *ClientState) SetCurrentRoom(room string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.CurrentRoom = room
}

// ConnectToServer connects to the Socket.IO server
func ConnectToServer(host string, port int, clientState *ClientState) (*socketio_client.Client, error) {
	serverURL := fmt.Sprintf("http://%s:%d/socket.io/", host, port)
	log.Printf("CONNECTION: Connecting to server at %s", serverURL)

	opts := &socketio_client.Options{
		Transport: "websocket",
		Query:     make(map[string]string),
	}
	opts.Query["username"] = clientState.GetUsername()

	client, err := socketio_client.NewClient(serverURL, opts)
	if err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}

	clientState.SetClient(client)

	// Register event handlers here, if needed
	client.RegisterHandler("message", func(c *socketio_client.Client, msg interface{}) {
		log.Printf("Received message: %v", msg)
	})

	return client, nil
}
