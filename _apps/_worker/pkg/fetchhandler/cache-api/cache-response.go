//go:build js && wasm

package httpcache

import (
	"bytes"
	"io"
	"net/http"

	"codeberg.org/darckfast/workers-go/platform/cloudflare/cache"
	"codeberg.org/darckfast/workers-go/platform/cloudflare/fetch"
)

var GET_CACHE = func(w http.ResponseWriter, r *http.Request) {
	c := cache.New()
	res, _ := c.Match(r, nil)
	xcache := "miss"
	if res == nil {
		w.Header().Add("x-cache", xcache)
		rs, _ := http.NewRequest("GET", "https://darckfast.com", nil)
		res, _ := fetch.NewClient().Do(rs)

		defer res.Body.Close()
		bodyBytes, _ := io.ReadAll(res.Body)
		dummyR := http.Response{
			Status: res.Status,
			Header: http.Header{
				"cache-control": []string{"max-age=1500"},
			},
			Body: io.NopCloser(bytes.NewBuffer(bodyBytes)),
		}
		_ = c.Put(r, &dummyR)
		w.Write(bodyBytes)
	} else {
		xcache = "hit"
		w.Header().Add("x-cache", xcache)
		defer res.Body.Close()
		bodyBytes, _ := io.ReadAll(res.Body)

		w.Write(bodyBytes)
		// _, _ = io.Copy(w, res.Body)
		// defer res.Body.Close()
	}
}
