package jshttp

import (
	"net/http"
	"strings"
	"testing"

	jsclass "github.com/syumai/workers/internal/class"
)

func TestToHeaders(t *testing.T) {
	headers := jsclass.Headers.New()

	headers.Call("append", "Content-Type", "application/json")
	headers.Call("append", "Accept", "application/json, text/plain, */*")
	headers.Call("append", "Authorization", "Bearer fake-token-abc123")
	headers.Call("append", "X-Requested-With", "XMLHttpRequest")
	headers.Call("append", "Cache-Control", "no-cache")
	headers.Call("append", "Pragma", "no-cache")
	headers.Call("append", "X-App-Version", "2.5.1")
	headers.Call("append", "X-Client-Platform", "web")
	headers.Call("append", "X-Debug-Flag", "true")
	headers.Call("append", "X-Session-ID", "session-xyz-987")
	headers.Call("append", "X-Feature-Flag", "enable-new-ui")
	headers.Call("append", "X-Fake-Header", "this-is-fake")
	headers.Call("append", "X-Tracking-ID", "track-1234-5678")
	headers.Call("append", "X-Fake-User-Agent", "Mozilla/5.0 FakeBrowser/1.0")
	headers.Call("append", "X-Fake-Referer", "https://fake.example.com")
	headers.Call("append", "X-Fake-IP", "192.0.2.123")

	h := ToHeader(headers)

	for key, value := range h {
		hv := strings.Join(value, ",")
		jsh := headers.Call("get", key).String()
		if hv != jsh {
			t.Fatalf("conversion yielded wrong value: had '%s' expected '%s' on '%s'", hv, jsh, key)
		}
	}
}

func TestToJSHeaders(t *testing.T) {
	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("Accept", "application/json, text/plain, */*")
	headers.Set("Authorization", "Bearer fake-token-abc123")
	headers.Set("X-Requested-With", "XMLHttpRequest")
	headers.Set("Cache-Control", "no-cache")
	headers.Set("Pragma", "no-cache")
	headers.Set("X-App-Version", "2.5.1")
	headers.Set("X-Client-Platform", "web")
	headers.Set("X-Debug-Flag", "true")
	headers.Set("X-Session-ID", "session-xyz-987")
	headers.Set("X-Feature-Flag", "enable-new-ui")
	headers.Set("X-Fake-Header", "this-is-fake")
	headers.Set("X-Tracking-ID", "track-1234-5678")
	headers.Set("X-Fake-User-Agent", "Mozilla/5.0 FakeBrowser/1.0")
	headers.Set("X-Fake-Referer", "https://fake.example.com")
	headers.Set("X-Fake-IP", "192.0.2.123")

	h := ToJSHeader(headers)

	for key, value := range headers {
		hv := strings.Join(value, ",")
		jsh := h.Call("get", key).String()
		if hv != jsh {
			t.Fatalf("conversion yielded wrong value: had '%s' expected '%s' on '%s'", hv, jsh, key)
		}
	}
}
