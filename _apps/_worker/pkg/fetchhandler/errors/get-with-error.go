//go:build js && wasm

package errorshandler

import (
	"errors"
	"log"
	"net/http"
	"time"
)

var GET_ERROR = func(w http.ResponseWriter, r *http.Request) {
	log.Println("this is a log", time.Now().Nanosecond())

	e := errors.New("this is a error")

	log.Println("my error", e)

	panic(e)
}
