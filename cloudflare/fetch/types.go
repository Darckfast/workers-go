package fetch

//easyjson:json
type InitOptions struct {
	Redirect    string
	Credentials string
}

//easyjson:json
type RequestInitCF struct {
	Apps             bool           `json:"apps"`
	CacheEverything  bool           `json:"cacheEverything"`
	CacheKey         string         `json:"cacheKey"`
	CacheTags        []string       `json:"cacheTags"`
	CacheTtl         int            `json:"cacheTtl"`
	CacheTtlByStatus map[string]int `json:"cacheTtlByStatus"`
	Mirage           bool           `json:"mirage"`
	Polish           string         `json:"polish"`
	ResolveOverride  string         `json:"resolveOverride"`
	ScrapShield      bool           `json:"scrapShield"`
	Webp             bool           `json:"webp"`
}
