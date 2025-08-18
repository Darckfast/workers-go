//go:build js && wasm

package d1

import (
	"context"
	"database/sql"
	"database/sql/driver"
)

func init() {
	sql.Register("d1", &Driver{})
}

type Driver struct{}

var (
	_ driver.Driver = (*Driver)(nil)
)

func (d *Driver) Open(name string) (driver.Conn, error) {
	connector, err := OpenConnector(name)
	if err != nil {
		return nil, err
	}
	return connector.Connect(context.Background())
}
