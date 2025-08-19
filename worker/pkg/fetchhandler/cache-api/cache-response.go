//go:build js && wasm

package httpcache

import (
	"io"
	"net/http"

	_ "github.com/Darckfast/workers-go/cloudflare/d1" // register driver

	"github.com/Darckfast/workers-go/cloudflare/cache"
	"github.com/Darckfast/workers-go/cloudflare/fetch"
)

var GET_CACHE = func(w http.ResponseWriter, r *http.Request) {
	c := cache.New()

	res, _ := c.Match(r, nil)

	xcache := "miss"
	if res == nil {
		w.Header().Add("x-cache", xcache)
		rs, _ := http.NewRequest("GET", "https://darckfast.com", nil)
		res, _ := fetch.NewClient().Do(rs)

		tee := io.TeeReader(res.Body, w)
		dummyR := http.Response{
			Status: res.Status,
			Header: http.Header{
				"cache-control": []string{"max-age=1500"},
			},
			Body: io.NopCloser(tee),
		}
		_ = c.Put(r, &dummyR)
	} else {
		xcache = "hit"
		w.Header().Add("x-cache", xcache)
		io.Copy(w, res.Body)
	}
	// There might be a concurrency problem due pull being a promise
	// res.Body sometimes is nil when this function exits, but the returned body
	// is correct
	// defer res.Body.Close()
}
