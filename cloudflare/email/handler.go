//go:build js && wasm

package email

import (
	"fmt"
	"syscall/js"

	jsemail "github.com/syumai/workers/internal/email"
	jsutil "github.com/syumai/workers/internal/utils"
)

type EmailConsumer func(f *jsemail.ForwardableEmailMessage) error

var consumer EmailConsumer = func(_ *jsemail.ForwardableEmailMessage) error {
	return fmt.Errorf("no consumer implemented")
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
					reject.Invoke(jsutil.Error(err.Error()))
				} else {
					resolve.Invoke(true)
				}
			}()

			return nil
		})

		return jsutil.NewPromise(cb)
	})

	js.Global().Get("cf").Set("email", handleRequestPromise)
}

func handler(emailObj, envObj, ctxObj js.Value) error {
	jsutil.RuntimeEnv = envObj
	jsutil.RuntimeExcutionContext = ctxObj

	email := jsemail.NewForwardableEmailMessage(emailObj)
	defer email.Raw.Close()

	return consumer(email)
}

func ConsumeNonBlock(c EmailConsumer) {
	consumer = c
}
