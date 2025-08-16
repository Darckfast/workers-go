package fetch

import (
	"errors"
	"net/http"
	"syscall/js"

	jshttp "github.com/syumai/workers/internal/http"
	jsutil "github.com/syumai/workers/internal/utils"
)

// fetch is a function that reproduces cloudflare fetch.
// Docs: https://developers.cloudflare.com/workers/runtime-apis/fetch/
func fetch(namespace js.Value, req *http.Request, init *RequestInit) (*http.Response, error) {
	if namespace.IsUndefined() {
		return nil, errors.New("js namespace not set, fetch function not found")
	}
	fetchFunc := namespace.Get("fetch")
	promise := fetchFunc.Invoke(
		// The Request object to fetch.
		// Docs: https://developers.cloudflare.com/workers/runtime-apis/request
		jshttp.ToJSRequest(req),
		// The content of the request.
		// Docs: https://developers.cloudflare.com/workers/runtime-apis/request#requestinit
		init.ToJS(),
	)

	jsRes, err := jsutil.AwaitPromise(promise)
	if err != nil {
		return nil, err
	}

	return jshttp.ToResponse(jsRes)
}
