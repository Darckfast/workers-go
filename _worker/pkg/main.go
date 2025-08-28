//go:build js && wasm

package main

import (
	"worker/pkg/cronhandler"
	"worker/pkg/emailhandler"
	"worker/pkg/fetchhandler"
	"worker/pkg/queuehandler"
	"worker/pkg/tailhandler"
)

func main() {
	// Initialize the http Handler for the globalThis.cf.fetch()
	fetchhandler.New()
	// Initialize the consume for the globalThis.cf.email()
	emailhandler.New()
	// Initialize scheduled task for the globalThis.cf.scheduled()
	cronhandler.New()
	// Initialize the consumer for the globalThis.cf.queue()
	queuehandler.New()
	// Initialize the consumer for the globalThis.cf.tail()
	tailhandler.New()

	/**
	 * This code below is REQUIRED, it's what will keep your Go's WASM process
	 * running while Cloudflare Worker consumes it
	 *
	 * Without this code, the process will exit before the handlers can call it
	 */
	<-make(chan struct{})
}
