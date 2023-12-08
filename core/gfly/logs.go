package gfly

import (
	"app/core/log"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// ===========================================================================================================
// 										Logs
// ===========================================================================================================

func setupLog() {
	logChannel := os.Getenv("LOG_CHANNEL")

	// Log channel file
	if logChannel == "file" {
		logFile := fmt.Sprintf("storage/logs/%s", os.Getenv("LOG_FILE"))

		// Set the output destination to the console and file.
		file, _ := os.OpenFile(filepath.Clean(logFile), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
		iw := io.MultiWriter(os.Stdout, file)
		log.SetOutput(iw)
	}

	// Set log level
	switch os.Getenv("LOG_LEVEL") {
	case "Trace":
		log.SetLevel(log.LevelTrace)
	case "Debug":
		log.SetLevel(log.LevelDebug)
	case "Info":
		log.SetLevel(log.LevelInfo)
	case "Warn":
		log.SetLevel(log.LevelWarn)
	case "Error":
		log.SetLevel(log.LevelError)
	case "Fatal":
		log.SetLevel(log.LevelFatal)
	case "Panic":
		log.SetLevel(log.LevelPanic)
	}

	log.Trace("Setup Logs")
}
