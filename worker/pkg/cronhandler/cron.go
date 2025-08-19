//go:build js && wasm

package cronhandler

import (
	"log"

	"github.com/Darckfast/workers-go/cloudflare/cron"
	cloudflare "github.com/Darckfast/workers-go/cloudflare/ctx"
)

func New() {

	/*
	 * ScheduleTaskNonBlock functions must be called, it what will instantiate a cron task consumer
	 */
	cron.ScheduleTaskNonBlock(func(event *cron.CronEvent) error {
		cloudflare.WaitUntil(func() {
			log.Println("running my scheduled task at " + event.ScheduledTime.String())
		})

		return nil
	})
}
