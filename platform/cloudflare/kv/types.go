package kv

//easyjson:json
type GetOptions struct {
	Type     string `json:"type"`
	CacheTTL int    `json:"cacheTtl,omitempty"`
}

//easyjson:json
type StringWithMetadata struct {
	Metadata map[string]any `json:"metadata"`
	Value    string         `json:"value"`
}

//easyjson:json
type PutOptions struct {
	Metadata      map[string]any `json:"metadata,omitempty"`
	Expiration    int            `json:"expiration,omitempty"`
	ExpirationTTL int            `json:"expirationTtl,omitempty"`
}
