package d1

import (
	"context"
	"database/sql/driver"
	"syscall/js"

	jsclass "github.com/syumai/workers/internal/class"
)

type Connector struct {
	dbObj js.Value
}

var (
	_ driver.Connector = (*Connector)(nil)
)

// OpenConnector returns Connector of D1.
// This method checks DB existence. If DB was not found, this function returns error.
func OpenConnector(name string) (driver.Connector, error) {
	v := jsclass.Env.Get(name)
	if v.IsUndefined() {
		return nil, ErrDatabaseNotFound
	}
	return &Connector{dbObj: v}, nil
}

// Connect returns Conn of D1.
// This method doesn't check DB existence, so this function never return errors.
func (c *Connector) Connect(context.Context) (driver.Conn, error) {
	return &Conn{dbObj: c.dbObj}, nil
}

func (c *Connector) Driver() driver.Driver {
	return &Driver{}
}
