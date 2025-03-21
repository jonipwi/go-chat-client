// File: main.go (client)
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jonipwi/go-chat-client/server_connection"
	"github.com/jonipwi/go-chat-client/state"
)

func main() {
	log.Println("[CHAT-CLIENT] STARTUP: Starting Go Socket.IO Chat Client with debug logging...")

	// Create client state with default username
	clientState := state.NewClientState("GoClient")

	// Start the stats reporting in a goroutine
	go server_connection.ReportStats(clientState)

	// Connect to server
	log.Println("[CHAT-CLIENT] STARTUP: Initiating connection to server...")
	client, err := server_connection.ConnectToServer("127.0.0.1", 8000, clientState)
	if err != nil {
		log.Fatalf("[CHAT-CLIENT] CONNECTION ERROR: Failed on initial connection to server: %v", err)
	}

	clientState.SetClient(client)

	// Start the heartbeat mechanism in a goroutine
	go server_connection.StartHeartbeat(clientState)

	// Wait a moment for connection to stabilize
	time.Sleep(1 * time.Second)

	// Print welcome message and instructions
	fmt.Println("\n=== Welcome to Go Chat Client ===")
	fmt.Println("Available commands:")
	fmt.Println("  /help - Show this help")
	fmt.Println("  /quit - Exit the client")
	fmt.Println("  /stats - Show connection stats")
	fmt.Println("  /username <new_name> - Change your username")
	fmt.Println("  /join <room_id> - Join a room")
	fmt.Println("  /list groups - List available groups")
	fmt.Println("  /list guilds - List available guilds")
	fmt.Println("  /create group <name> - Create a new group")
	fmt.Println("  /create guild <name> - Create a new guild")
	fmt.Println("  /msg <user_id> <message> - Send private message")
	fmt.Println("  /ping - Send a ping to the server")
	fmt.Println("  /errors - Show recent connection errors")
	fmt.Println("================================================")
	fmt.Printf("You are connected as: %s\n", clientState.GetUsername())
	fmt.Println("Type your message and press Enter to send to current room")
	fmt.Print("> ")

	// Start scanner for user input
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()

		// Handle different commands
		if strings.HasPrefix(input, "/") {
			parts := strings.Fields(input)
			if len(parts) == 0 {
				continue
			}

			command := strings.TrimPrefix(parts[0], "/")

			switch command {
			case "quit", "exit":
				fmt.Println("Disconnecting and exiting...")
				clientState.CloseConnection()
				return

			case "help":
				fmt.Println("\nAvailable commands:")
				fmt.Println("  /help - Show this help")
				fmt.Println("  /quit - Exit the client")
				fmt.Println("  /stats - Show connection stats")
				fmt.Println("  /username <new_name> - Change your username")
				fmt.Println("  /join <room_id> - Join a room")
				fmt.Println("  /list groups - List available groups")
				fmt.Println("  /list guilds - List available guilds")
				fmt.Println("  /create group <name> - Create a new group")
				fmt.Println("  /create guild <name> - Create a new guild")
				fmt.Println("  /msg <user_id> <message> - Send private message")
				fmt.Println("  /ping - Send a ping to the server")
				fmt.Println("  /errors - Show recent connection errors")

			case "stats":
				fmt.Println(clientState.GetStats())

			case "username":
				if len(parts) < 2 {
					fmt.Println("Usage: /username <new_name>")
					continue
				}
				newUsername := parts[1]
				err := client.Emit("username_change", newUsername)
				if err != nil {
					fmt.Printf("Error changing username: %v\n", err)
					continue
				}
				clientState.SetUsername(newUsername)
				fmt.Printf("Username change request sent to: %s\n", newUsername)

			case "join":
				if len(parts) < 2 {
					fmt.Println("Usage: /join <room_id>")
					continue
				}
				roomID := parts[1]
				err := client.Emit("join_room", roomID)
				if err != nil {
					fmt.Printf("Error joining room: %v\n", err)
					continue
				}
				fmt.Printf("Join request sent for room: %s\n", roomID)

			case "list":
				if len(parts) < 2 {
					fmt.Println("Usage: /list groups|guilds")
					continue
				}
				roomType := parts[1]
				if roomType != "groups" && roomType != "guilds" {
					fmt.Println("Invalid room type. Use 'groups' or 'guilds'")
					continue
				}
				// Convert to singular for the server
				if roomType == "groups" {
					roomType = "group"
				} else {
					roomType = "guild"
				}
				err := client.Emit("list_rooms", roomType)
				if err != nil {
					fmt.Printf("Error listing rooms: %v\n", err)
					continue
				}
				fmt.Printf("Listing %s...\n", roomType)

			case "create":
				if len(parts) < 3 {
					fmt.Println("Usage: /create group|guild <name>")
					continue
				}
				roomType := parts[1]
				if roomType != "group" && roomType != "guild" {
					fmt.Println("Invalid room type. Use 'group' or 'guild'")
					continue
				}
				roomName := parts[2]
				err := client.Emit("create_room", roomType, roomName)
				if err != nil {
					fmt.Printf("Error creating room: %v\n", err)
					continue
				}
				fmt.Printf("Creating %s: %s\n", roomType, roomName)

			case "msg":
				if len(parts) < 3 {
					fmt.Println("Usage: /msg <user_id> <message>")
					continue
				}
				targetUserID := parts[1]
				messageText := strings.Join(parts[2:], " ")
				err := client.Emit("private_message", targetUserID, messageText)
				if err != nil {
					fmt.Printf("Error sending private message: %v\n", err)
					continue
				}
				fmt.Printf("Private message sent to %s\n", targetUserID)

			case "ping":
				err := client.Emit("ping", "Ping from client")
				if err != nil {
					fmt.Printf("Error sending ping: %v\n", err)
					continue
				}
				fmt.Println("Ping sent to server")

			case "errors":
				errors := clientState.GetConnectionErrors()
				if len(errors) == 0 {
					fmt.Println("No connection errors recorded")
				} else {
					fmt.Println("Recent connection errors:")
					for i, err := range errors {
						fmt.Printf("%d. %s\n", i+1, err)
					}
				}

			default:
				fmt.Printf("Unknown command: %s. Type /help for available commands.\n", command)
			}
		} else if input != "" {
			// Not a command, send as a chat message to current room
			currentRoom := clientState.GetCurrentRoom()
			var err error

			if currentRoom == "" || currentRoom == "global" {
				// Send to global chat
				err = client.Emit("global_message", input)
				if err != nil {
					fmt.Printf("Error sending message: %v\n", err)
				} else {
					clientState.TrackMessageSent()
				}
			} else {
				// Check if it's a group or guild (simplified - you might want to improve this)
				if strings.HasPrefix(currentRoom, "demo-guild") {
					err = client.Emit("guild_message", currentRoom, input)
				} else {
					err = client.Emit("group_message", currentRoom, input)
				}

				if err != nil {
					fmt.Printf("Error sending message: %v\n", err)
				} else {
					clientState.TrackMessageSent()
				}
			}
		}

		fmt.Print("> ")
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading input: %v", err)
	}

	// Close connection before exiting
	clientState.CloseConnection()
}
