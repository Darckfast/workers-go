//go:build js && wasm

package queues

// contentType represents the content type of a message produced to a queue.
// This information mostly affects how the message body is represented in the Cloudflare UI and is NOT
// propagated to the consumer side.
//   - https://developers.cloudflare.com/queues/configuration/javascript-apis/#queuescontenttype
type contentType string

const (
	// contentTypeJSON is the default content type for the produced queue message.
	// The message body is NOT being marshaled before sending and is passed to js.ValueOf directly.
	// Make sure the body is serializable to JSON.
	//  - https://pkg.go.dev/syscall/js#ValueOf
	contentTypeJSON contentType = "json"

	// contentTypeV8 is currently treated the same as QueueContentTypeJSON.
	contentTypeV8 contentType = "v8"

	// contentTypeText is used to send a message as a string.
	// Supported body types are string, []byte and io.Reader.
	contentTypeText contentType = "text"

	// contentTypeBytes is used to send a message as a byte array.
	// Supported body types are string, []byte, and io.Reader.
	contentTypeBytes contentType = "bytes"
)
