package poller

import (
	"fmt"
	"github.com/marcw/ezmail"
	"github.com/marcw/pagerduty"
	"net"
	"net/smtp"
	"os"
	"strings"
	"time"
)

// SMTP
type smtpAlerter struct {
	addr    string
	auth    smtp.Auth
	message ezmail.Message
}

func NewSmtpAlerter() (Alerter, error) {
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

func (m *smtpAlerter) Alert(event *Event) {
	msg := m.message
	msg.Subject = fmt.Sprintf("[ALERT] %s is down", event.Check.Key)
	msg.Body = fmt.Sprintf("Poller alert: %s is down since %s", event.Check.AlertDescription(), event.Check.DownSince.Format(time.RFC822))

	if err := smtp.SendMail(m.addr, m.auth, msg.From.String(), msg.Recipients(), msg.Bytes()); err != nil {
		println(err)
	}
}

// PagerDuty
type pagerDutyAlerter struct {
	serviceKey string
}

func NewPagerDutyAlerter() (Alerter, error) {
	envServiceKey := os.Getenv("PAGERDUTY_SERVICE_KEY")
	if envServiceKey == "" {
		return nil, fmt.Errorf("Please define the PAGERDUTY_SERVICE_KEY environment variable.")
	}

	return &pagerDutyAlerter{envServiceKey}, nil
}

func (pda *pagerDutyAlerter) Alert(event *Event) {
	description := fmt.Sprintf("%s is DOWN since %s.", event.Check.AlertDescription(), event.Check.DownSince.Format(time.RFC3339))
	e := pagerduty.NewTriggerEvent(pda.serviceKey, description)
	e.Details["checked_at"] = event.Time.Format(time.RFC3339)
	e.Details["duration"] = event.Duration.String()
	e.Details["status_code"] = event.StatusCode
	e.Details["was_up_for"] = event.Check.WasUpFor.String()
	e.IncidentKey = event.Check.Key
	for {
		_, statusCode, _ := pagerduty.Submit(e)
		if statusCode < 500 {
			break
		} else {
			// Wait a bit before trying again
			time.Sleep(3 * time.Second)
		}
	}
}
