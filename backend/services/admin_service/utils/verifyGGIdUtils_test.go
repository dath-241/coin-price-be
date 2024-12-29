package utils

import (
    "testing"
    "github.com/jarcoal/httpmock"
)

// TestVerifyGoogleIDToken: Kiểm tra hàm VerifyGoogleIDToken để đảm bảo xử lý đúng các token hợp lệ và không hợp lệ.
func TestVerifyGoogleIDToken(t *testing.T) {
    // Kích hoạt httpmock để thay thế mọi request HTTP thật bằng các phản hồi giả lập.
    httpmock.Activate()
    // Đảm bảo httpmock sẽ được tắt và reset sau khi kết thúc test.
    defer httpmock.DeactivateAndReset()

    // Đăng ký phản hồi giả lập cho token hợp lệ.
    httpmock.RegisterResponder("GET", "https://oauth2.googleapis.com/tokeninfo?id_token=validToken",
        httpmock.NewStringResponder(200, `{"email": "test@example.com", "name": "Test User"}`))

    // Đăng ký phản hồi giả lập cho token không hợp lệ.
    httpmock.RegisterResponder("GET", "https://oauth2.googleapis.com/tokeninfo?id_token=invalidToken",
        httpmock.NewStringResponder(401, ``))

    // Danh sách các trường hợp kiểm tra.
    tests := []struct {
        name     string                 // Tên của test case, để dễ nhận biết khi chạy test.
        idToken  string                 // Token được dùng để kiểm tra.
        expected map[string]interface{} // Kết quả mong muốn nếu không có lỗi.
        hasError bool                   // Kỳ vọng có lỗi hay không.
    }{
        {
            name:    "Valid Token", // Trường hợp token hợp lệ.
            idToken: "validToken",
            expected: map[string]interface{}{
                "email": "test@example.com",
                "name":  "Test User",
            },
            hasError: false, // Không kỳ vọng lỗi.
        },
        {
            name:     "Invalid Token", // Trường hợp token không hợp lệ.
            idToken:  "invalidToken",
            expected: nil,             // Kết quả mong muốn là nil.
            hasError: true,            // Kỳ vọng lỗi.
        },
    }

    // Lặp qua từng test case để kiểm tra.
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Gọi hàm VerifyGoogleIDToken với token hiện tại.
            result, err := VerifyGoogleIDToken(tt.idToken)
            
            // Kiểm tra xem trạng thái lỗi có đúng như kỳ vọng không.
            if (err != nil) != tt.hasError {
                t.Errorf("Expected error: %v, got: %v", tt.hasError, err)
            }

            // Nếu không có lỗi, kiểm tra xem kết quả trả về có đúng như mong muốn không.
            if !tt.hasError && result["email"] != tt.expected["email"] {
                t.Errorf("Expected: %v, got: %v", tt.expected, result)
            }
        })
    }
}
