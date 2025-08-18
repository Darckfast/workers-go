//go:build js && wasm

package cron

import (
	"syscall/js"
	"time"
)

type CronEvent struct {
	Cron          string
	ScheduledTime time.Time
}

func NewEvent(obj js.Value) *CronEvent {
	scheduledTimeVal := obj.Get("scheduledTime").Float()
	return &CronEvent{
		Cron:          obj.Get("cron").String(),
		ScheduledTime: time.Unix(int64(scheduledTimeVal)/1000, 0).UTC(),
	}
}
