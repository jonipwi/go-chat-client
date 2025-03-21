package state

import (
	"testing"
)

func TestClientState(t *testing.T) {
	clientState := NewClientState("testuser")

	// Test IsConnected
	if clientState.IsConnected() != false {
		t.Error("Expected initial connection status to be false")
	}

	// Test SetConnected
	clientState.SetConnected(true)
	if clientState.IsConnected() != true {
		t.Error("Expected connection status to be true after SetConnected(true)")
	}

	// Test TrackMessageSent
	clientState.TrackMessageSent()
	if clientState.messagesSent != 1 {
		t.Error("Expected messagesSent to be 1 after TrackMessageSent")
	}

	// Test TrackHeartbeatSent
	clientState.TrackHeartbeatSent()
	if clientState.heartbeatsSent != 1 {
		t.Error("Expected heartbeatsSent to be 1 after TrackHeartbeatSent")
	}

	// Test AddConnectionError
	clientState.AddConnectionError("test error")
	if len(clientState.connectionErrors) != 1 {
		t.Error("Expected connectionErrors to have 1 error after AddConnectionError")
	}

	// Test GetStats
	stats := clientState.GetStats()
	if stats == "" {
		t.Error("Expected GetStats to return a non-empty string")
	}
}
