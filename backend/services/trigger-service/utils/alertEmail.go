package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/mailjet/mailjet-apiv3-go/v3"
)

func SendAlertEmail(to, subject, htmlBody string) error {

	apiKey := os.Getenv("MAILJET_API_KEY")
	secretKey := os.Getenv("MAILJET_SECRET_KEY")
	senderEmail := os.Getenv("EMAIL_SENDER")
	if apiKey == "" || secretKey == "" {
		log.Fatal("Mailjet API keys are not set in environment variables")
	}
	if senderEmail == "" {
		log.Fatal("EMAIL_SENDER not set in environment variables")
	}

	mailjetClient := mailjet.NewMailjetClient(apiKey, secretKey)

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
						Name:  "", // Bạn có thể thêm tên người nhận nếu muốn
					},
				},
				Subject: subject,
				//TextPart: textBody,  // Nội dung văn bản
				HTMLPart: htmlBody, // Nội dung HTML
			},
		},
		SandBoxMode: false,
	}

	res, err := mailjetClient.SendMailV31(emailData)
	if err != nil {
		return err
	}

	responseData, _ := json.Marshal(res)
	fmt.Printf("Email đã được gửi thành công: %s\n", responseData)
	return nil
}
