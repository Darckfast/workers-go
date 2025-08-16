package jsutil

import (
	"bytes"
	"fmt"
	"io"
	"syscall/js"
)

type RawJSBodyWriter interface {
	WriteRawJSBody(body js.Value)
}

type RawJSBodyGetter interface {
	GetRawJSBody() js.Value
}

// ReadableStream implements io.Reader sourced from ReadableStreamDefaultReader.
//   - ReadableStreamDefaultReader: https://developer.mozilla.org/en-US/docs/Web/API/ReadableStreamDefaultReader
//   - This implementation is based on: https://deno.land/std@0.139.0/streams/conversion.ts#L76
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

// Read reads bytes from ReadableStreamDefaultReader.
func (rs *ReadableStream) Read(p []byte) (n int, err error) {
	if rs.reader == nil {
		r := rs.stream.Call("getReader")
		rs.reader = &r
	}
	if rs.buf.Len() == 0 {
		resultCh := make(chan js.Value)
		errCh := make(chan error)

		promise := rs.reader.Call("read")

		var then, catch js.Func

		then = js.FuncOf(func(_ js.Value, args []js.Value) any {
			defer then.Release()
			result := args[0]
			if result.Get("done").Bool() {
				errCh <- io.EOF
				return nil
			}
			resultCh <- result.Get("value")
			return nil
		})

		catch = js.FuncOf(func(_ js.Value, args []js.Value) any {
			defer catch.Release()
			result := args[0]
			errCh <- fmt.Errorf("JavaScript error on read: %s", result.Call("toString").String())
			return nil
		})

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
	sr.reader.Call("cancel")
	return nil
}

// readerWrapper is wrapper to disable readableStreamToReadCloser's WriteTo method.
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

// readerToReadableStream implements ReadableStream sourced from io.ReadCloser.
//   - ReadableStream: https://developer.mozilla.org/docs/Web/API/ReadableStream
//   - This implementation is based on: https://deno.land/std@0.139.0/streams/conversion.ts#L230
type readerToReadableStream struct {
	size        int64
	initialized bool
	reader      io.ReadCloser
	chunkBuf    []byte
}

// Pull implements ReadableStream's pull method.
//   - https://developer.mozilla.org/en-US/docs/Web/API/ReadableStream/ReadableStream#pull
func (rs *readerToReadableStream) Pull(controller js.Value) error {
	if !rs.initialized {
		ua := NewUint8Array(0)
		controller.Call("enqueue", ua)
		rs.initialized = true
		return nil
	}
	n, err := rs.reader.Read(rs.chunkBuf)
	if n != 0 {
		rs.size += int64(n)
		ua := NewUint8Array(n)
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
		controller.Call("error", Error(err.Error()))
		if err := rs.reader.Close(); err != nil {
			return err
		}
		return err
	}
	return nil
}

// Cancel implements ReadableStream's cancel method.
//   - https://developer.mozilla.org/en-US/docs/Web/API/ReadableStream/ReadableStream#cancel
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
	rsInit := NewObject()
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
					reject.Invoke(Error(err.Error()))
					return
				}
				resolve.Invoke()
			}()
			return nil
		})
		return NewPromise(cb)
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
					reject.Invoke(Error(err.Error()))
					return
				}
				resolve.Invoke()
			}()
			return nil
		})
		return NewPromise(cb)
	}))
	return ReadableStreamClass.New(rsInit)
}

// ConvertReaderToFixedLengthStream converts io.ReadCloser to TransformStream.
func ConvertReaderToFixedLengthStream(rc io.ReadCloser, size int64) js.Value {
	stream := MaybeFixedLengthStreamClass.New(js.ValueOf(size))
	go func(writer js.Value) {
		defer rc.Close()

		chunk := make([]byte, min(size, defaultChunkSize))
		AwaitPromise(writer.Get("ready"))
		for {
			n, err := rc.Read(chunk)
			if n > 0 {
				b := Uint8ArrayClass.New(n)
				js.CopyBytesToJS(b, chunk[:n])
				writer.Call("write", b)
			}
			if err != nil {
				AwaitPromise(writer.Get("ready"))
				writer.Call("close")
				return
			}
		}
	}(stream.Get("writable").Call("getWriter"))
	return stream.Get("readable")
}
