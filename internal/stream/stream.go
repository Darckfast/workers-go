package jsstream

import (
	"bytes"
	"io"
	"syscall/js"

	jsclass "github.com/syumai/workers/internal/class"
)

type ReadableStream struct {
	buf    bytes.Buffer
	stream js.Value
	reader *js.Value
}

var (
	_ io.ReadCloser = (*ReadableStream)(nil)
)

func (rs *ReadableStream) Read(p []byte) (n int, err error) {
	if rs.reader == nil {
		r := rs.stream.Call("getReader")
		rs.reader = &r
	}
	if rs.buf.Len() == 0 {
		resultCh := make(chan js.Value)
		defer close(resultCh)

		errCh := make(chan error)
		defer close(errCh)

		var then, catch js.Func

		then = js.FuncOf(func(_ js.Value, args []js.Value) any {
			result := args[0]
			if result.Get("done").Bool() {
				errCh <- io.EOF
				return nil
			}
			resultCh <- result.Get("value")
			return nil
		})
		defer then.Release()

		catch = js.FuncOf(func(_ js.Value, args []js.Value) any {
			errCh <- js.Error{Value: args[0]}
			return nil
		})
		defer catch.Release()

		promise := rs.reader.Call("read")
		promise.Call("then", then).Call("catch", catch)

		select {
		case result := <-resultCh:
			chunk := make([]byte, result.Get("byteLength").Int())
			_ = js.CopyBytesToGo(chunk, result)
			// The length written is always the same as the length of chunk, so it can be discarded.
			//   - https://pkg.go.dev/bytes#Buffer.Write
			_, err := rs.buf.Write(chunk)
			if err != nil {
				return 0, err
			}
		case err := <-errCh:
			return 0, err
		}
	}
	return rs.buf.Read(p)
}

func (sr *ReadableStream) Close() error {
	if sr.reader == nil {
		return nil
	}
	//this returns a promise, maybe it should be waited
	sr.reader.Call("cancel")
	return nil
}

type readerWrapper struct {
	io.Reader
}

func ReadableStreamToReadCloser(stream js.Value) io.ReadCloser {
	return &ReadableStream{
		stream: stream,
	}
}

// https://deno.land/std@0.139.0/streams/conversion.ts#L5
const defaultChunkSize = 16_640

func ReadCloserToReadableStream(reader io.ReadCloser) js.Value {
	pull := js.FuncOf(func(this js.Value, args []js.Value) any {
		// Because the ReadableStream implemented above, this func
		// need to be a promise, otherwise it will cause a deadlock
		// if it tries to convert a ReadCloser that came from a Response
		return jsclass.Promise.New(js.FuncOf(func(_ js.Value, pargs []js.Value) any {
			go func() {
				resolve := pargs[0]
				controller := args[0]
				chunk := make([]byte, defaultChunkSize)

				n, err := reader.Read(chunk)
				if err != nil && err != io.EOF && err != io.ErrClosedPipe {
					controller.Call("error", err.Error())
					return
				}

				if n > 0 {
					uint8Array := js.Global().Get("Uint8Array").New(n)
					js.CopyBytesToJS(uint8Array, chunk[:n])
					controller.Call("enqueue", uint8Array)
				}

				if err == io.EOF || err == io.ErrClosedPipe {
					controller.Call("close")
				}

				resolve.Invoke(true)
			}()
			return nil
		}))
	})

	cancel := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		_ = reader.Close()
		return nil
	})

	rsInit := jsclass.Object.New()
	rsInit.Set("pull", pull)
	rsInit.Set("cancel", cancel)

	return jsclass.ReadableStream.New(rsInit)
}

func ReadCloserToFixedLengthStream(rc io.ReadCloser, size int64) js.Value {
	stream := jsclass.MaybeFixedLengthStream.New(size)
	go func(writer js.Value) {
		defer rc.Close()

		chunk := make([]byte, min(size, 16_640))
		jsclass.Await(writer.Get("ready"))
		for {
			n, err := rc.Read(chunk)

			if n > 0 {
				b := jsclass.Uint8Array.New(n)
				js.CopyBytesToJS(b, chunk[:n])
				writer.Call("write", b)
			}

			if err != nil {
				jsclass.Await(writer.Get("ready"))
				writer.Call("close")
				return
			}
		}
	}(stream.Get("writable").Call("getWriter"))

	return stream.Get("readable")
}
