package cron

import (
	"context"
	"fmt"
	"syscall/js"

	runtimecontext "github.com/syumai/workers/internal/runtime"
	jsutil "github.com/syumai/workers/internal/utils"
)

type Task func(ctx context.Context) error

var scheduledTask Task = func(ctx context.Context) error {
	return fmt.Errorf("no scheduled implemented")
}

func runScheduler(eventObj js.Value, envObj js.Value, ctxObj js.Value) error {
	jsutil.RuntimeEnv = envObj
	jsutil.RuntimeExcutionContext = ctxObj

	ctx := runtimecontext.New(context.Background(), eventObj)

	if err := scheduledTask(ctx); err != nil {
		return err
	}
	return nil
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
					reject.Invoke(jsutil.Error(err.Error()))
				} else {
					resolve.Invoke(js.Undefined())
				}
			}()
			return nil
		})

		return jsutil.NewPromise(cb)
	})

	js.Global().Get("cf").Set("scheduled", runSchedulerCallback)
}

func ScheduleTaskNonBlock(task Task) {
	scheduledTask = task
}
