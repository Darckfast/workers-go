//go:build js && wasm

package d1

import (
	"database/sql"
	"errors"
	"syscall/js"
)

type result struct {
	resultObj js.Value
}

var (
	_ sql.Result = (*result)(nil)
)

// LastInsertId returns id of result's last row.
// If lastRowId can't be retrieved, this method returns error.
func (r *result) LastInsertId() (int64, error) {
	v := r.resultObj.Get("meta").Get("last_row_id")
	if v.IsNull() || v.IsUndefined() {
		return 0, errors.New("d1: lastRowId cannot be retrieved")
	}
	id := v.Int()
	return int64(id), nil
}

func (r *result) RowsAffected() (int64, error) {
	return int64(r.resultObj.Get("changes").Int()), nil
}
