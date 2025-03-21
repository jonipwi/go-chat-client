package utils

import (
	"log"
	"os"
)

// Logger is a custom logger with timestamp and file information
var Logger = log.New(os.Stdout, "[CHAT-CLIENT] ", log.LstdFlags|log.Lshortfile)
