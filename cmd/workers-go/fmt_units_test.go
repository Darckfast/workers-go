//go:build !js && !wasm

package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFmtDuration(t *testing.T) {
	cases := []struct {
		input    time.Duration
		expected string
	}{
		{2 * time.Minute, "2.00 m"},
		{90 * time.Second, "1.50 m"},
		{time.Second, "1.00 s"},
		{1500 * time.Millisecond, "1.50 s"},
		{time.Millisecond, "1.00 ms"},
		{time.Microsecond, "1.00 µs"},
		{time.Nanosecond, "1.00 ns"},
		{0, "0 ns"},
	}
	for _, tc := range cases {
		t.Run(tc.expected, func(t *testing.T) {
			assert.Equal(t, tc.expected, fmtDuration(tc.input))
		})
	}
}

func TestFmtBytes(t *testing.T) {
	cases := []struct {
		input    int64
		expected string
	}{
		{0, "0.00 B"},
		{512, "512.00 B"},
		{1023, "1023.00 B"},
		{1024, "1.00 KiB"},
		{1536, "1.50 KiB"},
		{1024 * 1024, "1.00 MiB"},
		{1024 * 1024 * 1024, "1024.00 MiB"}, // caps at MiB
	}
	for _, tc := range cases {
		t.Run(tc.expected, func(t *testing.T) {
			assert.Equal(t, tc.expected, fmtBytes(tc.input))
		})
	}
}
