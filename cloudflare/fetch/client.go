package fetch

import (
	"context"
	"io"
	"net/http"
	"syscall/js"

	jsutil "github.com/syumai/workers/internal/utils"
)

// Client is an HTTP client.
type Client struct {
	// namespace - Objects that Fetch API belongs to. Default is Global
	namespace js.Value
}

// applyOptions applies client options.
func (c *Client) applyOptions(opts []ClientOption) {
	for _, opt := range opts {
		opt(c)
	}
}

func (c *Client) WithBinding(bindname string) *Client {
	c.namespace = jsutil.RuntimeEnv.Get(bindname)
	return c
}

// HTTPClient returns *http.Client.
func (c *Client) HTTPClient(redirect RedirectMode) *http.Client {
	return &http.Client{
		Transport: &CFTransport{
			namespace: c.namespace,
			redirect:  redirect,
		},
	}
}

// Do sends an HTTP request and returns an HTTP response
func (c *Client) Do(req *http.Request, init *RequestInit) (*http.Response, error) {
	return fetch(c.namespace, req, init)
}

// ClientOption is a type that represents an optional function.
type ClientOption func(*Client)

// WithBinding changes the objects that Fetch API belongs to.
// This is useful for service bindings, mTLS, etc.
func WithBinding(bind js.Value) ClientOption {
	return func(c *Client) {
		c.namespace = bind
	}
}

// NewClient returns new Client
func NewClient(opts ...ClientOption) *Client {
	c := &Client{
		namespace: js.Global(),
	}
	c.applyOptions(opts)

	return c
}

// NewRequest returns new Request given a method, URL, and optional body
func NewRequest(ctx context.Context, method string, url string, body io.Reader) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, method, url, body)
}
