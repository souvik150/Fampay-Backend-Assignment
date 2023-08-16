package utils

import (
	"fmt"
	initializers "github.com/souvik150/Fampay-Backend-Assignment/config"
	"net/smtp"
)

func SendEmail(recipient, subject, msg string) (string, error) {
	config, _ := initializers.LoadConfig(".")
	auth := smtp.PlainAuth(
		"",
		config.Email,
		config.EmailPassword,
		"smtp.gmail.com",
	)

	// Create the complete email message
	emailMsg := []byte(fmt.Sprintf(
		"Subject: %s\r\n"+
			"To: %s\r\n"+
			"\r\n"+
			"%s", subject, recipient, msg))

	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		config.Email,
		[]string{recipient},
		emailMsg,
	)

	if err != nil {
		fmt.Println(err)
	}

	return "", nil
}
