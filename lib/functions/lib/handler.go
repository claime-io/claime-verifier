package lib

import (
	"os"
)

// Headers with headers
func Headers() map[string]string {
	return map[string]string{
		"Access-Control-Allow-Headers":     "*",
		"Access-Control-Allow-Methods":     "GET,POST,PUT,DELETE",
		"Access-Control-Allow-Credentials": "true",
		"Access-Control-Allow-Origin":      os.Getenv("AllowedOrigin"),
	}
}
