package emailhandler

import (
	"io"
	"strconv"
	"strings"

	"github.com/Darckfast/workers-go/cloudflare/email"
	jsemail "github.com/Darckfast/workers-go/internal/email"
)

func New() {
	/*
	 * This functions must be called to instantiate a email handler consumer
	 */
	email.ConsumeNonBlock(func(f *jsemail.ForwardableEmailMessage) error {
		f.Headers.Add("x-test-id", "12345-asdfg-56789-ghjkl")

		err := f.Forward("<YOUR_EMAIL>", &f.Headers)

		if err != nil {
			return err
		}

		emailBody := strings.NewReader("this is a test, and this email has " + strconv.Itoa(f.RawSize))
		reply := jsemail.EmailMessage{
			From: "me",
			To:   "you",
			Raw:  io.NopCloser(emailBody),
		}
		err = f.Reply(&reply)

		if err != nil {
			return err
		}

		f.SetReject("this reject is just for testing")
		return nil
	})
}
