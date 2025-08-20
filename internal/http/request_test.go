//go:build js && wasm

package jshttp

import (
	"io"
	"net/http"
	"strings"
	"testing"

	jsclass "github.com/Darckfast/workers-go/internal/class"
	jsconv "github.com/Darckfast/workers-go/internal/conv"
	"github.com/stretchr/testify/assert"
)

func TestToRequest(t *testing.T) {
	rawStr := `{"_id":"68a12e82045aa4c97496a889","index":0,"guid":"08a17b8a-7fec-4999-9937-2e9d437b9f80","isActive":true,"balance":"$3,495.88","picture":"http://placehold.it/32x32","age":37,"eyeColor":"blue","name":"Sampson Sheppard","gender":"male","company":"YOGASM","email":"sampsonsheppard@yogasm.com","phone":"+1 (962) 571-3499","address":"850 Gatling Place, Grenelefe, Arkansas, 1371","about":"Id deserunt tempor est pariatur aliqua consectetur nisi veniam proident cillum. Sit fugiat eiusmod consequat aute incididunt sint est. Incididunt id tempor aliquip qui ipsum. Elit voluptate pariatur enim ullamco reprehenderit elit proident minim. Nostrud officia commodo quis adipisicing voluptate ipsum quis deserunt exercitation consequat sit id. Ex laborum ut ad aliquip officia ipsum nostrud est velit pariatur tempor. Enim mollit esse et non Lorem sit ullamco labore qui occaecat.\\r\\n","registered":"2023-08-28T08:39:25 +03:00","latitude":-58.806364,"longitude":-111.129464,"tags":["cillum","aute","duis","nostrud","irure","nulla","nulla"],"friends":[{"id":0,"name":"Deleon Black"},{"id":1,"name":"Maddox Wade"},{"id":2,"name":"Poole Bowman"}],"greeting":"Hello, Sampson Sheppard! You have 10 unread messages.","favoriteFruit":"strawberry"}`
	req := jsclass.Request.New("http://my-service.com/path", jsconv.MapToJSValue(map[string]any{
		"headers": map[string]string{
			"content-type": "application/json",
		},
		"method": "POST",
		"body":   rawStr,
	}))

	r := ToRequest(req)

	if r.Header.Get("content-type") != "application/json" {
		t.Fatalf("conversion is missing header: had %s expected %s", r.Header.Get("content-type"), "application/json")
	}

	b, _ := io.ReadAll(r.Body)
	bs := string(b)

	assert.Equal(t, rawStr, bs)
}

func TestToJSRequest(t *testing.T) {
	rawStr := `{"_id":"68a12e82045aa4c97496a889","index":0,"guid":"08a17b8a-7fec-4999-9937-2e9d437b9f80","isActive":true,"balance":"$3,495.88","picture":"http://placehold.it/32x32","age":37,"eyeColor":"blue","name":"Sampson Sheppard","gender":"male","company":"YOGASM","email":"sampsonsheppard@yogasm.com","phone":"+1 (962) 571-3499","address":"850 Gatling Place, Grenelefe, Arkansas, 1371","about":"Id deserunt tempor est pariatur aliqua consectetur nisi veniam proident cillum. Sit fugiat eiusmod consequat aute incididunt sint est. Incididunt id tempor aliquip qui ipsum. Elit voluptate pariatur enim ullamco reprehenderit elit proident minim. Nostrud officia commodo quis adipisicing voluptate ipsum quis deserunt exercitation consequat sit id. Ex laborum ut ad aliquip officia ipsum nostrud est velit pariatur tempor. Enim mollit esse et non Lorem sit ullamco labore qui occaecat.\\r\\n","registered":"2023-08-28T08:39:25 +03:00","latitude":-58.806364,"longitude":-111.129464,"tags":["cillum","aute","duis","nostrud","irure","nulla","nulla"],"friends":[{"id":0,"name":"Deleon Black"},{"id":1,"name":"Maddox Wade"},{"id":2,"name":"Poole Bowman"}],"greeting":"Hello, Sampson Sheppard! You have 10 unread messages.","favoriteFruit":"strawberry"}`
	reader := io.NopCloser(strings.NewReader(rawStr))

	r, _ := http.NewRequest("POST", "http://my-service/path", reader)
	r.Header.Set("content-type", "application/json")

	rjs := ToJSRequest(r)
	j, _ := jsclass.Await(rjs.Call("json"))

	js := jsclass.JSON.Call("stringify", j)
	assert.Equal(t, rawStr, js.String())
}
