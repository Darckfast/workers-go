package jstail

import (
	"fmt"
	"net/http"
	"syscall/js"

	jsclass "github.com/syumai/workers/internal/class"
	jsconv "github.com/syumai/workers/internal/conv"
	jshttp "github.com/syumai/workers/internal/http"
)

type ScriptVersion struct {
	Id      string
	Tag     string
	Message string
}

type TraceItemTailEventInfoTailItem struct {
	ScriptName string
}

type TraceItemFetchEventInfoResponse struct {
	Status int
}
type TraceItemFetchEventInfoRequest struct {
	Cf      map[string]any
	Headers http.Header
	Method  string
	Url     string
}

type TraceItemEvent struct {
	Type string
	//rpc
	RpcMethod string
	//email
	MailFrom string
	RcptTo   string
	RawSize  int
	//queue
	Queue     string
	BatchSize int
	// cron and alarm
	ScheduledTime int64
	// cron
	Cron string
	// tail
	ConsumedEvents *[]TraceItemTailEventInfoTailItem
	// fetch
	Response *TraceItemFetchEventInfoResponse
	Request  *TraceItemFetchEventInfoRequest
	// websocket
	GetWebSocketEvent *TraceItemGetWebSocketEvent
}

type TraceItemGetWebSocketEvent struct {
	WebSocketEventType string
	Code               int
	WasClean           bool
}

type TraceLog struct {
	Timestamp int64
	Level     string
	Message   string
}

type TraceException struct {
	Timestamp int64
	Message   string
	Name      string
	Stack     string
}

type TraceDiagnosticeChannelEvent struct {
	Timestamp int64
	Channel   string
	Message   string
}

type TailItem struct {
	ScriptName               string
	Entrypoint               string
	Event                    *TraceItemEvent
	EventTimeStamp           int64
	Logs                     []TraceLog
	Exceptions               []TraceException
	DiagnosticsChannelEvents []TraceDiagnosticeChannelEvent
	Outcome                  string
	Truncated                bool
	CpuTime                  int64
	WallTime                 int64
	ExecutionModel           string
	ScriptTags               []string
	DispatchNamespace        string
	ScriptVersion            *ScriptVersion
}

type TailEvent struct {
	self   js.Value
	Type   string
	Events *[]TailItem
	Traces *[]TailItem
}

func (t *TailEvent) WailUntil(task func() error) {
	t.self.Call("waitUntil", jsclass.Promise.New(js.FuncOf(func(this js.Value, args []js.Value) any {
		resolve := args[0]
		reject := args[1]

		err := task()

		if err == nil {
			resolve.Invoke(true)
		} else {
			reject.Invoke(jsclass.ToJSError(err))
		}

		return nil
	})))
}

func parseTailItems(tracesJs js.Value) *[]TailItem {
	traces := []TailItem{}

	if !tracesJs.Truthy() {
		return &traces
	}

	for j := range tracesJs.Length() {
		traceJs := tracesJs.Index(j)
		fmt.Println(traceJs)
		tailItem := TailItem{
			ScriptName:               traceJs.Get("scriptName").String(),
			Entrypoint:               traceJs.Get("entrypoint").String(),
			EventTimeStamp:           jsconv.MaybeInt64(traceJs.Get("eventTimestamp")),
			Logs:                     []TraceLog{},
			Exceptions:               []TraceException{},
			Outcome:                  traceJs.Get("outcome").String(),
			Truncated:                traceJs.Get("truncated").Bool(),
			DiagnosticsChannelEvents: []TraceDiagnosticeChannelEvent{},
			ScriptTags:               []string{},
			DispatchNamespace:        traceJs.Get("dispatchNamespace").String(),
			CpuTime:                  jsconv.MaybeInt64(traceJs.Get("cpuTime")),
			WallTime:                 jsconv.MaybeInt64(traceJs.Get("wallTime")),
			ExecutionModel:           traceJs.Get("executionModel").String(),
		}

		scriptVerJs := traceJs.Get("scriptVersion")

		if scriptVerJs.Truthy() {
			tailItem.ScriptVersion = &ScriptVersion{
				Id:      scriptVerJs.Get("id").String(),
				Tag:     scriptVerJs.Get("tag").String(),
				Message: scriptVerJs.Get("message").String(),
			}
		}

		logsJs := traceJs.Get("logs")
		if logsJs.Truthy() {
			for li := range logsJs.Length() {
				item := logsJs.Index(li)

				tailItem.Logs = append(tailItem.Logs, TraceLog{
					Timestamp: jsconv.MaybeInt64(item.Get("timestamp")),
					Level:     item.Get("level").String(),
					Message:   item.Get("message").String(),
				})
			}
		}

		excJs := traceJs.Get("exceptions")
		if excJs.Truthy() {
			for li := range excJs.Length() {
				item := excJs.Index(li)

				tailItem.Exceptions = append(tailItem.Exceptions, TraceException{
					Timestamp: jsconv.MaybeInt64(item.Get("timestamp")),
					Name:      item.Get("name").String(),
					Stack:     item.Get("stack").String(),
					Message:   item.Get("message").String(),
				})
			}
		}

		diagJs := traceJs.Get("diagnosticsChannelEvents")
		if diagJs.Truthy() {
			for li := range diagJs.Length() {
				item := diagJs.Index(li)

				tailItem.DiagnosticsChannelEvents = append(tailItem.DiagnosticsChannelEvents, TraceDiagnosticeChannelEvent{
					Timestamp: jsconv.MaybeInt64(item.Get("timestamp")),
					Channel:   item.Get("channel").String(),
					Message:   item.Get("message").String(),
				})
			}
		}

		tagsJs := traceJs.Get("scriptTags")
		if tagsJs.Truthy() {
			tailItem.ScriptTags = jsconv.MaybeStringList(tagsJs)
		}

		eventJs := traceJs.Get("event")

		if eventJs.Truthy() {
			tailItem.Event = GetEvent(eventJs)
		}

		traces = append(traces, tailItem)
	}

	return &traces
}

