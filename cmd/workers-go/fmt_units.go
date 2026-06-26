//go:build !js && !wasm

package main

import (
	_ "embed"
	"strconv"
	"time"
)

type unit struct {
	name string
	dur  time.Duration
}

var byteUnits = []string{"B", "KiB", "MiB"}

const byteSize = 1 << 10

var durationUnits = []unit{
	{"m", time.Minute},
	{"s", time.Second},
	{"ms", time.Millisecond},
	{"µs", time.Microsecond},
	{"ns", time.Nanosecond},
}

func fmtDuration(d time.Duration) string {
	for _, u := range durationUnits {
		if d >= u.dur {
			return strconv.FormatFloat(float64(d)/float64(u.dur), 'f', 2, 64) + " " + u.name
		}
	}

	return "0 ns"
}

func fmtBytes(b int64) string {
	v := float64(b)
	i := 0

	for v >= byteSize && i < len(byteUnits)-1 {
		v /= byteSize
		i++
	}

	return strconv.FormatFloat(v, 'f', 2, 64) + " " + byteUnits[i]
}
