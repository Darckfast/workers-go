//go:build js && wasm

package email

import (
	"errors"
	"syscall/js"

	"github.com/syumai/workers/cloudflare/env"
	_ "github.com/syumai/workers/cloudflare/env"
	jsclass "github.com/syumai/workers/internal/class"
	jsemail "github.com/syumai/workers/internal/email"
)

type EmailConsumer func(f *jsemail.ForwardableEmailMessage) error

var consumer EmailConsumer = func(_ *jsemail.ForwardableEmailMessage) error {
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
	jsclass.Env = envObj
	jsclass.ExcutionContext = ctxObj

	env.LoadEnvs()
	email := jsemail.NewForwardableEmailMessage(emailObj)
	defer email.Raw.Close()

	return consumer(email)
}

func ConsumeNonBlock(c EmailConsumer) {
	consumer = c
}
