//go:build js && wasm

package cronhandler

import (
	"fmt"

	"github.com/Darckfast/workers-go/cloudflare/cron"
	"github.com/Darckfast/workers-go/cloudflare/kv"
)

func New() {

	/*
	 * ScheduleTaskNonBlock functions must be called, it what will instantiate a cron task consumer
	 */
	cron.ScheduleTaskNonBlock(func(event *cron.CronEvent) error {
		kvStore, _ := kv.NewNamespace("TEST_NAMESPACE")
		kvStore.PutString("cron:result", fmt.Sprintf("%d", event.ScheduledTime.UnixMilli()), nil)

		return nil
	})
}
