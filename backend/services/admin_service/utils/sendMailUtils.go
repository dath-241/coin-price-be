package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/mailjet/mailjet-apiv3-go/v3"
)

// SendEmail gửi email đến người nhận.
func SendEmail(to, subject, name, otp string) error {
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

	// Tạo nội dung email từ template
	htmlBody, err := GeneratePasswordResetEmailBody(name, otp)
	if err != nil {
		return fmt.Errorf("error generating email body: %v", err)
	}

	// Gửi email qua Mailjet
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
						Email: to,
						Name:  "",
					},
				},
				Subject: subject,
				HTMLPart: htmlBody,
			},
		},
		SandBoxMode: false,
	}

	// Gửi email
	_, err = mailjetClient.SendMailV31(emailData)
	if err != nil {
		return fmt.Errorf("error sending email: %v", err)
	}

	// Kiểm tra phản hồi
	//responseData, _ := json.Marshal(res)
    //fmt.Printf("Email đã được gửi thành công: %s\n", string(responseData))
	return nil
}

// GeneratePasswordResetEmailBody tạo nội dung email reset password từ template HTML trực tiếp
func GeneratePasswordResetEmailBody(name, otp string) (string, error) {
	// Định nghĩa template HTML trực tiếp trong chuỗi string
	const emailTemplate = `
	<!DOCTYPE html>
	<html lang="en">
	<head>
	  <meta charset="UTF-8">
	  <meta name="viewport" content="width=device-width, initial-scale=1.0">
	  <style>
	    body {
	      font-family: Arial, sans-serif;
	      margin: 0;
	      padding: 0;
	      background-color: #f4f4f4;
	    }
	    .container {
	      width: 100%;
	      max-width: 600px;
	      margin: 0 auto;
	      background-color: white;
	      padding: 20px;
	      border-radius: 8px;
	      box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
	    }
	    .header {
	      text-align: center;
	      color: #6a1b9a;
	      font-size: 24px;
	      font-weight: bold;
	    }
	    .otp-container {
	      text-align: center;
	      background-color: #e8f5e9;
	      padding: 20px;
	      border-radius: 8px;
	      margin: 20px 0;
	    }
	    .otp {
	      font-size: 32px;
	      font-weight: bold;
	      color: #388e3c;
	      padding: 10px 20px;
	      border-radius: 8px;
	      background-color: #e8f5e9;
	      display: inline-block;
	    }
	    .message {
	      text-align: center;
	      font-size: 16px;
	      color: #555;
	      margin-top: 20px;
	    }
	    .footer {
	      text-align: center;
	      font-size: 12px;
	      color: #888;
	      margin-top: 30px;
	    }
	    .footer a {
	      color: #888;
	      text-decoration: none;
	    }
	  </style>
	</head>
	<body>

	<div class="container">
	  <div class="header">
	    Password Reset Request
	  </div>

	  <p>Dear {{.Name}},</p>
	  <p>You have requested to reset your password. Please use the following OTP to complete the process:</p>

	  <div class="otp-container">
	    <span class="otp">{{.OTP}}</span>
	  </div>

	  <p class="message">
	    This OTP is valid until {{.ExpirationTime}} minutes. Do not share this OTP with anyone.
	  </p>

	  <p class="footer">
	    If you did not request a password reset, please ignore this email or contact support.
	  </p>

	  <p class="footer">
	    <a href="#">Privacy Policy</a> | <a href="#">Terms of Service</a>
	  </p>
	</div>

	</body>
	</html>
	`

	// Phân tích template HTML và điền dữ liệu vào
	t, err := template.New("OTP-email").Parse(emailTemplate)
	if err != nil {
		return "", fmt.Errorf("error parsing email template: %v", err)
	}

	var bodyBuffer bytes.Buffer
	err = t.Execute(&bodyBuffer, map[string]interface{}{
		"Name": name,
		"OTP":  otp,
        "ExpirationTime": 5,
	})
	if err != nil {
		return "", fmt.Errorf("error executing template: %v", err)
	}

	// Trả về nội dung email đã điền dữ liệu
	return bodyBuffer.String(), nil
}