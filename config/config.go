package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// init loads the .env file in the config directory into the environment, if one exists.
// This should be called automatically when the package is initialized.
func init() {
	if err := godotenv.Load("./config/.env"); err != nil {
		log.Print("No .env file found ")
	}
}

// Get retrieves the value of the environment variable named by the key.
// It returns the value and a boolean indicating whether the key was found.
func Get(key string) (string, bool) {
	return os.LookupEnv(key)
}

// MustGet retrieves the value of the environment variable named by the key.
// If the key is not found, it will panic with a message indicating the missing key.
func MustGet(key string) string {
	if value, ok := os.LookupEnv(key); !ok {
		panic("Missing environment variable: " + key)
	} else {
		return value
	}
}

func MustGetInt(key string) int {
	if value, ok := os.LookupEnv(key); !ok {
		panic("Missing environment variable: " + key)
	} else {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			panic("Environment variable is not an int: " + err.Error())
		}
		return intValue
	}
}

// DefaultGet retrieves the value of the environment variable named by the key.
// If the key is not found, it returns the provided defaultValue.
func DefaultGet(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); !ok {
		return defaultValue
	} else {
		return value
	}
}
