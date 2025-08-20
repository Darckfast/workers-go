//go:build js && wasm

package jsconv

import (
	"encoding/json"
	"math"
	"testing"
	"time"

	jsclass "github.com/Darckfast/workers-go/internal/class"
	"github.com/stretchr/testify/assert"
)

func TestMaybeInt(t *testing.T) {
	obj := jsclass.Object.New()
	obj.Set("int32", math.MaxInt32)

	cInt := MaybeInt64(obj.Get("int32"))

	assert.Equal(t, int64(math.MaxInt32), cInt)
}

func TestMaybeInt64(t *testing.T) {
	obj := jsclass.Object.New()
	obj.Set("int64", math.MaxInt64)

	cInt := MaybeInt64(obj.Get("int64"))

	assert.Equal(t, int64(math.MaxInt64), cInt)
}

func TestMaybeInt64With32Plus1(t *testing.T) {
	obj := jsclass.Object.New()
	obj.Set("value", math.MaxInt32+1)

	cInt := MaybeInt64(obj.Get("value"))

	assert.Equal(t, int64(math.MaxInt32+1), cInt)
}

func TestJSValueToMap(t *testing.T) {
	jStr := `{"_id":"68a12e82045aa4c97496a889","index":0,"guid":"08a17b8a-7fec-4999-9937-2e9d437b9f80","isActive":true,"balance":"$3,495.88","picture":"http://placehold.it/32x32","age":37,"eyeColor":"blue","name":"Sampson Sheppard","gender":"male","company":"YOGASM","email":"sampsonsheppard@yogasm.com","phone":"+1 (962) 571-3499","address":"850 Gatling Place, Grenelefe, Arkansas, 1371","about":"Id deserunt tempor est pariatur aliqua consectetur nisi veniam proident cillum. Sit fugiat eiusmod consequat aute incididunt sint est. Incididunt id tempor aliquip qui ipsum. Elit voluptate pariatur enim ullamco reprehenderit elit proident minim. Nostrud officia commodo quis adipisicing voluptate ipsum quis deserunt exercitation consequat sit id. Ex laborum ut ad aliquip officia ipsum nostrud est velit pariatur tempor. Enim mollit esse et non Lorem sit ullamco labore qui occaecat.\\r\\n","registered":"2023-08-28T08:39:25 +03:00","latitude":-58.806364,"longitude":-111.129464,"tags":["cillum","aute","duis","nostrud","irure","nulla","nulla"],"friends":[{"id":0,"name":"Deleon Black"},{"id":1,"name":"Maddox Wade"},{"id":2,"name":"Poole Bowman"}],"greeting":"Hello, Sampson Sheppard! You have 10 unread messages.","favoriteFruit":"strawberry"}`
	obj := jsclass.JSON.Call("parse", jStr)

	mapValue, _ := JSValueToMap(obj)

	var j map[string]any
	_ = json.Unmarshal([]byte(jStr), &j)

	assert.Equal(t, j, mapValue)
	// if !reflect.DeepEqual(j, mapValue) {
	// 	t.Fatalf("conversion yielded different value from expected: had %s expected %s", mapValue, j)
	// }
}

func TestMaybeString(t *testing.T) {
	rawStr := "(❁´◡`❁)) this is my test string ☆*: .｡. o(≧▽≦)o .｡.:*☆)"
	obj := jsclass.String.Invoke(rawStr)
	str := MaybeString(obj)

	assert.Equal(t, rawStr, str)
}

func TestMaybeBool(t *testing.T) {
	obj := jsclass.Boolean.Invoke(true)
	b := MaybeBool(obj)

	assert.True(t, b)
}

func TestDateToEpoch(t *testing.T) {
	ts := time.Now().UnixMilli()
	obj := jsclass.Date.New(ts)
	jsts := DateToTimestamp(obj)

	assert.Equal(t, ts, jsts)
}

func TestDateToTime(t *testing.T) {
	n := time.Now()
	ts := n.UnixMilli()
	obj := jsclass.Date.New(ts)
	jsts := DateToTime(obj)

	assert.True(t, n.Equal(jsts))
}

func TestTimeToDate(t *testing.T) {
	now := time.Now()

	jsdate := TimeToDate(now)

	ms := MaybeInt64(jsdate.Call("getTime"))

	assert.Equal(t, now.UnixMilli(), ms)
}
