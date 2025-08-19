//go:build js && wasm

package jshttp

import (
	"io"
	"net/http"
	"strings"
	"testing"

	jsclass "github.com/Darckfast/workers-go/internal/class"
	jsconv "github.com/Darckfast/workers-go/internal/conv"
)

// BenchmarkToJSResponse       4382            260885 ns/op
func BenchmarkToJSResponse(t *testing.B) {
	rawStr := `{"_id":"68a12e82045aa4c97496a889","index":0,"guid":"08a17b8a-7fec-4999-9937-2e9d437b9f80","isActive":true,"balance":"$3,495.88","picture":"http://placehold.it/32x32","age":37,"eyeColor":"blue","name":"Sampson Sheppard","gender":"male","company":"YOGASM","email":"sampsonsheppard@yogasm.com","phone":"+1 (962) 571-3499","address":"850 Gatling Place, Grenelefe, Arkansas, 1371","about":"Id deserunt tempor est pariatur aliqua consectetur nisi veniam proident cillum. Sit fugiat eiusmod consequat aute incididunt sint est. Incididunt id tempor aliquip qui ipsum. Elit voluptate pariatur enim ullamco reprehenderit elit proident minim. Nostrud officia commodo quis adipisicing voluptate ipsum quis deserunt exercitation consequat sit id. Ex laborum ut ad aliquip officia ipsum nostrud est velit pariatur tempor. Enim mollit esse et non Lorem sit ullamco labore qui occaecat.\\r\\n","registered":"2023-08-28T08:39:25 +03:00","latitude":-58.806364,"longitude":-111.129464,"tags":["cillum","aute","duis","nostrud","irure","nulla","nulla"],"friends":[{"id":0,"name":"Deleon Black"},{"id":1,"name":"Maddox Wade"},{"id":2,"name":"Poole Bowman"}],"greeting":"Hello, Sampson Sheppard! You have 10 unread messages.","favoriteFruit":"strawberry"}`
	reader := io.NopCloser(strings.NewReader(rawStr))
	rs := http.Response{
		Header: http.Header{
			"content-type":      []string{"application/json"},
			"Accept":            []string{"application/json", "text/plain, */*"},
			"Authorization":     []string{"Bearer fake-token-abc123"},
			"X-Requested-With":  []string{"XMLHttpRequest"},
			"Cache-Control":     []string{"no-cache"},
			"Pragma":            []string{"no-cache"},
			"X-App-Version":     []string{"2.5.1"},
			"X-Client-Platform": []string{"web"},
			"X-Debug-Flag":      []string{"true"},
			"X-Session-ID":      []string{"session-xyz-987"},
			"X-Feature-Flag":    []string{"enable-new-ui"},
			"X-Fake-Header":     []string{"this-is-fake"},
			"X-Tracking-ID":     []string{"track-1234-5678"},
			"X-Fake-User-Agent": []string{"Mozilla/5.0 FakeBrowser/1.0"},
			"X-Fake-Referer":    []string{"https://fake.example.com"},
			"X-Fake-IP":         []string{"192.0.2.123"},
		},
		StatusCode: 201,
		Status:     "201 Created",
		Body:       reader,
	}

	for t.Loop() {
		ToJSResponse(&rs)
	}
}

// BenchmarkToResponse         5256            192393 ns/op
func BenchmarkToResponse(t *testing.B) {
	rawStr := `{"_id":"68a12e82045aa4c97496a889","index":0,"guid":"08a17b8a-7fec-4999-9937-2e9d437b9f80","isActive":true,"balance":"$3,495.88","picture":"http://placehold.it/32x32","age":37,"eyeColor":"blue","name":"Sampson Sheppard","gender":"male","company":"YOGASM","email":"sampsonsheppard@yogasm.com","phone":"+1 (962) 571-3499","address":"850 Gatling Place, Grenelefe, Arkansas, 1371","about":"Id deserunt tempor est pariatur aliqua consectetur nisi veniam proident cillum. Sit fugiat eiusmod consequat aute incididunt sint est. Incididunt id tempor aliquip qui ipsum. Elit voluptate pariatur enim ullamco reprehenderit elit proident minim. Nostrud officia commodo quis adipisicing voluptate ipsum quis deserunt exercitation consequat sit id. Ex laborum ut ad aliquip officia ipsum nostrud est velit pariatur tempor. Enim mollit esse et non Lorem sit ullamco labore qui occaecat.\\r\\n","registered":"2023-08-28T08:39:25 +03:00","latitude":-58.806364,"longitude":-111.129464,"tags":["cillum","aute","duis","nostrud","irure","nulla","nulla"],"friends":[{"id":0,"name":"Deleon Black"},{"id":1,"name":"Maddox Wade"},{"id":2,"name":"Poole Bowman"}],"greeting":"Hello, Sampson Sheppard! You have 10 unread messages.","favoriteFruit":"strawberry"}`
	res := jsclass.Response.New(rawStr, jsconv.MapToJSValue(map[string]any{
		"status": 201,
		"headers": map[string]any{
			"Content-Type":      "application/json",
			"Accept":            "application/json, text/plain, */*",
			"Authorization":     "Bearer fake-token-abc123",
			"X-Requested-With":  "XMLHttpRequest",
			"Cache-Control":     "no-cache",
			"Pragma":            "no-cache",
			"X-App-Version":     "2.5.1",
			"X-Client-Platform": "web",
			"X-Debug-Flag":      "true",
			"X-Session-ID":      "session-xyz-987",
			"X-Feature-Flag":    "enable-new-ui",
			"X-Fake-Header":     "this-is-fake",
			"X-Tracking-ID":     "track-1234-5678",
			"X-Fake-User-Agent": "Mozilla/5.0 FakeBrowser/1.0",
			"X-Fake-Referer":    "https://fake.example.com",
			"X-Fake-IP":         "192.0.2.123",
		},
	}))

	for t.Loop() {
		ToResponse(res)
	}
}
