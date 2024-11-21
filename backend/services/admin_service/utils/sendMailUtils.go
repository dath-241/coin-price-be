package utils

import (
	"fmt"
	"os"
    "encoding/json"
    "log"

    "github.com/mailjet/mailjet-apiv3-go/v3"
)


// SendEmail gửi email đến người nhận.
func SendEmail(to, subject, htmlBody string) error {//, textBody
    if to == "" {
        return fmt.Errorf("recipient email is required")
    }
    
    // Khởi tạo Mailjet client
    
    apiKey := os.Getenv("MAILJET_API_KEY")
    secretKey := os.Getenv("MAILJET_SECRET_KEY")
    senderEmail := os.Getenv("EMAIL_SENDER")
    if senderEmail == "" {
        log.Fatal("EMAIL_SENDER not set in environment variables")
    }
    if apiKey == "" || secretKey == "" {
    log.Fatal("Mailjet API keys are not set in environment variables")
    }
    //fmt.Println("Generated Email Body:", htmlBody)

    mailjetClient := mailjet.NewMailjetClient(apiKey, secretKey)

    // Tạo thông điệp email
    emailData := &mailjet.MessagesV31{
        Info: []mailjet.InfoMessagesV31{
            {
                From: &mailjet.RecipientV31{
                    Email: senderEmail,
                    Name:  "Coin-Price",
                },
                To: &mailjet.RecipientsV31{ 
                    { 
                        Email: to, // Gán địa chỉ email của người nhận 
                        Name: "", // Bạn có thể thêm tên người nhận nếu muốn 
                    },
                },
                Subject:  subject,
                //TextPart: textBody,  // Nội dung văn bản
                HTMLPart: htmlBody,  // Nội dung HTML
            },
        },
        SandBoxMode: false, // Chế độ sandbox, chỉnh lại nếu cần
    }

    // Gửi email
    res, err := mailjetClient.SendMailV31(emailData)
    if err != nil {
        return err
    }

    // Kiểm tra phản hồi
    responseData, _ := json.Marshal(res)
    fmt.Printf("Email đã được gửi thành công: %s\n", responseData)
    return nil
}