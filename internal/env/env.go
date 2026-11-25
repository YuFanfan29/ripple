// Package env reads configuration from environment variables.
package env

import "os"

// String returns the value of the environment variable key, or fallback if
// the variable is unset or empty.
func String(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
