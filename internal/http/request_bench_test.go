package jshttp

import (
	"io"
	"net/http"
	"strings"
	"testing"

	jsclass "github.com/syumai/workers/internal/class"
	jsconv "github.com/syumai/workers/internal/conv"
)

func BenchmarkToRequestGET(b *testing.B) {
	req := jsclass.Request.New("http://my-service.com/path", jsconv.MapToJSValue(map[string]any{
		"headers": map[string]string{
			"content-type":      "application/json",
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
		"method": "GET",
	}))

	for b.Loop() {
		ToRequest(req)
	}
}

func BenchmarkToRequestPOST(b *testing.B) {
	rawStr := `{"_id":"68a12e82045aa4c97496a889","index":0,"guid":"08a17b8a-7fec-4999-9937-2e9d437b9f80","isActive":true,"balance":"$3,495.88","picture":"http://placehold.it/32x32","age":37,"eyeColor":"blue","name":"Sampson Sheppard","gender":"male","company":"YOGASM","email":"sampsonsheppard@yogasm.com","phone":"+1 (962) 571-3499","address":"850 Gatling Place, Grenelefe, Arkansas, 1371","about":"Id deserunt tempor est pariatur aliqua consectetur nisi veniam proident cillum. Sit fugiat eiusmod consequat aute incididunt sint est. Incididunt id tempor aliquip qui ipsum. Elit voluptate pariatur enim ullamco reprehenderit elit proident minim. Nostrud officia commodo quis adipisicing voluptate ipsum quis deserunt exercitation consequat sit id. Ex laborum ut ad aliquip officia ipsum nostrud est velit pariatur tempor. Enim mollit esse et non Lorem sit ullamco labore qui occaecat.\\r\\n","registered":"2023-08-28T08:39:25 +03:00","latitude":-58.806364,"longitude":-111.129464,"tags":["cillum","aute","duis","nostrud","irure","nulla","nulla"],"friends":[{"id":0,"name":"Deleon Black"},{"id":1,"name":"Maddox Wade"},{"id":2,"name":"Poole Bowman"}],"greeting":"Hello, Sampson Sheppard! You have 10 unread messages.","favoriteFruit":"strawberry"}`
	req := jsclass.Request.New("http://my-service.com/path", jsconv.MapToJSValue(map[string]any{
		"headers": map[string]string{
			"content-type":      "application/json",
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
		"method": "POST",
		"body":   rawStr,
	}))

	for b.Loop() {
		ToRequest(req)
	}
}

func BenchmarkToRequestV2POST(b *testing.B) {
	rawStr := `{"_id":"68a12e82045aa4c97496a889","index":0,"guid":"08a17b8a-7fec-4999-9937-2e9d437b9f80","isActive":true,"balance":"$3,495.88","picture":"http://placehold.it/32x32","age":37,"eyeColor":"blue","name":"Sampson Sheppard","gender":"male","company":"YOGASM","email":"sampsonsheppard@yogasm.com","phone":"+1 (962) 571-3499","address":"850 Gatling Place, Grenelefe, Arkansas, 1371","about":"Id deserunt tempor est pariatur aliqua consectetur nisi veniam proident cillum. Sit fugiat eiusmod consequat aute incididunt sint est. Incididunt id tempor aliquip qui ipsum. Elit voluptate pariatur enim ullamco reprehenderit elit proident minim. Nostrud officia commodo quis adipisicing voluptate ipsum quis deserunt exercitation consequat sit id. Ex laborum ut ad aliquip officia ipsum nostrud est velit pariatur tempor. Enim mollit esse et non Lorem sit ullamco labore qui occaecat.\\r\\n","registered":"2023-08-28T08:39:25 +03:00","latitude":-58.806364,"longitude":-111.129464,"tags":["cillum","aute","duis","nostrud","irure","nulla","nulla"],"friends":[{"id":0,"name":"Deleon Black"},{"id":1,"name":"Maddox Wade"},{"id":2,"name":"Poole Bowman"}],"greeting":"Hello, Sampson Sheppard! You have 10 unread messages.","favoriteFruit":"strawberry"}`
	req := jsclass.Request.New("http://my-service.com/path", jsconv.MapToJSValue(map[string]any{
		"headers": map[string]string{
			"content-type":      "application/json",
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
		"method": "POST",
		"body":   rawStr,
	}))

	for b.Loop() {
		ToRequestV2(req)
	}
}

func BenchmarkToRequestV2GET(b *testing.B) {
	req := jsclass.Request.New("http://my-service.com/path", jsconv.MapToJSValue(map[string]any{
		"headers": map[string]string{
			"content-type":      "application/json",
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
		"method": "GET",
	}))

	for b.Loop() {
		ToRequestV2(req)
	}
}

func BenchmarkToJSRequestGET(t *testing.B) {
	r, _ := http.NewRequest("GET", "http://my-service/path", nil)
	r.Header.Set("content-type", "application/json")
	r.Header.Set("Accept", "application/json, text/plain, */*")
	r.Header.Set("Authorization", "Bearer fake-token-abc123")
	r.Header.Set("X-Requested-With", "XMLHttpRequest")
	r.Header.Set("Cache-Control", "no-cache")
	r.Header.Set("Pragma", "no-cache")
	r.Header.Set("X-App-Version", "2.5.1")
	r.Header.Set("X-Client-Platform", "web")
	r.Header.Set("X-Debug-Flag", "true")
	r.Header.Set("X-Session-ID", "session-xyz-987")
	r.Header.Set("X-Feature-Flag", "enable-new-ui")
	r.Header.Set("X-Fake-Header", "this-is-fake")
	r.Header.Set("X-Tracking-ID", "track-1234-5678")
	r.Header.Set("X-Fake-User-Agent", "Mozilla/5.0 FakeBrowser/1.0")
	r.Header.Set("X-Fake-Referer", "https://fake.example.com")
	r.Header.Set("X-Fake-IP", "192.0.2.123")

	for t.Loop() {
		ToJSRequest(r)
	}
}

func BenchmarkToJSRequestPOST(t *testing.B) {
	rawStr := `{"_id":"68a12e82045aa4c97496a889","index":0,"guid":"08a17b8a-7fec-4999-9937-2e9d437b9f80","isActive":true,"balance":"$3,495.88","picture":"http://placehold.it/32x32","age":37,"eyeColor":"blue","name":"Sampson Sheppard","gender":"male","company":"YOGASM","email":"sampsonsheppard@yogasm.com","phone":"+1 (962) 571-3499","address":"850 Gatling Place, Grenelefe, Arkansas, 1371","about":"Id deserunt tempor est pariatur aliqua consectetur nisi veniam proident cillum. Sit fugiat eiusmod consequat aute incididunt sint est. Incididunt id tempor aliquip qui ipsum. Elit voluptate pariatur enim ullamco reprehenderit elit proident minim. Nostrud officia commodo quis adipisicing voluptate ipsum quis deserunt exercitation consequat sit id. Ex laborum ut ad aliquip officia ipsum nostrud est velit pariatur tempor. Enim mollit esse et non Lorem sit ullamco labore qui occaecat.\\r\\n","registered":"2023-08-28T08:39:25 +03:00","latitude":-58.806364,"longitude":-111.129464,"tags":["cillum","aute","duis","nostrud","irure","nulla","nulla"],"friends":[{"id":0,"name":"Deleon Black"},{"id":1,"name":"Maddox Wade"},{"id":2,"name":"Poole Bowman"}],"greeting":"Hello, Sampson Sheppard! You have 10 unread messages.","favoriteFruit":"strawberry"}`
	reader := io.NopCloser(strings.NewReader(rawStr))

	r, _ := http.NewRequest("POST", "http://my-service/path", reader)
	r.Header.Set("content-type", "application/json")
	r.Header.Set("Accept", "application/json, text/plain, */*")
	r.Header.Set("Authorization", "Bearer fake-token-abc123")
	r.Header.Set("X-Requested-With", "XMLHttpRequest")
	r.Header.Set("Cache-Control", "no-cache")
	r.Header.Set("Pragma", "no-cache")
	r.Header.Set("X-App-Version", "2.5.1")
	r.Header.Set("X-Client-Platform", "web")
	r.Header.Set("X-Debug-Flag", "true")
	r.Header.Set("X-Session-ID", "session-xyz-987")
	r.Header.Set("X-Feature-Flag", "enable-new-ui")
	r.Header.Set("X-Fake-Header", "this-is-fake")
	r.Header.Set("X-Tracking-ID", "track-1234-5678")
	r.Header.Set("X-Fake-User-Agent", "Mozilla/5.0 FakeBrowser/1.0")
	r.Header.Set("X-Fake-Referer", "https://fake.example.com")
	r.Header.Set("X-Fake-IP", "192.0.2.123")
	for t.Loop() {
		ToJSRequest(r)
	}
}

func BenchmarkToJSRequestV2POST(t *testing.B) {
	rawStr := `{"_id":"68a12e82045aa4c97496a889","index":0,"guid":"08a17b8a-7fec-4999-9937-2e9d437b9f80","isActive":true,"balance":"$3,495.88","picture":"http://placehold.it/32x32","age":37,"eyeColor":"blue","name":"Sampson Sheppard","gender":"male","company":"YOGASM","email":"sampsonsheppard@yogasm.com","phone":"+1 (962) 571-3499","address":"850 Gatling Place, Grenelefe, Arkansas, 1371","about":"Id deserunt tempor est pariatur aliqua consectetur nisi veniam proident cillum. Sit fugiat eiusmod consequat aute incididunt sint est. Incididunt id tempor aliquip qui ipsum. Elit voluptate pariatur enim ullamco reprehenderit elit proident minim. Nostrud officia commodo quis adipisicing voluptate ipsum quis deserunt exercitation consequat sit id. Ex laborum ut ad aliquip officia ipsum nostrud est velit pariatur tempor. Enim mollit esse et non Lorem sit ullamco labore qui occaecat.\\r\\n","registered":"2023-08-28T08:39:25 +03:00","latitude":-58.806364,"longitude":-111.129464,"tags":["cillum","aute","duis","nostrud","irure","nulla","nulla"],"friends":[{"id":0,"name":"Deleon Black"},{"id":1,"name":"Maddox Wade"},{"id":2,"name":"Poole Bowman"}],"greeting":"Hello, Sampson Sheppard! You have 10 unread messages.","favoriteFruit":"strawberry"}`
	reader := io.NopCloser(strings.NewReader(rawStr))

	r, _ := http.NewRequest("POST", "http://my-service/path", reader)
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json, text/plain, */*")
	r.Header.Set("Authorization", "Bearer fake-token-abc123")
	r.Header.Set("X-Requested-With", "XMLHttpRequest")
	r.Header.Set("Cache-Control", "no-cache")
	r.Header.Set("Pragma", "no-cache")
	r.Header.Set("X-App-Version", "2.5.1")
	r.Header.Set("X-Client-Platform", "web")
	r.Header.Set("X-Debug-Flag", "true")
	r.Header.Set("X-Session-ID", "session-xyz-987")
	r.Header.Set("X-Feature-Flag", "enable-new-ui")
	r.Header.Set("X-Fake-Header", "this-is-fake")
	r.Header.Set("X-Tracking-ID", "track-1234-5678")
	r.Header.Set("X-Fake-User-Agent", "Mozilla/5.0 FakeBrowser/1.0")
	r.Header.Set("X-Fake-Referer", "https://fake.example.com")
	r.Header.Set("X-Fake-IP", "192.0.2.123")

	for t.Loop() {
		ToJSRequestV2(r)
	}
}

func BenchmarkToJSRequestV2GET(t *testing.B) {
	r, _ := http.NewRequest("GET", "http://my-service/path", nil)
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json, text/plain, */*")
	r.Header.Set("Authorization", "Bearer fake-token-abc123")
	r.Header.Set("X-Requested-With", "XMLHttpRequest")
	r.Header.Set("Cache-Control", "no-cache")
	r.Header.Set("Pragma", "no-cache")
	r.Header.Set("X-App-Version", "2.5.1")
	r.Header.Set("X-Client-Platform", "web")
	r.Header.Set("X-Debug-Flag", "true")
	r.Header.Set("X-Session-ID", "session-xyz-987")
	r.Header.Set("X-Feature-Flag", "enable-new-ui")
	r.Header.Set("X-Fake-Header", "this-is-fake")
	r.Header.Set("X-Tracking-ID", "track-1234-5678")
	r.Header.Set("X-Fake-User-Agent", "Mozilla/5.0 FakeBrowser/1.0")
	r.Header.Set("X-Fake-Referer", "https://fake.example.com")
	r.Header.Set("X-Fake-IP", "192.0.2.123")

	for t.Loop() {
		ToJSRequestV2(r)
	}
}
