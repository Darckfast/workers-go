package jshttp

//easyjson:json
type JSRequest struct {
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
	URL     string            `json:"url"`
}

// //easyjson:json
// type JSResponse struct {}
