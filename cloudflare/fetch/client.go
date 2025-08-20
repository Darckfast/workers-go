//go:build js && wasm

package fetch

import (
	"encoding/json"
	"net/http"
	"syscall/js"
	"time"

	"github.com/Darckfast/workers-go/cloudflare/lifecycle"
	jsclass "github.com/Darckfast/workers-go/internal/class"
	jshttp "github.com/Darckfast/workers-go/internal/http"
)

type RequestInitCF struct {
	// Using json tags will use reflect
	// by default this is slower than manually setting the values
	// in a string template or similar, but it's easier to maintain
	Apps             bool           `json:"apps"`
	CacheEverything  bool           `json:"cacheEverything"`
	CacheKey         string         `json:"cacheKey"`
	CacheTags        []string       `json:"cacheTags"`
	CacheTtl         int            `json:"cacheTtl"`
	CacheTtlByStatus map[string]int `json:"cacheTtlByStatus"`
	Mirage           bool           `json:"mirage"`
	Polish           string         `json:"polish"`
	ResolveOverride  string         `json:"resolveOverride"`
	ScrapShield      bool           `json:"scrapShield"`
	Webp             bool           `json:"webp"`
}

type Client struct {
	RedirectMode string
	Timeout      time.Duration
	namespace    js.Value
	CF           *RequestInitCF
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
	c.namespace = lifecycle.Env.Get(bindname)
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

	initOptions := map[string]string{
		"redirect":    c.RedirectMode,
		"credentials": "omit",
	}

	initJson, _ := json.Marshal(initOptions)
	initObj, _ := jsclass.JSON.Parse(string(initJson))

	if c.Timeout != 0 {
		initObj.Set("signal", jsclass.AbortSignal.Call("timeout", c.Timeout.Milliseconds()))
	}

	if c.CF != nil {
		cfJson, _ := json.Marshal(c.CF)
		cfObj, _ := jsclass.JSON.Parse(string(cfJson))
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
