package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"task-manager/constants"
	"task-manager/dto"
)

func SendSummaryEmailToUser(templatePath string, emailData *dto.EmailCronTemplateData) error {

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, emailData)
	if err != nil {
		return err
	}

	emailBody := buf.String()

	return sendEmail(emailBody, os.Getenv(constants.FromEmail), emailData.UserEmail, constants.EmailSummarySubject)
}

func sendEmail(emailBody, fromEmail, toEmail, subject string) error {
	auth := smtp.PlainAuth(
		"",
		fromEmail,
		os.Getenv(constants.EmailHostPassword),
		os.Getenv(constants.EmailHost),
	)

	headers := make(map[string]string)
	headers["From"] = fromEmail
	headers["To"] = toEmail
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"utf-8\""

	emailCompleteBody := ""
	for k, v := range headers {
		emailCompleteBody += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	emailCompleteBody += "\r\n" + emailBody
	err := smtp.SendMail(
		os.Getenv(constants.EmailHost),
		auth,
		fromEmail,
		[]string{toEmail},
		[]byte(emailCompleteBody),
	)
	if err != nil {
		return err
	}

	return nil
}
