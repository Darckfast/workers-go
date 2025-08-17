package cronhandler

import (
	"fmt"

	"github.com/syumai/workers/cloudflare/cron"
	cloudflare "github.com/syumai/workers/cloudflare/ctx"
)

func New() {
	cron.ScheduleTaskNonBlock(func(ctx *cron.CronEvent) error {
		fmt.Println("cronjob executed")

		cloudflare.WaitUntil(func() {
			fmt.Println("cronjob executed")
		})

		return nil
	})
}
