package jsemail

import (
	"io"
	"net/http"
	"syscall/js"

	jshttp "github.com/syumai/workers/internal/http"
	jsutil "github.com/syumai/workers/internal/utils"
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

	_, err := jsutil.AwaitPromise(promise)

	return err
}

func (f *ForwardableEmailMessage) Reply(emailMsg *EmailMessage) error {
	emlMsgObj := jsutil.NewObject()
	emlMsgObj.Set("from", emailMsg.From)
	emlMsgObj.Set("to", emailMsg.To)
	readableStream := jsutil.ReadCloserToReadableStream(emailMsg.Raw)
	emlMsgObj.Set("raw", readableStream)

	promise := f.self.Call("reply", emlMsgObj)

	_, err := jsutil.AwaitPromise(promise)

	return err
}

func NewForwardableEmailMessage(msg js.Value) *ForwardableEmailMessage {
	frwMsg := ForwardableEmailMessage{
		From:    msg.Get("from").String(),
		To:      msg.Get("to").String(),
		Headers: jshttp.ToHeader(msg.Get("headers")),
		Raw:     jshttp.ToBody(msg.Get("raw")),
		RawSize: msg.Get("rawSize").Int(),
		self:    &msg,
	}

	return &frwMsg
}
