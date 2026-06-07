package tail

//easyjson:json
type ScriptVersion struct {
	ID      string `json:"id,omitempty"`
	Tag     string `json:"tag,omitempty"`
	Message string `json:"message,omitempty"`
}

//easyjson:json
type TraceItemTailEventInfoTailItem struct {
	ScriptName string `json:"scriptName,omitempty"`
}

//easyjson:json
type TraceItemFetchEventInfoResponse struct {
	Status int `json:"status,omitempty"`
}

//easyjson:json
type TraceItemFetchEventInfoRequest struct {
	Cf      map[string]any    `json:"cf,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
	Method  string            `json:"method,omitempty"`
	URL     string            `json:"url,omitempty"`
}

//easyjson:json
type TraceItemEvent struct {
	ConsumedEvents    *[]TraceItemTailEventInfoTailItem `json:"consumedEvents,omitempty"`
	GetWebSocketEvent *TraceItemGetWebSocketEvent       `json:"getWebSocketEvent,omitempty"`
	Request           *TraceItemFetchEventInfoRequest   `json:"request,omitempty"`
	Response          *TraceItemFetchEventInfoResponse  `json:"response,omitempty"`
	RcptTo            string                            `json:"rcptTo,omitempty"`
	Queue             string                            `json:"queue,omitempty"`
	Cron              string                            `json:"cron,omitempty"`
	Type              string                            `json:"-"`
	MailFrom          string                            `json:"mailFrom,omitempty"`
	RPCMethod         string                            `json:"rpcMethod,omitempty"`
	BatchSize         int                               `json:"batchSize,omitempty"`
	ScheduledTime     int64                             `json:"scheduledTime,omitempty"`
	RawSize           int                               `json:"rawSize,omitempty"`
}

//easyjson:json
type TraceItemGetWebSocketEvent struct {
	WebSocketEventType string `json:"webSocketEventType,omitempty"`
	Code               int    `json:"code,omitempty"`
	WasClean           bool   `json:"wasClean,omitempty"`
}

//easyjson:json
type TraceLog struct {
	Level     string `json:"level,omitempty"`
	Message   []any  `json:"message"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

//easyjson:json
type TraceException struct {
	Message   string `json:"message,omitempty"`
	Name      string `json:"name,omitempty"`
	Stack     string `json:"stack,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

//easyjson:json
type TraceDiagnosticeChannelEvent struct {
	Channel   string `json:"channel,omitempty"`
	Message   string `json:"message,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

//easyjson:json
type TraceItem struct {
	Event                    *TraceItemEvent                `json:"event,omitempty"`
	ScriptVersion            *ScriptVersion                 `json:"scriptVersion,omitempty"`
	ExecutionModel           string                         `json:"executionModel,omitempty"`
	Entrypoint               string                         `json:"entrypoint,omitempty"`
	ScriptName               string                         `json:"scriptName"`
	Outcome                  string                         `json:"outcome,omitempty"`
	DispatchNamespace        string                         `json:"dispatchNamespace,omitempty"`
	Logs                     []TraceLog                     `json:"logs"`
	Exceptions               []TraceException               `json:"exceptions"`
	DiagnosticsChannelEvents []TraceDiagnosticeChannelEvent `json:"diagnosticsChannelEvents"`
	ScriptTags               []string                       `json:"scriptTags,omitempty"`
	EventTimeStamp           int64                          `json:"eventTimestamp,omitempty"`
	WallTime                 int64                          `json:"wallTime"`
	CPUTime                  int64                          `json:"cpuTime"`
	Truncated                bool                           `json:"truncated"`
}

//easyjson:json
type Traces []TraceItem
