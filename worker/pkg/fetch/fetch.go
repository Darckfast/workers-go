package fetchhandler

import (
	"bufio"
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/syumai/workers/cloudflare/cache"
	_ "github.com/syumai/workers/cloudflare/d1" // register driver
	durableobjects "github.com/syumai/workers/cloudflare/durable_objects"
	"github.com/syumai/workers/cloudflare/fetch"
	"github.com/syumai/workers/cloudflare/kv"
	"github.com/syumai/workers/cloudflare/queues"
	"github.com/syumai/workers/cloudflare/r2"
	"github.com/syumai/workers/cloudflare/sockets"
)

func New() {
	http.HandleFunc("GET /application/json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"vitest": true})
	})

	http.HandleFunc("POST /application/json", func(w http.ResponseWriter, r *http.Request) {
		payload := map[string]any{}
		defer r.Body.Close()
		json.NewDecoder(r.Body).Decode(&payload)

		b, _ := json.Marshal(payload)
		h := r.Header.Get("X-Test-Id")
		size := len(strconv.Quote(string(b)))
		w.Header().Set("Content-Type", "application/json")
		result := map[string]any{"raw": string(b), "size": size, "header": h, "query": r.URL.Query().Encode()}
		json.NewEncoder(w).Encode(result)
	})

	http.HandleFunc("POST /application/x-www-form-urlencoded", func(w http.ResponseWriter, r *http.Request) {
		id := r.FormValue("id")
		alpha := r.FormValue("alpha")
		url := r.FormValue("url")
		name := r.FormValue("fullname")
		num := r.FormValue("number")

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id":       id,
			"alpha":    alpha,
			"url":      url,
			"fullname": name,
			"number":   num,
		})
	})
	http.HandleFunc("POST /multipart/form-data", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		f, fh, err := r.FormFile("img")
		defer f.Close()
		buf := bytes.NewBuffer(make([]byte, 0))
		io.Copy(buf, f)

		jsonStr := r.FormValue("json")
		var j map[string]any
		json.Unmarshal([]byte(jsonStr), &j)
		jb, _ := json.Marshal(j)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"has_error":    err != nil,
			"size":         fh.Size,
			"filename":     fh.Filename,
			"actual-size":  buf.Len(),
			"content-type": fh.Header.Get("content-type"),
			"json":         string(jb),
		})
	})

	http.HandleFunc("DELETE /kv", func(w http.ResponseWriter, r *http.Request) {
		namespace, _ := kv.NewNamespace("TEST_NAMESPACE")
		err := namespace.Delete("count")

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"has_error": err != nil})
	})

	http.HandleFunc("POST /kv", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		namespace, _ := kv.NewNamespace("TEST_NAMESPACE")

		countStr, _ := namespace.GetString("count", nil)
		count, _ := strconv.Atoi(countStr)

		nextCountStr := strconv.Itoa(count + 1)

		err := namespace.PutString("count", nextCountStr, nil)
		json.NewEncoder(w).Encode(map[string]any{"has_error": err != nil, "count": nextCountStr})
	})

	http.HandleFunc("GET /r2", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		bucket, _ := r2.NewBucket("TEST_BUCKET")

		result, err := bucket.Get("count")
		rawBody, _ := io.ReadAll(result.Body)
		b64 := base64.StdEncoding.EncodeToString(rawBody)

		json.NewEncoder(w).Encode(map[string]any{
			"has_error": err != nil,
			"result":    result,
			"body":      b64,
		})
	})

	http.HandleFunc("POST /r2", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		bucket, _ := r2.NewBucket("TEST_BUCKET")

		b64 := r.FormValue("b64")
		data, err := base64.StdEncoding.DecodeString(b64)

		reader := io.NopCloser(bytes.NewReader(data))
		result, err := bucket.Put("count", reader, int64(len(data)), nil)

		json.NewEncoder(w).Encode(map[string]any{
			"has_error": err != nil,
			"result":    result,
		})
	})

	http.HandleFunc("GET /cache", func(w http.ResponseWriter, r *http.Request) {
		c := cache.New()

		res, err := c.Match(r, nil)

		xcache := "miss"
		if res == nil {
			w.Header().Add("x-cache", xcache)
			rs, _ := fetch.NewRequest(r.Context(), "GET", "https://google.com", nil)
			res, _ = fetch.NewClient().Do(rs, nil)

			tee := io.TeeReader(res.Body, w)
			dummyR := http.Response{
				Status: res.Status,
				Header: http.Header{
					"cache-control": []string{"max-age=1500"}},
				Body: io.NopCloser(tee),
			}
			err = c.Put(r, &dummyR)

			if err != nil {
				fmt.Println(err)
			}
		} else {
			xcache = "hit"
			w.Header().Add("x-cache", xcache)
			io.Copy(w, res.Body)
		}
		defer res.Body.Close()
	})

	http.HandleFunc("GET /d1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		db, _ := sql.Open("d1", "DB")

		result := db.QueryRow("SELECT current_timestamp")
		fmt.Println(result)

		var a any
		result.Scan(&a)
		json.NewEncoder(w).Encode(map[string]any{
			"result": a,
		})
	})

	http.HandleFunc("GET /queue", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		namespace, _ := kv.NewNamespace("TEST_NAMESPACE")
		result, _ := namespace.GetString("queue:result", nil)

		if result == "<null>" {
			w.WriteHeader(404)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"result": result,
		})
	})

	http.HandleFunc("POST /queue", func(w http.ResponseWriter, r *http.Request) {
		q, err := queues.NewProducer("TEST_QUEUE")
		fmt.Println("queue", err, q)
		content, _ := io.ReadAll(r.Body)
		err = q.SendText(string(content))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(202)
		json.NewEncoder(w).Encode(map[string]any{
			"has_error": err != nil,
		})
	})

	http.HandleFunc("GET /socket", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := sockets.Connect(r.Context(), "tcpbin.com:4242", nil)
		defer conn.Close()
		conn.SetDeadline(time.Now().Add(1 * time.Hour))
		conn.Write([]byte("hello.\n"))
		rd := bufio.NewReader(conn)
		bts, _ := rd.ReadBytes('.')
		w.Write(bts)
	})

	http.HandleFunc("GET /tail", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		namespace, _ := kv.NewNamespace("TEST_NAMESPACE")
		result, _ := namespace.GetString("tail:result", nil)

		if result == "<null>" {
			w.WriteHeader(404)
		}

		json.NewEncoder(w).Encode(map[string]any{
			"result": result,
		})
	})

	http.HandleFunc("GET /container", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		c, err := durableobjects.GetContainer("GO_CONTAINER", "test")

		rs, _ := c.ContainerFetch(r)

		json.NewEncoder(w).Encode(map[string]any{
			"has_error": err != nil,
			"result":    rs.Status,
		})
	})

	http.HandleFunc("GET /do", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		n, _ := durableobjects.NewDurableObjectNamespace("TEST_DO")
		objId := n.IdFromName("id")
		stub, _ := n.Get(objId)

		rs, err := stub.Call("sayHello")

		json.NewEncoder(w).Encode(map[string]any{
			"has_error": err != nil,
			"result":    rs,
		})
	})

	fetch.ServeNonBlock(nil)
}
