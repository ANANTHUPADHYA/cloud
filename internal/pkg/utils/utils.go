package utils

import (
	"os"

	"github.com/twinj/uuid"
)

// GetEnvOrDefault will return the environmental value of given variable
// if not present it'll return the default value
func GetEnvOrDefault(envVar string, defaultValue string) string {
	if v := os.Getenv(envVar); v != "" {
		return v
	}
	return defaultValue
}

func GenerateUUID() string {
	return uuid.NewV4().String()
}

// DefaultOkCodes are the default http ok codes
func DefaultOkCodes(method string) []int {
	switch {
	case method == "GET":
		return []int{200}
	case method == "POST":
		return []int{200, 201, 202}
	case method == "PUT":
		return []int{200, 201, 202, 204}
	case method == "PATCH":
		return []int{200, 204}
	case method == "DELETE":
		return []int{202, 204}
	}
	return []int{}
}
