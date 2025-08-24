//go:build js && wasm

package httpsocket

import (
	"bufio"
	"net/http"
	"time"

	"github.com/Darckfast/workers-go/cloudflare/sockets"
)

var GET_SOCKET_TCPBIN = func(w http.ResponseWriter, r *http.Request) {
	conn, _ := sockets.Connect(r.Context(), "tcpbin.com:4242", nil)
	defer func() {
		_ = conn.Close()
	}()
	_ = conn.SetDeadline(time.Now().Add(1 * time.Hour))
	_, _ = conn.Write([]byte("hello.\n"))
	rd := bufio.NewReader(conn)
	bts, _ := rd.ReadBytes('.')
	_, _ = w.Write(bts)
}
