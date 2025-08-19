//go:build js && wasm

package d1

import (
	"errors"
)

var (
	ErrDatabaseNotFound = errors.New("d1: database not found")
)
