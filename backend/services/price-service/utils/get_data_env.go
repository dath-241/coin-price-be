package utils

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

// func get key api of coin market cap
func GetKeyApi() (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", errors.New("Error from loading env variable")
	}
	key := os.Getenv("API_KEY_COINMARKET_CAP")
	return key, nil
}
