//go:build js && wasm

/*
Package email is the glue code for Cloudflare's Worker Email bindings
*/
package email

import (
	"context"
	"errors"
	"syscall/js"

	"codeberg.org/darckfast/workers-go/internal/jsclass"
)

type EmailConsumer func(c context.Context, f *ForwardableEmailMessage) error

var consumer EmailConsumer = func(c context.Context, _ *ForwardableEmailMessage) error {
	return errors.New("no consumer implemented")
}

func init() {
	var promise = js.FuncOf(func(this js.Value, args []js.Value) any {
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

	jsclass.CF.Set("email", promise)
}

func handler(emailObj, envObj, ctxObj js.Value) error {
	jsclass.Env.LoadEnvs(envObj)
	jsclass.Ctx.Init(ctxObj)

	email := NewForwardableEmailMessage(emailObj)
	defer func() {
		err := email.Raw.Close()
		if err != nil {
			println("error closing email raw body reader", err.Error())
		}
	}()

	ctx := context.Background()
	return consumer(ctx, email)
}

func ConsumeNonBlock(c EmailConsumer) {
	consumer = c
}