func GetEvent(event js.Value) *TraceItemEvent {
	if !event.Truthy() {
		return nil
	}

	if event.Get("scheduledTime").Truthy() {
		if event.Get("cron").Truthy() {
			return &TraceItemEvent{
				Type:          "cron",
				Cron:          event.Get("cron").String(),
				ScheduledTime: jsconv.MaybeInt64(event.Get("scheduledTime")),
			}
		}

		return &TraceItemEvent{
			Type: "alarm",
			// In the Alarm types, it's defined as Date instead of number
			ScheduledTime: jsconv.DateToTimestamp(event.Get("scheduledTime")),
		}
	} else if event.Get("queue").Truthy() {
		return &TraceItemEvent{
			Type:      "queue",
			Queue:     event.Get("queue").String(),
			BatchSize: event.Get("batchSize").Int(),
		}
	} else if event.Get("mailFrom").Truthy() {
		return &TraceItemEvent{
			Type:     "email",
			MailFrom: event.Get("mailFrom").String(),
			RcptTo:   event.Get("rcptTo").String(),
			RawSize:  event.Get("rawSize").Int(),
		}
	} else if event.Get("consumedEvents").Truthy() {
		items := event.Get("consumedEvents")

		event := &TraceItemEvent{}
		for i := range items.Length() {
			*event.ConsumedEvents = append(*event.ConsumedEvents, TraceItemTailEventInfoTailItem{
				ScriptName: items.Index(i).Get("scriptName").String(),
			})
		}

		return event
	} else if event.Get("request").Truthy() {
		return &TraceItemEvent{
			Type: "fetch",
			Response: &TraceItemFetchEventInfoResponse{
				Status: event.Get("response").Get("status").Int(),
			},
			Request: &TraceItemFetchEventInfoRequest{
				Url:     event.Get("request").Get("url").String(),
				Method:  event.Get("request").Get("method").String(),
				Headers: jshttp.ToHeader(event.Get("request").Get("headers")),
			},
		}
	} else if event.Get("rpcMethod").Truthy() {
		return &TraceItemEvent{
			Type:      "rpc",
			RpcMethod: event.Get("rpcMethod").String(),
		}
	} else if event.Get("getWebSocketEvent").Truthy() {
		return &TraceItemEvent{
			Type: "websocket",
			GetWebSocketEvent: &TraceItemGetWebSocketEvent{
				WebSocketEventType: event.Get("getWebSocketEvent").Get("webSocketEventType").String(),
				Code:               jsconv.MaybeInt(event.Get("getWebSocketEvent").Get("code")),
				WasClean:           jsconv.MaybeBool(event.Get("getWebSocketEvent").Get("wasClean")),
			},
		}
	}

	return nil
}

func NewEvents(eventsJs js.Value) *[]TailEvent {
	tailEvents := []TailEvent{}

	for i := range eventsJs.Length() {
		event := eventsJs.Index(i)
		traces := parseTailItems(event.Get("traces"))
		events := parseTailItems(event.Get("events"))
		tailEvents = append(tailEvents, TailEvent{
			Events: events,
			Type:   event.Get("type").String(),
			Traces: traces,
		})
	}

	return &tailEvents
}
