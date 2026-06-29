//go:build js && wasm

package fetch

import (
	"net/http"
	"syscall/js"
	"time"

	"codeberg.org/darckfast/workers-go/internal/jsclass"
	"codeberg.org/darckfast/workers-go/internal/jshttp"
	"codeberg.org/darckfast/workers-go/internal/jsruntime"
	"codeberg.org/darckfast/workers-go/platform/cloudflare/bind"
	"github.com/mailru/easyjson"
)

type Client struct {
	namespace    js.Value
	CF           *RequestInitCF
	RedirectMode string
	Timeout      time.Duration
}

type Transport struct{}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	c := Client{}

	return c.Do(req)
}

var _ http.RoundTripper = (*Transport)(nil)

// Deprecated: This can be used normally, just be aware by transforming it into
// http.Client, the compiler will also include the crypto lib, and
// it can increase the final binary size from 5.6MB to 11MB
// the compressed file can increase from 1.6MB to 2.8MB
func (c *Client) ToHTTPClient() *http.Client {
	return &http.Client{
		Timeout:   c.Timeout,
		Transport: &Transport{},
	}
}

func (c *Client) WithBinding(bindname string) *Client {
	c.namespace = bind.Env.Get(bindname)
	return c
}

func (c *Client) WithCF(cf *RequestInitCF) *Client {
	c.CF = cf
	return c
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	if c.namespace.IsUndefined() {
		c.namespace = js.Global()
	}

	// This client is incompatible with the current container.fetch
	fetchFunc := c.namespace.Get("fetch")

	if c.RedirectMode == "" {
		c.RedirectMode = "follow"
	}

	initOptions := InitOptions{
		Redirect:    c.RedirectMode,
		Credentials: "omit",
	}

	initJSON, _ := easyjson.Marshal(initOptions)
	initObj, _ := jsclass.JSON.Parse(string(initJSON))

	if c.Timeout != 0 {
		timeoutSignal := jsclass.AbortSignal.Timeout(c.Timeout.Milliseconds())
		reqSignal := req.Context().Value(jsruntime.CtxSignal{}).(js.Value)
		initObj.Set("signal", jsclass.AbortSignal.Any([]any{timeoutSignal, reqSignal}))
	} else {
		initObj.Set("signal", req.Context().Value("signal"))
	}

	if c.CF != nil {
		cfJSON, _ := easyjson.Marshal(c.CF)
		cfObj, _ := jsclass.JSON.Parse(string(cfJSON))
		initObj.Set("cf", cfObj)
	}

	promise := fetchFunc.Invoke(
		jshttp.ToJSRequest(req),
		initObj,
	)

	jsRes, err := jsclass.Await(promise)
	if err != nil {
		return nil, err
	}

	return jshttp.ToResponse(jsRes), nil
}

func NewClient() *Client {
	return &Client{}
}
