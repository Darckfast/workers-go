//go:build js && wasm

package cron

import (
	"syscall/js"
	"time"

	jsconv "github.com/Darckfast/workers-go/internal/conv"
)

type CronEvent struct {
	Cron          string
	ScheduledTime time.Time
}

func NewEvent(obj js.Value) *CronEvent {
	scheduledTimeVal := jsconv.MaybeInt64(obj.Get("scheduledTime"))
	return &CronEvent{
		Cron:          obj.Get("cron").String(),
		ScheduledTime: time.Unix(0, scheduledTimeVal*int64(time.Millisecond)).UTC(),
	}
}
