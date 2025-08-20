//go:build js && wasm

package d1

import (
	"context"
	"database/sql/driver"
	"errors"
	"syscall/js"

	jsclass "github.com/Darckfast/workers-go/internal/class"
)

type stmt struct {
	stmtObj js.Value
}

var (
	_ driver.Stmt             = (*stmt)(nil)
	_ driver.StmtExecContext  = (*stmt)(nil)
	_ driver.StmtQueryContext = (*stmt)(nil)
)

func (s *stmt) Close() error {
	// do nothing
	return nil
}

// NumInput is not supported and always returns -1.
func (s *stmt) NumInput() int {
	return -1
}

func (s *stmt) Exec([]driver.Value) (driver.Result, error) {
	return nil, errors.New("d1: Exec is deprecated and not implemented")
}

// TODO: implement batch
// TODO: implement withSession and D1DatabaseSession

// ExecContext executes prepared statement.
// Given []driver.NamedValue's `Name` field will be ignored because Cloudflare D1 client doesn't support it.
func (s *stmt) ExecContext(_ context.Context, args []driver.NamedValue) (driver.Result, error) {
	argValues := make([]any, len(args))
	for i, arg := range args {
		if src, ok := arg.Value.([]byte); ok {
			dst := jsclass.Uint8Array.New(len(src))
			if n := js.CopyBytesToJS(dst, src); n != len(src) {
				return nil, errors.New("incomplete copy into Uint8Array")
			}
			argValues[i] = dst
		} else {
			argValues[i] = arg.Value
		}
	}
	resultPromise := s.stmtObj.Call("bind", argValues...).Call("run")
	resultObj, err := jsclass.Await(resultPromise)
	if err != nil {
		return nil, err
	}
	return &result{
		resultObj: resultObj,
	}, nil
}

func (s *stmt) Query([]driver.Value) (driver.Rows, error) {
	return nil, errors.New("d1: Query is deprecated and not implemented")
}

func (s *stmt) QueryContext(_ context.Context, args []driver.NamedValue) (driver.Rows, error) {
	argValues := make([]any, len(args))
	for i, arg := range args {
		if src, ok := arg.Value.([]byte); ok {
			dst := jsclass.Uint8Array.New(len(src))
			if n := js.CopyBytesToJS(dst, src); n != len(src) {
				return nil, errors.New("incomplete copy into Uint8Array")
			}
			argValues[i] = dst
		} else {
			argValues[i] = arg.Value
		}
	}
	resultPromise := s.stmtObj.Call("bind", argValues...).Call("raw", map[string]any{"columnNames": true})
	rowsArray, err := jsclass.Await(resultPromise)
	if err != nil {
		return nil, err
	}
	// If there are no rows to retrieve, length is 0.
	if rowsArray.Length() == 0 {
		return &rows{
			_columns:  nil,
			rowsArray: rowsArray,
		}, nil
	}

	// First item of rowsArray is column names
	colsArray := rowsArray.Call("shift")
	colsLen := colsArray.Length()
	cols := make([]string, colsLen)
	for i := range colsLen {
		cols[i] = colsArray.Index(i).String()
	}
	return &rows{
		_columns:  cols,
		rowsArray: rowsArray,
	}, nil
}

// {
//   success: boolean, // true if the operation was successful, false otherwise
//   meta: {
//     served_by: string // the version of Cloudflare's backend Worker that returned the result
//     served_by_region: string // the region of the database instance that executed the query
//     served_by_primary: boolean // true if (and only if) the database instance that executed the query was the primary
//     timings: {
//       sql_duration_ms: number // the duration of the SQL query execution by the database instance (not including any network time)
//     }
//     duration: number, // the duration of the SQL query execution only, in milliseconds
//     changes: number, // the number of changes made to the database
//     last_row_id: number, // the last inserted row ID, only applies when the table is defined without the `WITHOUT ROWID` option
//     changed_db: boolean, // true if something on the database was changed
//     size_after: number, // the size of the database after the query is successfully applied
//     rows_read: number, // the number of rows read (scanned) by this query
//     rows_written: number // the number of rows written by this query
//   }
//   results: array | null, // [] if empty, or null if it does not apply
// }
