package main

import (
	"fmt"
	"net"
	"net/smtp"
	"regexp"
	"strings"
)

func EmailVerificationProcess(email string) (*EmailVerifier, error) {
	emailVer := NewEmailVerifier(email)

	if err := emailVer.EmailFormatVerifier(); err != nil {
		return emailVer, err
	}

	emailVer.EmailParser()

	mx, err := emailVer.EmailDomainVerifier()
	if err != nil {
		return emailVer, err
	}

	err = emailVer.EmailSMPTVerifier(mx.Host)
	if err != nil {
		return emailVer, err
	}

	return emailVer, nil
}

func (ev *EmailVerifier) EmailFormatVerifier() error {
	pattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(pattern)
	if !re.MatchString(ev.Email) {
		return fmt.Errorf("invalid pattern")
	}
	return nil
}

func (ev *EmailVerifier) EmailParser() {
	parts := strings.Split(ev.Email, "@")
	(*ev).UserName = parts[0]
	(*ev).Domain = parts[1]
}

func (ev *EmailVerifier) EmailDomainVerifier() (*net.MX, error) {
	mxRecords, err := net.LookupMX(ev.Domain)
	if len(mxRecords) == 0 || err != nil {
		return nil, fmt.Errorf("no mx record found")
	}

	return mxRecords[0], nil
}

func (ev *EmailVerifier) EmailSMPTVerifier(addr string) error {
	sender := GetConfig().SenderEmail

	client, err := smtp.Dial(addr + ":25")
	if err != nil {
		return err
	}

	defer client.Close()

	if err = client.Hello(ev.Domain); err != nil {
		return err
	}

	if err = client.Mail(sender); err != nil {
		return err
	}

	if err = client.Rcpt(ev.Email); err != nil {
		return err
	}

	return client.Quit()
}

func NewEmailVerifier(email string) *EmailVerifier {
	return &EmailVerifier{
		Email: email,
	}
}

func (ev *EmailVerifier) NewEmailStatus(valid bool) *EmailStatus {
	return &EmailStatus{
		Email:    ev.Email,
		Domain:   ev.Domain,
		UserName: ev.UserName,
		Valid:    valid,
	}
}
