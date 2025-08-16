package fetch

import (
	"net/http"
	"syscall/js"
)

// CFTransport is an implementation of http.RoundTripper
type CFTransport struct {
	// namespace - Objects that Fetch API belongs to. Default is Global
	namespace js.Value
	redirect  RedirectMode
}

// RoundTrip replaces http.DefaultTransport.RoundTrip to use cloudflare fetch
func (t *CFTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return fetch(t.namespace, req, &RequestInit{
		Redirect: t.redirect,
	})
}
