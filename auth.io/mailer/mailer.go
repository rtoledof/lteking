package mailer

import (
	"crypto/tls"
	"log"

	"gopkg.in/gomail.v2"
)

var Messages = make(chan *gomail.Message, 10000)

type Mailer struct {
}

func NewMailer(host, user, pass string, port int) {
	d := gomail.NewDialer(host, port, user, pass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	go func() {
		for m := range Messages {
			if err := d.DialAndSend(m); err != nil {
				log.Println(err)
			}
		}
	}()
}

func GenMessage(sender, receiver, plainTemplate, htmlTemplate string) {
	m := gomail.NewMessage()
	m.SetHeader("From", sender)
	m.SetHeader("To", receiver)
	m.SetBody("text/plain", plainTemplate)
	m.AddAlternative("text/html", htmlTemplate)

	Messages <- m
}
