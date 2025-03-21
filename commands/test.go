package commands

import (
	"testing"

	"github.com/jonipwi/go-chat-client/state"
)

func TestHandlePing(t *testing.T) {
	clientState := state.NewClientState("testuser")
	handlePing(clientState) // Should not panic even if not connected
}

func TestHandleStats(t *testing.T) {
	clientState := state.NewClientState("testuser")
	handleStats(clientState) // Should not panic
}

func TestHandleUsernameChange(t *testing.T) {
	clientState := state.NewClientState("testuser")
	handleUsernameChange(clientState, []string{"/username", "newuser"})

	if clientState.GetUsername() != "newuser" {
		t.Error("Expected username to be updated to 'newuser'")
	}
}

func TestHandleGlobalMessage(t *testing.T) {
	clientState := state.NewClientState("testuser")
	handleGlobalMessage(clientState, []string{"/global", "Hello World"})

	if clientState.messagesSent != 1 {
		t.Error("Expected messagesSent to be 1 after sending a global message")
	}
}
