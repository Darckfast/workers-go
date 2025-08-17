package jsutil

import (
	"fmt"
	"strconv"
	"syscall/js"
	"time"

	jsclass "github.com/syumai/workers/internal/class"
)

func init() {
	if js.Global().Get("cf").IsUndefined() {
		cfObj, _ := jsclass.JSON.Parse(`{"ctx":{},"env":{},"handlers":{},"connect":{}}`)
		js.Global().Set("cf", cfObj)

		return
	}

	RuntimeEnv = js.Global().Get("cf").Get("env")
	RuntimeExcutionContext = js.Global().Get("cf").Get("ctx")
	RuntimeConnect = js.Global().Get("cf").Get("connect")
}

var (
	RuntimeEnv             js.Value
	RuntimeExcutionContext js.Value
	RuntimeConnect         js.Value
	RuntimeCache           = js.Global().Get("caches")
	ObjectClass            = js.Global().Get("Object")
	PromiseClass           = js.Global().Get("Promise")
	RequestClass           = js.Global().Get("Request")
	ResponseClass          = js.Global().Get("Response")
	HeadersClass           = js.Global().Get("Headers")
	ArrayClass             = js.Global().Get("Array")
	Uint8ArrayClass        = js.Global().Get("Uint8Array")
	Uint8ClampedArrayClass = js.Global().Get("Uint8ClampedArray")
	ErrorClass             = js.Global().Get("Error")
	ReadableStreamClass    = js.Global().Get("ReadableStream")
	DateClass              = js.Global().Get("Date")
	Null                   = js.ValueOf(nil)
	// MaybeFixedLengthStreamClass is a class for FixedLengthStream.
	// * This class is only available in Cloudflare Workers.
	// * If this class is not available, the value will be undefined.
	MaybeFixedLengthStreamClass = js.Global().Get("FixedLengthStream")
)

func NewObject() js.Value {
	return ObjectClass.New()
}

func NewArray(size int) js.Value {
	return ArrayClass.New(size)
}

func NewUint8Array(size int) js.Value {
	return Uint8ArrayClass.New(size)
}

func NewPromise(fn js.Func) js.Value {
	return PromiseClass.New(fn)
}

func Error(msg string) js.Value {
	return ErrorClass.New(msg)
}

func Errorf(format string, args ...any) js.Value {
	return ErrorClass.New(fmt.Sprintf(format, args...))
}

// ArrayFrom calls Array.from to given argument and returns result Array.
func ArrayFrom(v js.Value) js.Value {
	return ArrayClass.Call("from", v)
}

func AwaitPromise(promise js.Value) (js.Value, error) {
	resultCh := make(chan js.Value)
	defer close(resultCh)

	then := js.FuncOf(func(_ js.Value, args []js.Value) any {
		resultCh <- args[0]
		return nil
	})
	defer then.Release()

	errCh := make(chan error)
	defer close(errCh)

	catch := js.FuncOf(func(_ js.Value, args []js.Value) any {
		errCh <- fmt.Errorf("failed on promise: %s", args[0].Call("toString").String())
		return nil
	})
	defer catch.Release()

	promise.Call("then", then).Call("catch", catch)

	select {
	case result := <-resultCh:
		return result, nil
	case err := <-errCh:
		return js.Value{}, err
	}
}

// StrRecordToMap converts JavaScript side's Record<string, string> into map[string]string.
func StrRecordToMap(v js.Value) map[string]string {
	if v.IsUndefined() || v.IsNull() {
		return map[string]string{}
	}
	entries := ObjectClass.Call("entries", v)
	entriesLen := entries.Get("length").Int()
	result := make(map[string]string, entriesLen)
	for i := range entriesLen {
		entry := entries.Index(i)
		key := entry.Index(0).String()
		value := entry.Index(1).String()
		result[key] = value
	}
	return result
}

// MaybeString returns string value of given JavaScript value or returns "" if the value is undefined.
func MaybeString(v js.Value) string {
	if v.IsUndefined() {
		return ""
	}
	return v.String()
}

// MaybeInt returns int value of given JavaScript value or returns nil if the value is undefined.
func MaybeInt(v js.Value) int {
	if v.IsUndefined() {
		return 0
	}
	return v.Int()
}

func MaybeInt64(v js.Value) int64 {
	if v.IsUndefined() {
		return 0
	}

	vi, _ := strconv.ParseInt(v.String(), 10, 64)

	return vi
}

// MaybeDate returns time.Time value of given JavaScript Date value or returns nil if the value is undefined.
func MaybeDate(v js.Value) (time.Time, error) {
	if v.IsUndefined() {
		return time.Time{}, nil
	}
	return DateToTime(v)
}

// DateToTime converts JavaScript side's Data object into time.Time.
func DateToTime(v js.Value) (time.Time, error) {
	milli := v.Call("getTime").Float()
	return time.UnixMilli(int64(milli)), nil
}

// TimeToDate converts Go side's time.Time into Date object.
func TimeToDate(t time.Time) js.Value {
	return DateClass.New(t.UnixMilli())
}
