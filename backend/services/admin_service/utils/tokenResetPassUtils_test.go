package utils

import (
    "testing"
)

func TestGenerateRandomString(t *testing.T) {
    tests := []struct {
        name     string
        length   int
        expected string
        shouldErr bool
    }{
        {
            name:     "Valid length",
            length:   16,
            expected: "expected valid random string length",
            shouldErr: false,
        },
        {
            name:     "Length 0",
            length:   0,
            expected: "",
            shouldErr: true,
        },
        {
            name:     "Negative length",
            length:   -1,
            expected: "",
            shouldErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := GenerateRandomString(tt.length)

            if tt.shouldErr && err == nil {
                t.Errorf("Expected error, but got nil")
            } else if !tt.shouldErr && err != nil {
                t.Errorf("Expected no error, but got %v", err)
            }

            // For valid lengths, check the expected length of the hex string
            if !tt.shouldErr {
                expectedLength := tt.length * 2
                if len(result) != expectedLength {
                    t.Errorf("Expected length %d, but got %d", expectedLength, len(result))
                }
            }

            // Log the result and check if it matches the expected condition
            t.Logf("Test: %s | Got: %s | Expected Length: %d | Expected Error: %v", tt.name, result, len(result), tt.shouldErr)
        })
    }
}

func TestHashString(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {
            name:     "Test with empty string",
            input:    "",
            expected: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", // SHA256 của chuỗi rỗng
        },
        {
            name:     "Test with regular string",
            input:    "hello",
            expected: "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824", // SHA256 của "hello"
        },
        {
            name:     "Test with different string",
            input:    "world",
            expected: "486ea46224d1bb4fb680f34f7c9ad96a8f24ec88be73ea8e5a6c65260e9cb8a7", // SHA256 của "world"
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := HashString(tt.input)

			// In ra kết quả expected và got
            t.Logf("Test: %s | Input: %s | Expected: %s | Got: %s", tt.name, tt.input, tt.expected, got)

            if got != tt.expected {
                t.Errorf("HashString(%v) = %v; want %v", tt.input, got, tt.expected)
            }
        })
    }
}