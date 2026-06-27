//go:build js && wasm

/*
Package cron is the glue code for Cloudflare's Worker cron handler
*/
package cron

import (
	"syscall/js"
	"time"

	"codeberg.org/darckfast/workers-go/internal/jsconv"
)

type CronEvent struct {
	ScheduledTime time.Time
	Cron          string
}

func NewEvent(obj js.Value) *CronEvent {
	scheduledTimeVal := jsconv.MaybeInt64(obj.Get("scheduledTime"))
	return &CronEvent{
		Cron:          obj.Get("cron").String(),
		ScheduledTime: time.Unix(0, scheduledTimeVal*int64(time.Millisecond)).UTC(),
	}
}
