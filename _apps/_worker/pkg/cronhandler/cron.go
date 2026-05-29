//go:build js && wasm

package cronhandler

import (
	"fmt"

	"codeberg.org/darckfast/workers-go/platform/cloudflare/cron"
	"codeberg.org/darckfast/workers-go/platform/cloudflare/kv"
)

func New() {

	/*
	 * ScheduleTaskNonBlock functions must be called, it what will instantiate a cron task consumer
	 */
	cron.ScheduleTaskNonBlock(func(event *cron.CronEvent) error {
		kvStore, _ := kv.NewNamespace("TEST_NAMESPACE")
		return kvStore.Put("cron:result", fmt.Sprintf("%d", event.ScheduledTime.UnixMilli()), nil)
	})
}
