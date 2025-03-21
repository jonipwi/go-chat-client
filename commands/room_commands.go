package commands

import (
	"fmt"
	"strings"

	"github.com/jonipwi/go-chat-client/state"
	"github.com/jonipwi/go-chat-client/utils"
)

// handleCreateRoom creates a new room (group or guild)
func handleCreateRoom(clientState *state.ClientState, parts []string) {
	if !checkClientConnected(clientState) {
		return
	}

	if len(parts) < 3 {
		utils.Logger.Println("COMMAND ERROR: Invalid room creation format")
		fmt.Println("Usage: /create <type> <name>")
		return
	}

	roomType := parts[1]
	// Join all remaining parts as the room name
	roomName := strings.Join(parts[2:], " ")

	if roomType != "group" && roomType != "guild" {
		utils.Logger.Printf("COMMAND ERROR: Invalid room type: %s", roomType)
		fmt.Println("Room type must be 'group' or 'guild'")
		return
	}

	utils.Logger.Printf("COMMAND: Creating %s room: %s", roomType, roomName)

	// Pack multiple arguments into a slice
	clientState.Client().Emit("create_room", []interface{}{
		map[string]interface{}{
			"type": roomType,
			"name": roomName,
		},
	})
	clientState.TrackMessageSent()
	utils.Logger.Printf("ROOM CREATION REQUEST SENT: type=%s, name=%s", roomType, roomName)
}

// handleJoinRoom joins an existing room
func handleJoinRoom(clientState *state.ClientState, parts []string) {
	if !checkClientConnected(clientState) {
		return
	}

	if len(parts) < 2 {
		utils.Logger.Println("COMMAND ERROR: Invalid join room format")
		fmt.Println("Usage: /join <room_id>")
		return
	}

	roomID := parts[1]
	utils.Logger.Printf("COMMAND: Joining room: %s", roomID)

	clientState.Client().Emit("join_room", []interface{}{roomID})
	clientState.TrackMessageSent()
	utils.Logger.Printf("JOIN ROOM REQUEST SENT: room_id=%s", roomID)
}

// handleListRooms retrieves and displays available rooms
func handleListRooms(clientState *state.ClientState, parts []string) {
	if !checkClientConnected(clientState) {
		return
	}

	if len(parts) < 2 {
		utils.Logger.Println("COMMAND ERROR: Invalid list rooms format")
		fmt.Println("Usage: /list <type>")
		return
	}

	roomType := parts[1]
	if roomType != "group" && roomType != "guild" {
		utils.Logger.Printf("COMMAND ERROR: Invalid room type: %s", roomType)
		fmt.Println("Room type must be 'group' or 'guild'")
		return
	}

	utils.Logger.Printf("COMMAND: Listing rooms of type: %s", roomType)
	clientState.Client().Emit("get_rooms", []interface{}{roomType})
	clientState.TrackMessageSent()
	utils.Logger.Printf("LIST ROOMS REQUEST SENT: type=%s", roomType)
}
