package tail

//easyjson:json
type ScriptVersion struct {
	Id      string `json:"id,omitempty"`
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
	Url     string            `json:"url,omitempty"`
}

//easyjson:json
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

//easyjson:json
type TraceItemGetWebSocketEvent struct {
	WebSocketEventType string `json:"webSocketEventType,omitempty"`
	Code               int    `json:"code,omitempty"`
	WasClean           bool   `json:"wasClean,omitempty"`
}

//easyjson:json
type TraceLog struct {
	Timestamp int64  `json:"timestamp,omitempty"`
	Level     string `json:"level,omitempty"`
	Message   []any  `json:"message"`
}

//easyjson:json
type TraceException struct {
	Timestamp int64  `json:"timestamp,omitempty"`
	Message   string `json:"message,omitempty"`
	Name      string `json:"name,omitempty"`
	Stack     string `json:"stack,omitempty"`
}

//easyjson:json
type TraceDiagnosticeChannelEvent struct {
	Timestamp int64  `json:"timestamp,omitempty"`
	Channel   string `json:"channel,omitempty"`
	Message   string `json:"message,omitempty"`
}

//easyjson:json
type TraceItem struct {
	ScriptName               string                         `json:"scriptName"`
	Entrypoint               string                         `json:"entrypoint,omitempty"`
	Event                    *TraceItemEvent                `json:"event,omitempty"`
	EventTimeStamp           int64                          `json:"eventTimestamp,omitempty"`
	Logs                     []TraceLog                     `json:"logs"`
	Exceptions               []TraceException               `json:"exceptions"`
	DiagnosticsChannelEvents []TraceDiagnosticeChannelEvent `json:"diagnosticsChannelEvents"`
	Outcome                  string                         `json:"outcome,omitempty"`
	Truncated                bool                           `json:"truncated"`
	CpuTime                  int64                          `json:"cpuTime"`
	WallTime                 int64                          `json:"wallTime"`
	ExecutionModel           string                         `json:"executionModel,omitempty"`
	ScriptTags               []string                       `json:"scriptTags,omitempty"`
	DispatchNamespace        string                         `json:"dispatchNamespace,omitempty"`
	ScriptVersion            *ScriptVersion                 `json:"scriptVersion,omitempty"`
}

//easyjson:json
type Traces []TraceItem
