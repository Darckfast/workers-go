//go:build js && wasm

package tail

import (
	"net/http"
	"syscall/js"

	jsconv "github.com/Darckfast/workers-go/internal/conv"
	jshttp "github.com/Darckfast/workers-go/internal/http"
)

type ScriptVersion struct {
	Id      string `json:"id,omitempty"`
	Tag     string `json:"tag,omitempty"`
	Message string `json:"message,omitempty"`
}

type TraceItemTailEventInfoTailItem struct {
	ScriptName string `json:"scriptName,omitempty"`
}

type TraceItemFetchEventInfoResponse struct {
	Status int `json:"status,omitempty"`
}
type TraceItemFetchEventInfoRequest struct {
	Cf      map[string]any `json:"cf,omitempty"`
	Headers http.Header    `json:"headers,omitempty"`
	Method  string         `json:"method,omitempty"`
	Url     string         `json:"url,omitempty"`
}

type TraceItemEvent struct {
	Type string `json:"-"`
	//rpc
	RpcMethod string `json:"rpcMethod,omitempty"`
	//email
	MailFrom string `json:"mailFrom,omitempty"`
	RcptTo   string `json:"rcptTo,omitempty"`
	RawSize  int    `json:"rawSize,omitempty"`
	//queue
	Queue     string `json:"queue,omitempty"`
	BatchSize int    `json:"batchSize,omitempty"`
	// cron and alarm
	ScheduledTime int64 `json:"scheduledTime,omitempty"`
	// cron
	Cron string `json:"cron,omitempty"`
	// tail
	ConsumedEvents *[]TraceItemTailEventInfoTailItem `json:"consumedEvents,omitempty"`
	// fetch
	Response *TraceItemFetchEventInfoResponse `json:"response,omitempty"`
	Request  *TraceItemFetchEventInfoRequest  `json:"request,omitempty"`
	// websocket
	GetWebSocketEvent *TraceItemGetWebSocketEvent `json:"getWebSocketEvent,omitempty"`
}

type TraceItemGetWebSocketEvent struct {
	WebSocketEventType string `json:"webSocketEventType,omitempty"`
	Code               int    `json:"code,omitempty"`
	WasClean           bool   `json:"wasClean,omitempty"`
}

type TraceLog struct {
	Timestamp int64  `json:"timestamp,omitempty"`
	Level     string `json:"level,omitempty"`
	Message   string `json:"message,omitempty"`
}

type TraceException struct {
	Timestamp int64  `json:"timestamp,omitempty"`
	Message   string `json:"message,omitempty"`
	Name      string `json:"name,omitempty"`
	Stack     string `json:"stack,omitempty"`
}

type TraceDiagnosticeChannelEvent struct {
	Timestamp int64  `json:"timestamp,omitempty"`
	Channel   string `json:"channel,omitempty"`
	Message   string `json:"message,omitempty"`
}

type TraceItem struct {
	ScriptName               string                         `json:"scriptName,omitempty"`
	Entrypoint               string                         `json:"entrypoint,omitempty"`
	Event                    *TraceItemEvent                `json:"event,omitempty"`
	EventTimeStamp           int64                          `json:"eventTimestamp,omitempty"`
	Logs                     []TraceLog                     `json:"logs,omitempty"`
	Exceptions               []TraceException               `json:"exceptions,omitempty"`
	DiagnosticsChannelEvents []TraceDiagnosticeChannelEvent `json:"diagnosticsChannelEvents,omitempty"`
	Outcome                  string                         `json:"outcome,omitempty"`
	Truncated                bool                           `json:"truncated,omitempty"`
	CpuTime                  int64                          `json:"cpuTime,omitempty"`
	WallTime                 int64                          `json:"wallTime,omitempty"`
	ExecutionModel           string                         `json:"executionModel,omitempty"`
	ScriptTags               []string                       `json:"scriptTags,omitempty"`
	DispatchNamespace        string                         `json:"dispatchNamespace,omitempty"`
	ScriptVersion            *ScriptVersion                 `json:"scriptVersion,omitempty"`
}

func parseTailItems(tracesJs js.Value) []TraceItem {
	traces := []TraceItem{}

	if !tracesJs.Truthy() {
		return traces
	}

	// jsconv.JSValueToMap(tracesJs)
	for j := range tracesJs.Length() {
		traceJs := tracesJs.Index(j)
		tailItem := TraceItem{
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

	return traces
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

func NewEvents(eventsJs js.Value) *[]TraceItem {
	events := parseTailItems(eventsJs)
	return &events
}
