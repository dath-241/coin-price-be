package utils

import (
    "regexp"
	"unicode"
	"strings"
)

// Hàm kiểm tra định dạng mật khẩu
func IsValidPassword(password string) bool {
    if len(password) < 8 {
        return false
    }

    hasLetter := false
    hasDigit := false
    hasSpecial := false
    specialChars := ".,/^@$!%*?&"


    for _, char := range password {
        if unicode.IsLetter(char) {
            hasLetter = true
        } else if unicode.IsDigit(char) {
            hasDigit = true
        } else if strings.ContainsRune(specialChars, char) {
            hasSpecial = true
        }
    }


    return hasLetter && hasDigit && hasSpecial
}
// Hàm kiểm tra độ dài của tên
func IsValidName(name string) bool {
    return len(name) >= 1 && len(name) <= 50
}
    
// Hàm kiểm tra xem tên có chỉ chứa các ký tự chữ cái không
func IsAlphabetical(name string) bool {
    re := regexp.MustCompile(`^[A-Za-z]+$`)
    return re.MatchString(name)
}

func IsValidUsername(username string) bool {
	// Regex chỉ cho phép ký tự alphanumeric và dấu gạch ngang
	const usernameRegex = `^[a-zA-Z0-9-]{3,20}$`
	matched, _ := regexp.MatchString(usernameRegex, username)
	return matched
}

// Kiểm tra định dạng số điện thoại Việt Nam
func IsValidPhoneNumber(phoneNumber string) bool {
    // Regex cho số điện thoại Việt Nam: +84 hoặc 0, theo sau là 9-10 chữ số
    regex := `^(?:\+84|0)(?:3[2-9]|5[6|8|9]|7[0|6-9]|8[1-5]|9[0-9])[0-9]{7}$`
    matched, _ := regexp.MatchString(regex, phoneNumber)
    return matched
}
