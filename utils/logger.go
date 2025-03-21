package utils

import (
	"io"
	"log"
	"os"
)

// Logger is a custom logger with timestamp and file information
var Logger *log.Logger

func init() {
	// Create a file for logging
	file, err := os.OpenFile("chat_client.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}

	// Create a multi-writer that writes to both file and stdout
	multiWriter := io.MultiWriter(os.Stdout, file)

	// Initialize the logger with the multi-writer
	Logger = log.New(multiWriter, "[CHAT-CLIENT] ", log.LstdFlags|log.Lshortfile)
}
