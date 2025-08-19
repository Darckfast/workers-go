//go:build js && wasm

package d1

import (
	"database/sql"
	"errors"
	"syscall/js"

	jsconv "github.com/Darckfast/workers-go/internal/conv"
)

type result struct {
	resultObj js.Value
}

var (
	_ sql.Result = (*result)(nil)
)

func (r *result) LastInsertId() (int64, error) {
	v := r.resultObj.Get("meta").Get("last_row_id")
	if v.IsNull() || v.IsUndefined() {
		return 0, errors.New("d1: lastRowId cannot be retrieved")
	}
	return jsconv.MaybeInt64(v), nil
}

func (r *result) RowsAffected() (int64, error) {
	return jsconv.MaybeInt64(r.resultObj.Get("meta").Get("changes")), nil
}
