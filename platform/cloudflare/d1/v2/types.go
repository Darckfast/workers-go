package d1

//easyjson:json
type D1ExecResult struct {
	Count    int `json:"count"`
	Duration int `json:"duration"`
}

//easyjson:json
type D1Result struct {
	Success bool  `json:"success"`
	Results []any `json:"results"`
	Meta    struct {
		ServedBy        string `json:"served_by"`
		ServedByRegion  string `json:"served_by_region"`
		ServedByPrimary bool   `json:"served_by_primary"`
		Timings         struct {
			SqlDurationMs int64 `json:"sql_duration_ms"`
		} `json:"timings"`
		Duration    int64 `json:"duration"`
		Changes     int64 `json:"changes"`
		LastRowId   int64 `json:"last_row_id"`
		ChangedDb   bool  `json:"changed_db"`
		SizeAfter   int64 `json:"size_after"`
		RowsRead    int64 `json:"rows_read"`
		RowsWritten int64 `json:"rows_written"`
	} `json:"meta"`
}

//easyjson:json
type D1BatchResults []D1Result

//easyjson:json
type D1FirstResult map[string]any

//easyjson:json
type D1RawResults []map[string]any
