//go:build js && wasm

package main

import (
	cronhandler "github.com/Darckfast/workers-go/worker/pkg/cron"
	emailhandler "github.com/Darckfast/workers-go/worker/pkg/email"
	fetchhandler "github.com/Darckfast/workers-go/worker/pkg/fetch"
	queuehandler "github.com/Darckfast/workers-go/worker/pkg/queue"
	tailhandler "github.com/Darckfast/workers-go/worker/pkg/tail"
)

func main() {
	fetchhandler.New()
	emailhandler.New()
	cronhandler.New()
	queuehandler.New()
	tailhandler.New()

	/**
	 * This code below is REQUIRED, it's what will keep your Go's WASM binary
	 * running while Cloudflare Worker consumes it
	 *
	 * Without this code, the process will exit before the handlers call it
	 */
	<-make(chan struct{})
}
