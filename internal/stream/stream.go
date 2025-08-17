package jsstream

import (
	"bytes"
	"io"
	"syscall/js"

	jsclass "github.com/syumai/workers/internal/class"
)

type RawJSBodyWriter interface {
	WriteRawJSBody(body js.Value)
}

type RawJSBodyGetter interface {
	GetRawJSBody() js.Value
}

type ReadableStream struct {
	buf    bytes.Buffer
	stream js.Value
	reader *js.Value
}

var (
	_ io.ReadCloser   = (*ReadableStream)(nil)
	_ io.WriterTo     = (*ReadableStream)(nil)
	_ RawJSBodyGetter = (*ReadableStream)(nil)
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

func (sr *ReadableStream) WriteTo(w io.Writer) (n int64, err error) {
	if w, ok := w.(RawJSBodyWriter); ok {
		w.WriteRawJSBody(sr.stream)
		return 0, nil
	}
	return io.Copy(w, &readerWrapper{sr})
}

func (sr *ReadableStream) GetRawJSBody() js.Value {
	return sr.stream
}

func ReadableStreamToReadCloser(stream js.Value) io.ReadCloser {
	return &ReadableStream{
		stream: stream,
	}
}

type readerToReadableStream struct {
	initialized bool
	reader      io.ReadCloser
	chunkBuf    []byte
}

func (rs *readerToReadableStream) Pull(controller js.Value) error {
	if !rs.initialized {
		ua := jsclass.Uint8Array.New(0)
		controller.Call("enqueue", ua)
		rs.initialized = true
		return nil
	}
	n, err := rs.reader.Read(rs.chunkBuf)
	if n != 0 {
		ua := jsclass.Uint8Array.New(n)
		js.CopyBytesToJS(ua, rs.chunkBuf[:n])
		controller.Call("enqueue", ua)
	}

	// Cloudflare Workers sometimes call `pull` to closed ReadableStream.
	// When the call happens, `io.ErrClosedPipe` should be ignored.
	if err == io.EOF || err == io.ErrClosedPipe {
		controller.Call("close")
		if err := rs.reader.Close(); err != nil {
			return err
		}
		return nil
	}
	if err != nil {
		controller.Call("error", jsclass.Error.New(err.Error()))
		if err := rs.reader.Close(); err != nil {
			return err
		}
		return err
	}
	return nil
}

func (rs *readerToReadableStream) Cancel() error {
	return rs.reader.Close()
}

// https://deno.land/std@0.139.0/streams/conversion.ts#L5
const defaultChunkSize = 16_640

func ReadCloserToReadableStream(reader io.ReadCloser) js.Value {
	stream := &readerToReadableStream{
		reader:   reader,
		chunkBuf: make([]byte, defaultChunkSize),
	}
	rsInit := jsclass.Object.New()
	rsInit.Set("pull", js.FuncOf(func(_ js.Value, args []js.Value) any {
		var cb js.Func
		cb = js.FuncOf(func(this js.Value, pArgs []js.Value) any {
			defer cb.Release()
			resolve := pArgs[0]
			reject := pArgs[1]
			controller := args[0]
			go func() {
				err := stream.Pull(controller)
				if err != nil {
					reject.Invoke(jsclass.Error.New(err.Error()))
					return
				}
				resolve.Invoke()
			}()
			return nil
		})
		return jsclass.Promise.New(cb)
	}))
	rsInit.Set("cancel", js.FuncOf(func(js.Value, []js.Value) any {
		var cb js.Func
		cb = js.FuncOf(func(this js.Value, pArgs []js.Value) any {
			defer cb.Release()
			resolve := pArgs[0]
			reject := pArgs[1]
			go func() {
				err := stream.Cancel()
				if err != nil {
					reject.Invoke(jsclass.Error.New(err.Error()))
					return
				}
				resolve.Invoke()
			}()
			return nil
		})
		return jsclass.Promise.New(cb)
	}))
	return jsclass.ReadableStream.New(rsInit)
}

func ReadCloserToFixedLengthStream(rc io.ReadCloser, size int64) js.Value {
	stream := jsclass.MaybeFixedLengthStream.New(js.ValueOf(size))
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
