package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"net/smtp"
	"os"
	"regexp"
	"strings"
)

var disposableList = make(map[string]struct{}, 3500)

var disposableErrorMessage = errors.New("disposable mail")

func init() {
	f, _ := os.Open("disposable_email_blocklist.txt")
	for scanner := bufio.NewScanner(f); scanner.Scan(); {
		disposableList[scanner.Text()] = struct{}{}
	}
	f.Close()
}

func EmailVerificationProcess(email string) (*EmailVerifier, error) {
	emailVer := NewEmailVerifier(email)

	if err := emailVer.EmailFormatVerifier(); err != nil {
		return emailVer, err
	}

	emailVer.EmailParser()

	if emailVer.isDisposableEmail() {
		return emailVer, disposableErrorMessage
	}

	_, err := emailVer.EmailDomainVerifier()
	if err != nil {
		return emailVer, err
	}

	fmt.Println("here")
	// err = emailVer.EmailSMPTVerifier(mx.Host)
	// fmt.Println(err)
	// if err != nil {
	// 	return emailVer, err
	// }

	return emailVer, nil
}

func (ev *EmailVerifier) isDisposableEmail() (disposable bool) {
	_, disposable = disposableList[strings.ToLower(ev.Domain)]
	return
}

func (ev *EmailVerifier) EmailFormatVerifier() error {
	if len(ev.Email) > 254 {
		return fmt.Errorf("invalid email pattern")
	}

	pattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(pattern)
	if !re.MatchString(ev.Email) {
		return fmt.Errorf("invalid email pattern")
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

func (ev *EmailVerifier) NewEmailStatus(err error) *EmailStatus {
	return &EmailStatus{
		Email:      ev.Email,
		Domain:     ev.Domain,
		UserName:   ev.UserName,
		Disposable: err == disposableErrorMessage,
		Valid:      err == nil,
		Error:      err.Error(),
	}
}
