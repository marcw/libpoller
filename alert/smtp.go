package alert

import (
	"fmt"
	"github.com/marcw/ezmail"
	"github.com/marcw/poller"
	"net"
	"net/smtp"
	"os"
	"strings"
	"time"
)

type smtpAlerter struct {
	addr    string
	auth    smtp.Auth
	message ezmail.Message
}

func NewSmtpAlerter() (poller.Alerter, error) {
	envHost := os.Getenv("SMTP_HOST")
	envPort := os.Getenv("SMTP_PORT")
	envAuth := os.Getenv("SMTP_AUTH")
	envUsername := os.Getenv("SMTP_USERNAME")
	envPassword := os.Getenv("SMTP_PASSWORD")
	envIdentity := os.Getenv("SMTP_PLAIN_IDENTITY")
	envTo := os.Getenv("SMTP_RECIPIENT")
	envFrom := os.Getenv("SMTP_FROM")

	if envHost == "" {
		return nil, fmt.Errorf("Please define SMTP_HOST env var")
	}
	if envPort == "" {
		return nil, fmt.Errorf("Please define SMTP_PORT env var")
	}

	if envAuth != "" && envAuth != "MD5" && envAuth != "PLAIN" {
		return nil, fmt.Errorf("Please either leave SMTP_AUTH env empty or set it to MD5 or PLAIN")
	}

	var addr string
	if envPort != "" {
		addr = net.JoinHostPort(envHost, envPort)
	} else {
		addr = envHost
	}

	var auth smtp.Auth
	if envAuth == "MD5" {
		auth = smtp.CRAMMD5Auth(envUsername, envPassword)
	} else if envAuth == "PLAIN" {
		auth = smtp.PlainAuth(envIdentity, envUsername, envPassword, envHost)
	}

	message := ezmail.NewMessage()
	message.SetFrom("", envFrom)
	for _, v := range strings.Split(envTo, ";") {
		message.AddTo("", v)
	}

	smtp := &smtpAlerter{
		addr:    addr,
		auth:    auth,
		message: *message}

	return smtp, nil
}

func (m *smtpAlerter) Alert(event *poller.Event) {
	msg := m.message
	msg.Subject = fmt.Sprintf("[ALERT] %s is down", event.Check.Url.String())
	msg.Body = fmt.Sprintf("Poller alert: %s (%s) is down since %s", event.Check.Key, event.Check.Url.String(), event.Check.DownSince.Format(time.RFC822))

	if err := smtp.SendMail(m.addr, m.auth, msg.From.String(), msg.Recipients(), msg.Bytes()); err != nil {
		println(err)
	}
}
