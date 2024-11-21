package utils

import (
    "crypto/rand"
    "crypto/sha256"
    "encoding/hex"
    "errors"
)

func GenerateRandomString(length int) (string, error) {
    if length <= 0 {
        return "", errors.New("length must be greater than 0")
    }

    bytes := make([]byte, length)
    _, err := rand.Read(bytes)
    if err != nil {
        return "", err
    }
    return hex.EncodeToString(bytes), nil
}

func HashString(input string) string {
    hash := sha256.Sum256([]byte(input))
    return hex.EncodeToString(hash[:])
}