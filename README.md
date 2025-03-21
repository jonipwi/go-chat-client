# Go Socket.IO Chat Client

## Project Structure

```
go-chat-client/
│
├── main.go                 # Main application entry point
├── server_connection.go    # Server connection and heartbeat logic
├── client_state.go         # Client state management
│
├── commands/
│   ├── commands.go         # Core command processing interface
│   ├── messaging_commands.go  # Message-related commands
│   ├── room_commands.go    # Room-related commands
│   └── utility_commands.go # Utility commands like ping, stats
│
├── events/
│   └── event_handlers.go   # Socket event listeners and handlers
│
├── utils/
│   ├── logger.go           # Logging utility
│   └── helpers.go          # Utility helper functions
│
└── go.mod                  # Go module file
```

## Features

- Real-time chat communication using Socket.IO
- Multiple chat rooms (global, group, guild)
- Private messaging
- Connection management and heartbeat
- Logging and error tracking

## Prerequisites

- Go 1.21+
- Socket.IO server running

## Installation

1. Clone the repository
2. Run `go mod tidy`
3. Configure server connection in `main.go`
4. Run `go run .`

## Commands

- `/global <message>`: Send a global message
- `/group <group_id> <message>`: Send a group message
- `/guild <guild_id> <message>`: Send a guild message
- `/private <user_id> <message>`: Send a private message
- `/create <type> <name>`: Create a new room
- `/join <room_id>`: Join a room
- `/list <type>`: List available rooms
- `/username <new_name>`: Change username
- `/stats`: Show connection statistics
- `/debug`: Show connection debug info
- `/exit`: Disconnect and exit

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## License

[Your License Here]