package client

import (
	"context"
	"fmt"
	"net/smtp"
)

type EmailClient struct {
	Host        string
	Port        int
	Username    string
	Password    string
	SenderEmail string
	SenderName  string
}

func NewEmailClient(host string, port int, username, password, senderEmail, senderName string) *EmailClient {
	return &EmailClient{
		Host:        host,
		Port:        port,
		Username:    username,
		Password:    password,
		SenderEmail: senderEmail,
		SenderName:  senderName,
	}
}

func (c *EmailClient) SendEmail(ctx context.Context, to string, orderID int) error {
	subject := fmt.Sprintf("Payment Successful - Order #%d", orderID)

	body := fmt.Sprintf("Your payment for Order #%d was successful.\n\nThank you for shopping with Fleet Ops.", orderID)

	message := []byte("From: " + c.SenderEmail + " <" + c.SenderEmail + ">\r\n" + "To: " + to + ">\r\n" + "Subject: " + subject + "\r\n" + "\r\n" + "Content-Type: text/plain; charset=UTF-8\r\n" + "\r\n" + body + "\r\n")

	auth := smtp.PlainAuth("", c.Username, c.Password, c.Host)

	addr := fmt.Sprintf("%s:%d", c.Host, c.Port)

	err := smtp.SendMail(addr, auth, c.SenderEmail, []string{to}, message)
	if err != nil {
		return fmt.Errorf("failed to send email to %s: %w", to, err)
	}

	return nil
}
