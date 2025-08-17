package main

import (
	"fmt"

	fetchhandler "github.com/syumai/workers/worker/pkg/fetch"
)

func main() {
	fmt.Println("oi")
	fetchhandler.New()
	// emailhandler.New()
	// cronhandler.New()
	// queuehandler.New()
	// tailhandler.New()

	<-make(chan struct{})
}
