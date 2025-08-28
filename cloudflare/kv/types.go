package kv

//easyjson:json
type GetOptions struct {
	Type     string `json:"type"`
	CacheTTL int    `json:"cacheTtl,omitempty"`
}

//easyjson:json
type StringWithMetadata struct {
	Value    string         `json:"value"`
	Metadata map[string]any `json:"metadata"`
}

//easyjson:json
type PutOptions struct {
	Expiration    int            `json:"expiration,omitempty"`
	ExpirationTTL int            `json:"expirationTtl,omitempty"`
	Metadata      map[string]any `json:"metadata,omitempty"`
}
