//go:build !js && !wasm

package main

import (
	_ "embed"
	"log"
	"strings"
)

const (
	Bold     = "\033[1m"
	Reset    = "\033[0m"
	Red      = "\033[38;2;240;82;82m"
	Green    = "\033[38;2;61;214;140m"
	Amber    = "\033[38;2;240;165;0m"
	Blue     = "\033[38;2;91;164;245m"
	Cyan     = "\033[38;2;34;211;238m"
	Pink     = "\033[38;2;232;121;160m"
	Purple   = "\033[38;2;167;139;250m"
	White    = "\033[38;2;232;232;232m"
	Gray     = "\033[38;2;136;136;136m"
	DarkGray = "\033[38;2;85;85;85m"
)

var rp = strings.NewReplacer(
	"{Bold}", Bold,
	"{Reset}", Reset,
	"{Red}", Red,
	"{Green}", Green,
	"{Amber}", Amber,
	"{Blue}", Blue,
	"{Cyan}", Cyan,
	"{Pink}", Pink,
	"{Purple}", Purple,
	"{White}", White,
	"{Gray}", Gray,
	"{DarkGray}", DarkGray,
)

func init() {
	log.SetFlags(0)
}

func erro(format string, v ...any) {
	log.Printf(Red+rp.Replace(format)+Reset+"\n", v...)
}

func warn(format string, v ...any) {
	log.Printf(Amber+rp.Replace(format)+Reset+"\n", v...)
}

func info(format string, v ...any) {
	log.Printf(Gray+rp.Replace(format)+Reset+"\n", v...)
}
