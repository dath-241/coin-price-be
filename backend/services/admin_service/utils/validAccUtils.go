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