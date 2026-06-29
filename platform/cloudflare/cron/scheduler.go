//go:build js && wasm

package cron

import (
	"context"
	"errors"
	"syscall/js"

	"codeberg.org/darckfast/workers-go/internal/jsclass"
	"codeberg.org/darckfast/workers-go/platform/cloudflare/bind"
)

type Task func(c context.Context, evt *CronEvent) error

var scheduledTask Task = func(c context.Context, _ *CronEvent) error {
	return errors.New("no scheduled implemented")
}

func runScheduler(jsEvent js.Value, envObj js.Value, ctxObj js.Value) error {
	bind.Env.LoadEnvs(envObj)
	bind.Ctx.Init(ctxObj)

	event := NewEvent(jsEvent)

	ctx := context.Background()
	return scheduledTask(ctx, event)
}

func init() {
	promise := js.FuncOf(func(_ js.Value, args []js.Value) any {
		controllerObj := args[0]
		envObj := args[1]
		ctxObj := args[2]

		var cb js.Func
		cb = js.FuncOf(func(_ js.Value, pArgs []js.Value) any {
			defer cb.Release()
			resolve := pArgs[0]
			reject := pArgs[1]

			go func() {
				err := runScheduler(controllerObj, envObj, ctxObj)

				if err != nil {
					reject.Invoke(jsclass.ToJSError(err))
				} else {
					resolve.Invoke(js.Undefined())
				}
			}()
			return nil
		})

		return jsclass.Promise.New(cb)
	})
	jsclass.CF.Set("scheduled", promise)

}

func ScheduleTaskNonBlock(task Task) {
	scheduledTask = task
}
