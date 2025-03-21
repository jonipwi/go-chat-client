package state

import (
	"fmt"
	"sync"
	"time"

	socketio "github.com/googollee/go-socket.io"
)

// ClientState keeps track of the client state
type ClientState struct {
	connected             bool
	client                *socketio.Client
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

// Getters
func (s *ClientState) IsConnected() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.connected
}

func (s *ClientState) Client() *socketio.Client {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.client
}

// Set connected state
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

// Update last activity time
func (s *ClientState) UpdateActivity() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.lastActivity = time.Now()
}

// Add a connection error to the history
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

// Get connection stats
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

// Message tracking methods
func (s *ClientState) TrackMessageSent() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.messagesSent++
	s.lastActivity = time.Now()
}

func (s *ClientState) TrackMessageReceived() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.messagesReceived++
	s.lastActivity = time.Now()
}

func (s *ClientState) TrackHeartbeatSent() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.heartbeatsSent++
	s.lastHeartbeatSent = time.Now()
	s.lastActivity = s.lastHeartbeatSent
}

func (s *ClientState) TrackHeartbeatReceived() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.heartbeatsReceived++
	s.lastHeartbeatReceived = time.Now()
	s.lastActivity = s.lastHeartbeatReceived
}
