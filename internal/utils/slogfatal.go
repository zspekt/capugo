package utils

import (
	"log/slog"
	"os"
)

// slog doesn't provide an equivalent to log.Fatal() so...
func SlogFatal(msg string, args ...any) {
	// lol the ...
	slog.Error(msg, args...)
	os.Exit(1)
}
