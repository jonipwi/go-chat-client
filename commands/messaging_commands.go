package commands

import (
	"fmt"
	"strings"

	"../state"
	"../utils"
)

// handleGlobalMessage sends a message to the global chat
func handleGlobalMessage(clientState *state.ClientState, parts []string) {
	if !checkClientConnected(clientState) {
		return
	}

	if len(parts) < 2 {
		utils.Logger.Println("COMMAND ERROR: Invalid global message format")
		fmt.Println("Usage: /global <message>")
		return
	}

	// Properly join all words after the command into a single message
	message := strings.Join(parts[1:], " ")
	utils.Logger.Printf("COMMAND: Sending global message: %s", message)

	clientState.Client().Emit("global_message", []interface{}{message})
	clientState.TrackMessageSent()
	utils.Logger.Println("GLOBAL MESSAGE SENT")
}

// handleGroupMessage sends a message to a specific group
func handleGroupMessage(clientState *state.ClientState, parts []string) {
	if !checkClientConnected(clientState) {
		return
	}

	if len(parts) < 3 {
		utils.Logger.Println("COMMAND ERROR: Invalid group message format")
		fmt.Println("Usage: /group <group_id> <message>")
		return
	}

	groupID := parts[1]
	message := strings.Join(parts[2:], " ")
	utils.Logger.Printf("COMMAND: Sending group message to %s: %s", groupID, message)

	clientState.Client().Emit("group_message", []interface{}{
		map[string]interface{}{
			"groupId": groupID,
			"message": message,
		},
	})
	clientState.TrackMessageSent()
	utils.Logger.Printf("GROUP MESSAGE SENT: group_id=%s", groupID)
}

// handleGuildMessage sends a message to a specific guild
func handleGuildMessage(clientState *state.ClientState, parts []string) {
	if !checkClientConnected(clientState) {
		return
	}

	if len(parts) < 3 {
		utils.Logger.Println("COMMAND ERROR: Invalid guild message format")
		fmt.Println("Usage: /guild <guild_id> <message>")
		return
	}

	guildID := parts[1]
	message := strings.Join(parts[2:], " ")
	utils.Logger.Printf("COMMAND: Sending guild message to %s: %s", guildID, message)

	clientState.Client().Emit("guild_message", []interface{}{
		map[string]interface{}{
			"guildId": guildID,
			"message": message,
		},
	})
	clientState.TrackMessageSent()
	utils.Logger.Printf("GUILD MESSAGE SENT: guild_id=%s", guildID)
}

// handlePrivateMessage sends a private message to a specific user
func handlePrivateMessage(clientState *state.ClientState, parts []string) {
	if !checkClientConnected(clientState) {
		return
	}

	if len(parts) < 3 {
		utils.Logger.Println("COMMAND ERROR: Invalid private message format")
		fmt.Println("Usage: /private <user_id> <message>")
		return
	}

	userID := parts[1]
	message := strings.Join(parts[2:], " ")
	utils.Logger.Printf("COMMAND: Sending private message to %s: %s", userID, message)

	clientState.Client().Emit("private_message", []interface{}{
		map[string]interface{}{
			"userId":  userID,
			"message": message,
		},
	})
	clientState.TrackMessageSent()
	utils.Logger.Printf("PRIVATE MESSAGE SENT: user_id=%s", userID)
}

// handleDefaultInput handles input without a specific command prefix
func handleDefaultInput(clientState *state.ClientState, input string) {
	if !checkClientConnected(clientState) {
		return
	}

	// If no command is specified, treat as a global message
	if !strings.HasPrefix(input, "/") {
		utils.Logger.Printf("COMMAND: Sending text as global message: %s", input)
		clientState.Client().Emit("global_message", []interface{}{input})
		clientState.TrackMessageSent()
		utils.Logger.Println("IMPLIED GLOBAL MESSAGE SENT")
	} else {
		utils.Logger.Printf("COMMAND ERROR: Unknown command: %s", input)
		fmt.Println("Unknown command. Type /help for available commands.")
	}
}
