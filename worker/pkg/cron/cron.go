package cronhandler

import (
	"github.com/syumai/workers/cloudflare/cron"
	cloudflare "github.com/syumai/workers/cloudflare/ctx"
)

func New() {
	cron.ScheduleTaskNonBlock(func(ctx *cron.CronEvent) error {
		cloudflare.WaitUntil(func() {

		})

		return nil
	})
}
