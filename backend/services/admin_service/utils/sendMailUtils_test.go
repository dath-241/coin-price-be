package utils

import (
    "os"
    "testing"
    "github.com/jarcoal/httpmock"
)

func TestSendEmail(t *testing.T) {
    // Thiết lập các biến môi trường cần thiết
    os.Setenv("MAILJET_API_KEY", "fake-api-key")
    os.Setenv("MAILJET_SECRET_KEY", "fake-secret-key")
    os.Setenv("EMAIL_SENDER", "test@coin-price.com")

    // Kích hoạt httpmock để thay thế các request HTTP thật
    httpmock.Activate()
    defer httpmock.DeactivateAndReset()

    // Đăng ký phản hồi giả lập cho Mailjet API
    httpmock.RegisterResponder("POST", "https://api.mailjet.com/v3.1/send",
        httpmock.NewStringResponder(200, `{"Messages": [{"Status": "success"}]}`))

    // Các trường hợp kiểm tra
    tests := []struct {
        name     string
        to       string
        subject  string
        htmlBody string
        expectErr bool
    }{
        {
            name:     "Valid Email",
            to:       "recipient@example.com",
            subject:  "Test Email",
            htmlBody: "<h1>This is a test</h1>",
            expectErr: false,
        },
        {
            name:     "Missing Recipient Email",
            to:       "",
            subject:  "Test Email",
            htmlBody: "<h1>This is a test</h1>",
            expectErr: true,
        },
    }

    // Chạy từng trường hợp kiểm tra
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := SendEmail(tt.to, tt.subject, tt.htmlBody)
            if (err != nil) != tt.expectErr {
                t.Errorf("Expected error: %v, got: %v", tt.expectErr, err)
            }
        })
    }
}
