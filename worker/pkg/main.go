package main

import (
	cronhandler "github.com/syumai/workers/worker/pkg/cron"
	emailhandler "github.com/syumai/workers/worker/pkg/email"
	fetchhandler "github.com/syumai/workers/worker/pkg/fetch"
	queuehandler "github.com/syumai/workers/worker/pkg/queue"
	tailhandler "github.com/syumai/workers/worker/pkg/tail"
)

func main() {
	fetchhandler.New()
	emailhandler.New()
	cronhandler.New()
	queuehandler.New()
	tailhandler.New()

	<-make(chan struct{})
}
