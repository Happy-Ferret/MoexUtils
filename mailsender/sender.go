package main

import (
	"flag"
	"fmt"
	"log"
	"net/smtp"
)

var (
	emailserver = "mail-server:25"
	from        = "from-notifier@server.com"
	to          = "to-recipient@server.com"
	mes         = "message"
)

func init() {
	flag.StringVar(&emailserver, "s", emailserver, "email server")
	flag.StringVar(&from, "f", from, "from rcpt")
	flag.StringVar(&to, "t", to, "to reciever")
	flag.StringVar(&mes, "m", mes, "message")
	flag.Parse()
}

func main() {

	// Connect to the remote SMTP server.
	c, err := smtp.Dial(emailserver)
	if err != nil {
		log.Fatal(err)
	}

	// Set the sender and recipient first
	if err := c.Mail(from); err != nil {
		log.Fatal(err)
	}

	// TODO добавить возможность отсылки нескольким клиентам
	if err := c.Rcpt(to); err != nil {
		log.Fatal(err)
	}

	// Send the email body.
	wc, err := c.Data()
	if err != nil {
		log.Fatal(err)
	}

	_, err = fmt.Fprintf(wc, mes)
	if err != nil {
		log.Fatal(err)
	}
	err = wc.Close()
	if err != nil {
		log.Fatal(err)
	}

	// Send the QUIT command and close the connection.
	err = c.Quit()
	if err != nil {
		log.Fatal(err)
	}
}
