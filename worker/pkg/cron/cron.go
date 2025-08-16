package cronhandler

import (
	"context"
	"fmt"

	"github.com/syumai/workers/cloudflare/cron"
	cloudflare "github.com/syumai/workers/cloudflare/ctx"
)

func New() {
	cron.ScheduleTaskNonBlock(func(ctx context.Context) error {
		fmt.Println("cronjob executed")

		cloudflare.WaitUntil(func() {
			fmt.Println("cronjob executed")
		})

		return nil
	})
}
