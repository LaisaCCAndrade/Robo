package service

import (
	"Robo/skimas"
	"log"
	"net/smtp"
	"strings"
)

func SendEmail(Body skimas.Data) {
	username := "seu usuario do mailtrap"
	password := "sua senha do mailtrap"
	smtpHost := "sandbox.smtp.mailtrap.io"

	auth := smtp.PlainAuth("", username, password, smtpHost)

	from := "emailquerecebera@mail.com"
	to := []string{Body.Email}

	message := strings.Join([]string{
		"To: " + Body.Email,
		"From: " + from,
		"domain: " + Body.Domain,
		"Subject: Testando email",
		"MIME-Version: 1.0",
		"Content-Type: text/html; charset=\"utf-8\"",
		"",
		`<html>
		<head>
			<style>
				body {
					font-family: Inter, Arial, sans-serif;
				}
				img {
					width: 200px;
				}
			</style>
		</head>
		<body>
				<h1>E-mail</h1>
		</body>
		</html>`,
	}, "\r\n")

	smtpUrl := smtpHost + ":25"
	err := smtp.SendMail(smtpUrl, auth, from, to, []byte(message))

	if err != nil {
		log.Fatal(err)
	}

	log.Println("URL:", smtpUrl)
	log.Println("E-mail enviado para:", to)
}
