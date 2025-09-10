//go:build js && wasm

package email

import (
	"errors"
	"log"
	"syscall/js"

	jsclass "github.com/Darckfast/workers-go/internal/class"
	"github.com/Darckfast/workers-go/platform/cloudflare/env"
	"github.com/Darckfast/workers-go/platform/cloudflare/lifecycle"
)

type EmailConsumer func(f *ForwardableEmailMessage) error

var consumer EmailConsumer = func(_ *ForwardableEmailMessage) error {
	return errors.New("no consumer implemented")
}

func init() {
	var handleRequestPromise = js.FuncOf(func(this js.Value, args []js.Value) any {
		fwrMsgObj := args[0]
		envObj := args[1]
		ctxObj := args[2]
		var cb js.Func
		cb = js.FuncOf(func(_ js.Value, pArgs []js.Value) any {
			defer cb.Release()
			resolve := pArgs[0]
			reject := pArgs[1]

			go func() {
				err := handler(fwrMsgObj, envObj, ctxObj)
				if err != nil {
					reject.Invoke(jsclass.ToJSError(err))
				} else {
					resolve.Invoke(true)
				}
			}()

			return nil
		})

		return jsclass.Promise.New(cb)
	})

	js.Global().Get("cf").Set("email", handleRequestPromise)
}

func handler(emailObj, envObj, ctxObj js.Value) error {
	lifecycle.Env = envObj
	lifecycle.Ctx = jsclass.ExecutionContextWrap{Ctx: ctxObj}

	err := env.LoadEnvs()
	if err != nil {
		return err
	}
	email := NewForwardableEmailMessage(emailObj)
	defer func() {
		err := email.Raw.Close()
		if err != nil {
			log.Println("error closing email raw body reader", err.Error())
		}
	}()

	return consumer(email)
}

func ConsumeNonBlock(c EmailConsumer) {
	consumer = c
}
