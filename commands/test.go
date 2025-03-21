package commands

import (
	"fmt"
	"strings"
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

	// Check if message was tracked
	stats := clientState.GetStats()
	if !strings.Contains(stats, "Messages Sent: 1") {
		t.Error("Expected messagesSent to be 1 after sending a global message")
	}
}

// TestCommand represents a test command
type TestCommand struct {
	Name        string
	Description string
}

// NewTestCommand creates a new test command
func NewTestCommand() *TestCommand {
	return &TestCommand{
		Name:        "test",
		Description: "Test command",
	}
}

// Execute executes the test command
func (c *TestCommand) Execute(clientState *state.ClientState) error {
	if !checkClientConnected(clientState) {
		return fmt.Errorf("not connected to server")
	}

	// Send a test event
	err := clientState.Client().Emit("test_event", []interface{}{
		fmt.Sprintf("Test event from %s", clientState.GetUsername()),
	})

	if err != nil {
		return fmt.Errorf("error sending test event: %v", err)
	}

	fmt.Println("Test event sent successfully")
	return nil
}
