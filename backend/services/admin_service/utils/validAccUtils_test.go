package utils

import (
    "testing"
)

func TestIsValidPassword(t *testing.T) {
    tests := []struct {
        name     string
        password string
        expected bool
    }{
        {
            name:     "Test with valid password",
            password: "Password123!",
            expected: true,
        },
        {
            name:     "Test with password too short",
            password: "Pass1!",
            expected: false,
        },
        {
            name:     "Test with password missing letter",
            password: "12345678!",
            expected: false,
        },
        {
            name:     "Test with password missing digit",
            password: "Password@!",
            expected: false,
        },
        {
            name:     "Test with password missing special character",
            password: "Password123",
            expected: false,
        },
        {
            name:     "Test with password having invalid characters",
            password: "Valid1Password",  // Missing special character
            expected: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := IsValidPassword(tt.password)
            if got != tt.expected {
                t.Errorf("IsValidPassword(%v) = %v; want %v", tt.password, got, tt.expected)
            }
        })
    }
}


func TestIsValidName(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected bool
    }{
        {
            name:     "Test with valid name (length 1)",
            input:    "A",
            expected: true,
        },
        {
            name:     "Test with valid name (length 50)",
            input:    "This is a valid name that has fifty(50) characters",
            expected: true,
        },
        {
            name:     "Test with name too short (length 0)",
            input:    "",
            expected: false,
        },
        {
            name:     "Test with name too long (length 51)",
            input:    "This name is too long to be valid because it has more than fifty characters",
            expected: false,
        },
        {
            name:     "Test with valid name (length 25)",
            input:    "Valid Name with 25 characters",
            expected: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := IsValidName(tt.input)
            if got != tt.expected {
                t.Errorf("IsValidName(%v) = %v; want %v", tt.input, got, tt.expected)
            }
        })
    }
}

func TestIsAlphabetical(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected bool
    }{
        {
            name:     "Test with name containing only letters",
            input:    "John",
            expected: true,
        },
        {
            name:     "Test with name containing uppercase letters",
            input:    "JOHN",
            expected: true,
        },
        {
            name:     "Test with name containing lowercase letters",
            input:    "john",
            expected: true,
        },
        {
            name:     "Test with name containing letters and numbers",
            input:    "John123",
            expected: false,
        },
        {
            name:     "Test with name containing special characters",
            input:    "John@Doe",
            expected: false,
        },
        {
            name:     "Test with empty name",
            input:    "",
            expected: false,
        },
        {
            name:     "Test with name containing spaces",
            input:    "John Doe",
            expected: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := IsAlphabetical(tt.input)
            if got != tt.expected {
                t.Errorf("IsAlphabetical(%v) = %v; want %v", tt.input, got, tt.expected)
            }
        })
    }
}

func TestIsValidUsername(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected bool
    }{
        {
            name:     "Valid username with letters",
            input:    "JohnDoe",
            expected: true,
        },
        {
            name:     "Valid username with letters and numbers",
            input:    "JohnDoe123",
            expected: true,
        },
        {
            name:     "Valid username with hyphen",
            input:    "John-Doe",
            expected: true,
        },
        {
            name:     "Too short username",
            input:    "Jo",
            expected: false,
        },
        {
            name:     "Too long username",
            input:    "ThisIsAUsernameThatIsWayTooLong",
            expected: false,
        },
        {
            name:     "Username with special characters",
            input:    "John@Doe",
            expected: false,
        },
        {
            name:     "Username with spaces",
            input:    "John Doe",
            expected: false,
        },
        {
            name:     "Empty username",
            input:    "",
            expected: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := IsValidUsername(tt.input)
            if got != tt.expected {
                t.Errorf("IsValidUsername(%v) = %v; want %v", tt.input, got, tt.expected)
            }
        })
    }
}


func TestIsValidPhoneNumber(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected bool
    }{
        // Số hợp lệ
        {
            name:     "Valid phone number with +84",
            input:    "+84901234567",
            expected: true,
        },
        {
            name:     "Valid phone number with 0",
            input:    "0901234567",
            expected: false,
        },
        {
            name:     "Valid phone number with 10 digits",
            input:    "+84345678901",
            expected: true,
        },
        // Số không hợp lệ
        {
            name:     "Phone number too short",
            input:    "+849012345",
            expected: false,
        },
        {
            name:     "Phone number too long",
            input:    "+8490123456789",
            expected: false,
        },
        {
            name:     "Phone number with invalid prefix",
            input:    "+84123456789",
            expected: false,
        },
        {
            name:     "Phone number with invalid characters",
            input:    "+84901234abc",
            expected: false,
        },
        // Số không hợp lệ với prefix không đúng
        {
            name:     "Phone number with invalid + prefix",
            input:    "+85123456789",
            expected: false,
        },
        {
            name:     "Phone number without leading 0 or +84",
            input:    "391234567",
            expected: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := IsValidPhoneNumber(tt.input)
            if got != tt.expected {
                t.Errorf("IsValidPhoneNumber(%v) = %v; want %v", tt.input, got, tt.expected)
            }
        })
    }
}