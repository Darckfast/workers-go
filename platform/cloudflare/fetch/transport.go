//go:build js && wasm

package fetch

import (
	"net/http"

	"codeberg.org/darckfast/workers-go/internal/jsclass"
	"codeberg.org/darckfast/workers-go/internal/jshelper"
	"codeberg.org/darckfast/workers-go/internal/jshttp"
	"github.com/mailru/easyjson"
)

type Transport struct{}

var jsFetch jshelper.LazyJSVal

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	initOptions := InitOptions{
		Redirect:    "follow",
		Credentials: "omit",
	}

	initJSON, _ := easyjson.Marshal(initOptions)
	initObj, _ := jsclass.JSON.Parse(string(initJSON))

	ac := jsclass.AbortController.New()
	initObj.Set("signal", ac.Get("signal"))
	cf := req.Context().Value("cf")

	if cf != nil {
		cfJSON, _ := easyjson.Marshal(cf.(*RequestInitCF))
		cfObj, _ := jsclass.JSON.Parse(string(cfJSON))
		initObj.Set("cf", cfObj)
	}

	jsFetch.Init("fetch")

	var err error
	select {
	case <-req.Context().Done():
		_ = ac.Call("abort")
		return nil, &TimeoutError{}
	default:
	}

	promise := jsFetch.Invoke(
		jshttp.ToJSRequest(req),
		initObj,
	)

	jsRes, err := jsclass.Await(promise)
	if err != nil {
		return nil, err
	}

	return jshttp.ToResponse(jsRes), nil
}

type TimeoutError struct{}

func (t *TimeoutError) Error() string {
	return "Request timeout"
}

var _ http.RoundTripper = (*Transport)(nil)
var DefaultCFTransport http.RoundTripper = &Transport{}
