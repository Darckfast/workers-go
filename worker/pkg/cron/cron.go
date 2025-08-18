package cronhandler

import (
	"github.com/Darckfast/workers-go/cloudflare/cron"
	cloudflare "github.com/Darckfast/workers-go/cloudflare/ctx"
)

func New() {

	/*
	 * This functions must be called, it what will instantiate a cron task consumer
	 */
	cron.ScheduleTaskNonBlock(func(ctx *cron.CronEvent) error {
		cloudflare.WaitUntil(func() {

		})

		return nil
	})
}
