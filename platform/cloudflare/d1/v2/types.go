package d1

//easyjson:json
type D1ExecResult struct {
	Count    int `json:"count"`
	Duration int `json:"duration"`
}

//easyjson:json
type Meta struct {
	ServedBy       string `json:"served_by"`
	ServedByRegion string `json:"served_by_region"`
	Timings        struct {
		SQLDurationMs float64 `json:"sql_duration_ms"`
	} `json:"timings"`
	Duration        float64 `json:"duration"`
	Changes         int64   `json:"changes"`
	LastRowID       int64   `json:"last_row_id"`
	SizeAfter       int64   `json:"size_after"`
	RowsRead        int64   `json:"rows_read"`
	RowsWritten     int64   `json:"rows_written"`
	ServedByPrimary bool    `json:"served_by_primary"`
	ChangedDB       bool    `json:"changed_db"`
}

//easyjson:json
type D1Result struct {
	ResultsString string
	Results       []any `json:"results"`
	Meta          Meta  `json:"meta"`
	Success       bool  `json:"success"`
}

//easyjson:json
type D1BatchResults []D1Result

//easyjson:json
type D1FirstResult map[string]any

//easyjson:json
type D1RawResults []map[string]any
