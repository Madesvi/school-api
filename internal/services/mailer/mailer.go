package mailer

import (
	"fmt"
	"net/smtp"
	"os"
)

// How to start
// mailpit - in the terminal

func SendWelcomeEmail(to string, message string) error {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	from := os.Getenv("SMTP_SENDER")

	// I will use map for headers
	header := make(map[string]string)
	header["From"] = from
	header["To"] = to
	header["Subject"] = "Password reset"
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=\"UTF-8\""

	headerStr := ""
	for k, v := range header {
		headerStr += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	body := fmt.Sprintf("<html><body><h1>Hello!</h1><p><b>%s</b></p></body></html>", message)

	msg := headerStr + "\r\n" + body

	// В Mailpit аутентификация не нужна, поэтому Auth = nil
	addr := fmt.Sprintf("%s:%s", host, port)
	return smtp.SendMail(addr, nil, from, []string{to}, []byte(msg))
}
