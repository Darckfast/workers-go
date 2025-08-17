package jsconv

import (
	"math"
	"testing"

	jsclass "github.com/syumai/workers/internal/class"
)

func TestMaybeInt(t *testing.T) {
	obj := jsclass.Object.New()
	obj.Set("int32", math.MaxInt32)

	convInt32 := MaybeInt64(obj.Get("int32"))

	if convInt32 != math.MaxInt32 {
		t.Fatalf("convertion yielded different value from expected: had %d expected %d", convInt32, math.MaxInt32)
	}
}

func TestMaybeInt64(t *testing.T) {
	obj := jsclass.Object.New()
	obj.Set("int64", math.MaxInt64)

	convInt64 := MaybeInt64(obj.Get("int64"))

	if convInt64 != math.MaxInt64 {
		t.Fatalf("convertion yielded different value from expected: had %d expected %d", convInt64, math.MaxInt64)
	}
}

func BenchmarkObjToMap1(b *testing.B) {
	obj := jsclass.JSON.Call("parse", `{"_id":"68a12e82045aa4c97496a889","index":0,"guid":"08a17b8a-7fec-4999-9937-2e9d437b9f80","isActive":true,"balance":"$3,495.88","picture":"http://placehold.it/32x32","age":37,"eyeColor":"blue","name":"Sampson Sheppard","gender":"male","company":"YOGASM","email":"sampsonsheppard@yogasm.com","phone":"+1 (962) 571-3499","address":"850 Gatling Place, Grenelefe, Arkansas, 1371","about":"Id deserunt tempor est pariatur aliqua consectetur nisi veniam proident cillum. Sit fugiat eiusmod consequat aute incididunt sint est. Incididunt id tempor aliquip qui ipsum. Elit voluptate pariatur enim ullamco reprehenderit elit proident minim. Nostrud officia commodo quis adipisicing voluptate ipsum quis deserunt exercitation consequat sit id. Ex laborum ut ad aliquip officia ipsum nostrud est velit pariatur tempor. Enim mollit esse et non Lorem sit ullamco labore qui occaecat.\\r\\n","registered":"2023-08-28T08:39:25 +03:00","latitude":-58.806364,"longitude":-111.129464,"tags":["cillum","aute","duis","nostrud","irure","nulla","nulla"],"friends":[{"id":0,"name":"Deleon Black"},{"id":1,"name":"Maddox Wade"},{"id":2,"name":"Poole Bowman"}],"greeting":"Hello, Sampson Sheppard! You have 10 unread messages.","favoriteFruit":"strawberry"}`)

	for b.Loop() {
		StrRecordToMap(obj)
	}
}

func BenchmarkObjToMap2(b *testing.B) {
	obj := jsclass.JSON.Call("parse", `{"_id":"68a12e82045aa4c97496a889","index":0,"guid":"08a17b8a-7fec-4999-9937-2e9d437b9f80","isActive":true,"balance":"$3,495.88","picture":"http://placehold.it/32x32","age":37,"eyeColor":"blue","name":"Sampson Sheppard","gender":"male","company":"YOGASM","email":"sampsonsheppard@yogasm.com","phone":"+1 (962) 571-3499","address":"850 Gatling Place, Grenelefe, Arkansas, 1371","about":"Id deserunt tempor est pariatur aliqua consectetur nisi veniam proident cillum. Sit fugiat eiusmod consequat aute incididunt sint est. Incididunt id tempor aliquip qui ipsum. Elit voluptate pariatur enim ullamco reprehenderit elit proident minim. Nostrud officia commodo quis adipisicing voluptate ipsum quis deserunt exercitation consequat sit id. Ex laborum ut ad aliquip officia ipsum nostrud est velit pariatur tempor. Enim mollit esse et non Lorem sit ullamco labore qui occaecat.\\r\\n","registered":"2023-08-28T08:39:25 +03:00","latitude":-58.806364,"longitude":-111.129464,"tags":["cillum","aute","duis","nostrud","irure","nulla","nulla"],"friends":[{"id":0,"name":"Deleon Black"},{"id":1,"name":"Maddox Wade"},{"id":2,"name":"Poole Bowman"}],"greeting":"Hello, Sampson Sheppard! You have 10 unread messages.","favoriteFruit":"strawberry"}`)

	for b.Loop() {
		ObjToMap(obj)
	}
}

func BenchmarkMaybeInt64(b *testing.B) {
	obj := jsclass.Object.New()
	obj.Set("int64", math.MaxInt64)

	for b.Loop() {
		MaybeInt64(obj.Get("int64"))
	}
}

func BenchmarkMaybeInt(b *testing.B) {
	obj := jsclass.Object.New()
	obj.Set("int64", math.MaxInt64)

	for b.Loop() {
		MaybeInt64(obj.Get("int32"))
	}
}
