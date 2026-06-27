//go:build js && wasm

package jsclass

import (
	"os"
	"syscall/js"

	"codeberg.org/darckfast/workers-go/internal/jshelper"
	"codeberg.org/darckfast/workers-go/internal/jstry"
	easyjson "github.com/mailru/easyjson"
)

type JSONJS struct {
	v jshelper.LazyJSVal
}

func (j *JSONJS) Stringify(args ...any) js.Value {
	j.v.Init("JSON")
	return j.v.Call("stringify", args...)
}

func (j *JSONJS) Parse(args ...any) (js.Value, error) {
	j.v.Init("JSON")
	return jstry.And.Catch(j.v.Value, "parse", args...)
}

type ObjectWrap struct {
	c jshelper.LazyJSVal
}

func (o *ObjectWrap) Class() js.Value {
	o.c.Init("Object")
	return o.c.Value
}

func (o *ObjectWrap) New(args ...any) js.Value {
	o.c.Init("Object")
	return o.c.New(args...)
}

func (o *ObjectWrap) FromEntries(args ...any) js.Value {
	o.c.Init("Object")
	return o.c.Call("fromEntries", args...)
}

func (o *ObjectWrap) GetPrototypeOf(arg js.Value) js.Value {
	o.c.Init("Object")
	return o.c.Call("getPrototypeOf", arg)
}

type ExecutionContextWrap struct {
	v jshelper.LazyJSVal
}

type EnvBinding struct {
	v jshelper.LazyJSVal
}

func JSValueToMapString(v js.Value) (GenericStringMap, error) {
	obj := GenericStringMap{}
	if !v.Truthy() {
		return obj, nil
	}

	jsonStr := JSON.Stringify(v).String()
	err := easyjson.Unmarshal([]byte(jsonStr), &obj)

	return obj, err
}

func (e *EnvBinding) LoadEnvs(j js.Value) {
	e.v.Do(func() {
		if !j.Truthy() {
			var workerapp = js.Global().Get("workerapp")

			if workerapp.Truthy() {
				e.v.Value = workerapp.Get("env")
			}
		} else {
			e.v.Value = j
		}

		if e.v.Truthy() {
			envs, err := JSValueToMapString(e.v.Value)
			if err != nil {
				println("error setting envs: " + err.Error())
			}

			for key, e := range envs {
				err = os.Setenv(key, e)
				if err != nil {
					println("error setting envs: " + err.Error())
				}
			}
		}
	})
}

func (e *EnvBinding) Get(n string) js.Value {
	e.LoadEnvs(js.Value{})
	return e.v.Get(n)
}

func (e *ExecutionContextWrap) Init(j js.Value) {
	e.v.Do(func() {
		if !j.Truthy() {

			var workerapp = js.Global().Get("workerapp")

			if workerapp.Truthy() {
				e.v.Value = workerapp.Get("ctx")
			}
		} else {
			e.v.Value = j
		}
	})
}

func (e *ExecutionContextWrap) WaitUntil(task func() error) {
	e.Init(js.Value{})
	var cb js.Func

	cb = js.FuncOf(func(this js.Value, pArgs []js.Value) any {
		resolve := pArgs[0]
		reject := pArgs[1]

		go func(res, rej js.Value) {
			defer cb.Release()

			err := task()
			if err != nil {
				reject.Invoke(ToJSError(err))
			} else {
				resolve.Invoke(true)
			}
		}(resolve, reject)

		return nil
	})

	e.v.Call("waitUntil", Promise.New(cb))
}

func (e *ExecutionContextWrap) PassThroughOnException() {
	e.Init(js.Value{})
	e.v.Call("passThroughOnException")
}

type ConsoleWrap struct {
	v jshelper.LazyJSVal
}

func (c *ConsoleWrap) Log(v ...any) {
	c.v.Init("console")
	c.v.Call("log", v...)
}

type ConnectWrap struct {
	c jshelper.LazyJSVal
}

func (c *ConnectWrap) init() {
	c.c.Do(func() {
		c.c.Value = js.Global().Get("cf").Get("connect")
	})
}

func (c *ConnectWrap) Get() js.Value {
	c.init()
	return c.c.Value
}
