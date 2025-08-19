//go:build js && wasm

package httpsocket

import (
	"bufio"
	"net/http"
	"time"

	_ "github.com/Darckfast/workers-go/cloudflare/d1" // register driver

	"github.com/Darckfast/workers-go/cloudflare/sockets"
)

var GET_SOCKET_TCPBIN = func(w http.ResponseWriter, r *http.Request) {
	conn, _ := sockets.Connect(r.Context(), "tcpbin.com:4242", nil)
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(1 * time.Hour))
	conn.Write([]byte("hello.\n"))
	rd := bufio.NewReader(conn)
	bts, _ := rd.ReadBytes('.')
	w.Write(bts)
}
