//go:build js && wasm

package cron

import (
	"errors"
	"syscall/js"

	jsclass "github.com/Darckfast/workers-go/internal/class"
	"github.com/Darckfast/workers-go/platform/cloudflare/env"
	"github.com/Darckfast/workers-go/platform/cloudflare/lifecycle"
)

type Task func(evt *CronEvent) error

var scheduledTask Task = func(_ *CronEvent) error {
	return errors.New("no scheduled implemented")
}

func runScheduler(jsEvent js.Value, envObj js.Value, ctxObj js.Value) error {
	lifecycle.Env = envObj
	lifecycle.Ctx = jsclass.ExecutionContextWrap{Ctx: ctxObj}
	err := env.LoadEnvs()
	if err != nil {
		return err
	}
	event := NewEvent(jsEvent)

	return scheduledTask(event)
}

func init() {
	runSchedulerCallback := js.FuncOf(func(_ js.Value, args []js.Value) any {
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

	js.Global().Get("cf").Set("scheduled", runSchedulerCallback)
}

func ScheduleTaskNonBlock(task Task) {
	scheduledTask = task
}
