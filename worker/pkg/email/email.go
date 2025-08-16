package emailhandler

import (
	"fmt"
	"io"
	"strings"

	"github.com/syumai/workers/cloudflare/email"
	jsemail "github.com/syumai/workers/internal/email"
)

func New() {
	email.ConsumeNonBlock(func(f *jsemail.ForwardableEmailMessage) error {
		f.Headers.Add("x-test-id", "12345-asdfg-56789-ghjkl")

		err := f.Forward("<YOUR_EMAIL>", &f.Headers)

		if err != nil {
			return err
		}

		emailBody := strings.NewReader(fmt.Sprintf("this is a test, and this email has %d", f.RawSize))
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
