//go:build js && wasm

package d1

import (
	"context"
	"database/sql/driver"
	"errors"
	"syscall/js"
)

type Conn struct {
	dbObj js.Value
}

var (
	_ driver.Conn               = (*Conn)(nil)
	_ driver.ConnBeginTx        = (*Conn)(nil)
	_ driver.ConnPrepareContext = (*Conn)(nil)
)

func (c *Conn) Prepare(query string) (driver.Stmt, error) {
	stmtObj := c.dbObj.Call("prepare", query)
	return &stmt{
		stmtObj: stmtObj,
	}, nil
}

func (c *Conn) PrepareContext(_ context.Context, query string) (driver.Stmt, error) {
	return c.Prepare(query)
}

func (c *Conn) Close() error {
	// do nothing
	return nil
}

func (c *Conn) Begin() (driver.Tx, error) {
	return nil, errors.New("d1: Begin is deprecated and not implemented")
}

func (c *Conn) BeginTx(context.Context, driver.TxOptions) (driver.Tx, error) {
	return nil, errors.New("d1: transaction is not currently supported")
}
