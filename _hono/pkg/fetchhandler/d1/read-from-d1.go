//go:build js && wasm

package httpd1

import (
	"log"
	"net/http"

	"github.com/Darckfast/workers-go/cloudflare/d1/v2"
	"github.com/mailru/easyjson"
	"go.opentelemetry.io/contrib/bridges/otelslog"
)

var GET_D1_TOTAL = func(w http.ResponseWriter, r *http.Request) {
	var logger = otelslog.NewLogger("indexer-go")
	w.Header().Set("Content-Type", "application/json")

	logger.InfoContext(r.Context(), "test")
	db, err := d1.GetDB("DB")
	if err != nil {
		m := GenericMap{}
		m["error"] = err.Error()

		easyjson.MarshalToHTTPResponseWriter(&m, w)
		return
	}

	defer func() {
		_ = r.Body.Close()
	}()

	result, err := db.Prepare(`SELECT COUNT(id) FROM messages`).
		Run()

	if err != nil {
		m := GenericMap{}
		m["error"] = err.Error()

		easyjson.MarshalToHTTPResponseWriter(&m, w)
		return
	}

	easyjson.MarshalToHTTPResponseWriter(result, w)
}

var GET_D1 = func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	q := r.URL.Query().Get("q")

	db, err := d1.GetDB("DB")
	if err != nil {
		m := GenericMap{}
		m["error"] = err.Error()

		easyjson.MarshalToHTTPResponseWriter(&m, w)
		return
	}

	defer func() {
		_ = r.Body.Close()
	}()

	result, err := db.Prepare(`SELECT id, content 
		FROM messages_fts WHERE messages_fts MATCH ? LIMIT 20`).
		Bind(q).
		Run()
	if err != nil {
		m := GenericMap{}
		m["error"] = err.Error()

		easyjson.MarshalToHTTPResponseWriter(&m, w)
		return
	}

	easyjson.MarshalToHTTPResponseWriter(result, w)
}

var POST_D1 = func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db, err := d1.GetDB("DB")

	if err != nil {
		m := GenericMap{}
		m["error"] = err.Error()
		log.Println(err.Error())

		w.WriteHeader(500)
		easyjson.MarshalToHTTPResponseWriter(&m, w)
		return
	}
	defer r.Body.Close()

	var messages Messages
	err = easyjson.UnmarshalFromReader(r.Body, &messages)

	if err != nil {
		m := GenericMap{}
		m["error"] = err.Error()
		log.Println(err.Error())

		w.WriteHeader(500)
		easyjson.MarshalToHTTPResponseWriter(&m, w)
		return
	}
	var stmts []d1.D1PreparedStatment
	for _, message := range messages {
		stmts = append(stmts, *db.Prepare(`INSERT INTO messages (
    message_id, type, content, timestamp, edited_timestamp,
    flags, pinned, mention_everyone, tts,
    channel_id, channel_name, is_bot, author_id
) VALUES (
    ?, ?, ?, ?,
    ?, ?, ?, ?, ?,
    ?, ?, ?, ?
)`).Bind(message.ID, message.Type, message.Content, message.Timestamp, message.EditedTimestamp,
			message.Flags, message.Pinned, message.MentionEveryone, message.Tts, message.ChannelID,
			message.ChannelName, message.IsBot, message.Author.ID))

		if len(message.Embeds) > 0 {
			for _, em := range message.Embeds {
				stmts = append(stmts, *db.Prepare(`INSERT INTO embed_images (
                    message_id,
                    url, 
                    proxy_url, 
                    width, 
                    height, 
                    flags
                ) VALUES (?, ?, ?, ?, ?, ?)`).
					Bind(message.ID, em.Thumbnail.URL, em.Thumbnail.ProxyURL,
						em.Thumbnail.Width, em.Thumbnail.Height, em.Thumbnail.Flags))
			}
		}
	}

	result, err := db.Batch(stmts)
	if err != nil {
		m := GenericMap{}
		m["error"] = err.Error()
		log.Println(err.Error())

		w.WriteHeader(500)
		easyjson.MarshalToHTTPResponseWriter(&m, w)
		return
	}

	easyjson.MarshalToHTTPResponseWriter(result, w)
	w.WriteHeader(http.StatusCreated)
}
