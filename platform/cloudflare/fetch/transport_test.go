//go:build js && wasm

package fetch

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCFTransport_Get_ShouldTimeout(t *testing.T) {
	c := http.Client{
		Timeout:   1,
		Transport: DefaultCFTransport,
	}

	r, err := http.NewRequest("GET", "http://localhost", nil)

	assert.Nil(t, err)
	_, err = c.Do(r)
	assert.ErrorIs(t, err, &TimeoutError{})
}

func TestCFTransport_Get_OK(t *testing.T) {
	c := http.Client{
		Transport: DefaultCFTransport,
	}

	r, err := http.NewRequest("GET", "http://localhost", nil)

	assert.Nil(t, err)
	_, err = c.Do(r)
	assert.Equal(t, err.Error(), `Get "http://localhost": JavaScript error: fetch failed`)
}

func TestCFTransport_Post_ShouldTimeout(t *testing.T) {
	c := http.Client{
		Timeout:   1,
		Transport: DefaultCFTransport,
	}

	r, err := http.NewRequest("POST", "http://localhost", bytes.NewBuffer([]byte("my payload")))

	assert.Nil(t, err)
	_, err = c.Do(r)
	assert.ErrorIs(t, err, &TimeoutError{})
}

func TestCFTransport_Post_OK(t *testing.T) {
	c := http.Client{
		Transport: DefaultCFTransport,
	}

	r, err := http.NewRequest("POST", "http://localhost", bytes.NewBuffer([]byte("my payload")))

	assert.Nil(t, err)
	_, err = c.Do(r)
	assert.Equal(t, err.Error(), `Post "http://localhost": JavaScript error: fetch failed`)
}
