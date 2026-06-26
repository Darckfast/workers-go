//go:build js && wasm

package main

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"codeberg.org/darckfast/workers-go/platform/cloudflare/durableobjects"
)

// Durable Object to be exported
type MyDurableObject struct {
	durableobjects.DurableObject // Required embed
	InitAt                       int64
}

// type Greeting = { Time: number, Custom: { ID: string, Name: string, Time: number, On: boolean, Ratio: number } }
type Greeting struct {
	Time   int64
	Custom struct {
		ID    string
		Name  string
		Time  int64
		On    bool
		Ratio float64
	}
}

// async SayHello(g: Greeting): Promise<string>
func (d *MyDurableObject) WhatTimsIsIt(_ context.Context, g *Greeting) Greeting {
	return Greeting{Time: time.Now().UnixMilli()}
}

// async Goodbye(g: Greeting, a: number): Promise<[string, number]>
func (d *MyDurableObject) Goodbye(_ context.Context, g *Greeting, a *int64) (string, int64) {
	return "Goodbye", *a
}

// async SayHello(g: Greeting): Promise<string>
func (d *MyDurableObject) SayHello(_ context.Context, g *Greeting) string {
	return "Hello world, " + strconv.FormatInt(d.InitAt, 10)
}

// async fetch(r: Request) Promise<Response>
func (d *MyDurableObject) fetch(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello from fetch request"))
}

// called once inside the constructor()
func (d *MyDurableObject) init(_ context.Context) {
	d.InitAt = time.Now().UnixMicro()
}
