//go:build js && wasm

package email

import (
	"io"
	"net/http"
	"syscall/js"

	jsclass "github.com/Darckfast/workers-go/internal/class"
	jshttp "github.com/Darckfast/workers-go/internal/http"
	jsstream "github.com/Darckfast/workers-go/internal/stream"
)

type EmailMessage struct {
	From string
	To   string
	Raw  io.ReadCloser
}

type ForwardableEmailMessage struct {
	From    string
	To      string
	Headers http.Header
	Raw     io.ReadCloser
	RawSize int
	self    *js.Value
}

func (f *ForwardableEmailMessage) SetReject(reason string) {
	f.self.Call("setReject", reason)
}

func (f *ForwardableEmailMessage) Forward(rcptTo string, headers *http.Header) error {
	promise := f.self.Call("forward", rcptTo, jshttp.ToJSHeader(*headers))

	_, err := jsclass.Await(promise)

	return err
}

func (f *ForwardableEmailMessage) Reply(emailMsg *EmailMessage) error {
	emlMsgObj := jsclass.Object.New()
	emlMsgObj.Set("from", emailMsg.From)
	emlMsgObj.Set("to", emailMsg.To)
	readableStream := jsstream.ReadCloserToReadableStream(emailMsg.Raw)
	emlMsgObj.Set("raw", readableStream)

	promise := f.self.Call("reply", emlMsgObj)

	_, err := jsclass.Await(promise)

	return err
}

func NewForwardableEmailMessage(msg js.Value) *ForwardableEmailMessage {
	h, _ := jshttp.ToHeader(msg.Get("headers"))
	frwMsg := ForwardableEmailMessage{
		From:    msg.Get("from").String(),
		To:      msg.Get("to").String(),
		Headers: h,
		Raw:     jshttp.ToBody(msg.Get("raw")),
		RawSize: msg.Get("rawSize").Int(),
		self:    &msg,
	}

	return &frwMsg
}
