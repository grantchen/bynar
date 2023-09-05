package email

import (
	"errors"
	"fmt"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/config"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// Send mail
func SendEmail(subject, toAddress, plainTextContent, htmlContent string) error {
	sendgridConfig := config.GetSendgridConfig()
	from := mail.NewEmail(sendgridConfig.FromName, sendgridConfig.FromAddress)
	to := mail.NewEmail(sendgridConfig.ToName, toAddress)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(sendgridConfig.ApiKey)
	response, err := client.Send(message)
	if err != nil {
		return err
	}
	if response.StatusCode != 202 {
		return errors.New(response.Body)
	}

	return nil
}

// Send mail with otp
func SendOtpEmail(toAddress string) (string, error) {
	verificationCode := GenerateTotpToken(toAddress)
	subject := "Your Security code is %s"
	subject = fmt.Sprintf(subject, verificationCode)
	htmlContent := `
		<p style="font-size: 20px;">Your Security code is</p>
		<strong style="font-size: 30px;">%s</strong>
	`
	htmlContent = fmt.Sprintf(htmlContent, verificationCode)
	if err := SendEmail(subject, toAddress, "", htmlContent); err != nil {
		return verificationCode, err
	}

	return verificationCode, nil
}

// Send mail with magic link
func SendMagicLinkEmail(toAddress string) (string, error) {
	verificationCode := GenerateTotpToken(toAddress)
	subject := "Please verify your email"
	htmlContent := `
		<p style="font-size: 20px;">Please verify your email</p>
		<a style="box-sizing: border-box; border-color: #348eda; font-weight: 400; text-decoration: none; 
			display: inline-block; margin: 0; color: #ffffff; background-color: #348eda; border: solid 1px #348eda; 
			border-radius: 2px; cursor: pointer; font-size: 14px; padding: 12px 45px;" href="%s"
		>
			Verify Email
		</a>
	`
	href := config.GetSendgridConfig().MagicLinkUrl + "?verificationCode=" + verificationCode
	htmlContent = fmt.Sprintf(htmlContent, href)
	if err := SendEmail(subject, toAddress, "", htmlContent); err != nil {
		return verificationCode, err
	}

	return verificationCode, nil
}
