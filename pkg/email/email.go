package email

import (
	"errors"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendEmail Send mail
func _(subject, toAddress, plainTextContent, htmlContent string) error {
	from := mail.NewEmail(os.Getenv("SENDGRID_FROM_NAME"), os.Getenv("SENDGRID_FROM_ADDRESS"))
	to := mail.NewEmail(os.Getenv("SENDGRID_TO_NAME"), toAddress)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		return err
	}
	if response.StatusCode != 202 {
		return errors.New(response.Body)
	}

	return nil
}
