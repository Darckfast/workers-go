package queues

//easyjson:json
type QueueSendResult struct {
	Metadata Metadata `json:"metadata"`
}

//easyjson:json
type Metadata struct {
	Metrics QueueMetrics `json:"metrics"`
}

//easyjson:json
type QueueMetrics struct {
	BacklogCount           int64 `json:"backlogCount"`
	BacklogBytes           int64 `json:"backlogBytes"`
	OldestMessageTimestamp int64 `json:"oldestMessageTimestamp"`
}
