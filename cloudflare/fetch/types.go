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
	Image            *CFImage       `json:"image,omitempty"`
	OriginAuth       string         `json:"origin-auth,omitempty"`
}

//easyjson:json
type GravityCoordinates struct {
	X float64 `json:"x,omitempty"`
	Y float64 `json:"y,omitempty"`
}

//easyjson:json
type BorderObject struct {
	Color  string `json:"color"`
	Width  int    `json:"width,omitempty"`
	Top    int    `json:"top,omitempty"`
	Right  int    `json:"right,omitempty"`
	Bottom int    `json:"bottom,omitempty"`
	Left   int    `json:"left,omitempty"`
}

//easyjson:json
type TrimBorder struct {
	Color     string  `json:"color,omitempty"`
	Tolerance float64 `json:"tolerance,omitempty"`
	Keep      float64 `json:"keep,omitempty"`
}

//easyjson:json
type Trim struct {
	Top    int         `json:"top,omitempty"`
	Bottom int         `json:"bottom,omitempty"`
	Left   int         `json:"left,omitempty"`
	Right  int         `json:"right,omitempty"`
	Width  int         `json:"width,omitempty"`
	Height int         `json:"height,omitempty"`
	Border interface{} `json:"border,omitempty"`
}

//easyjson:json
type Draw struct {
	URL        string      `json:"url"`
	Opacity    float64     `json:"opacity,omitempty"`
	Repeat     interface{} `json:"repeat,omitempty"`
	Top        int         `json:"top,omitempty"`
	Left       int         `json:"left,omitempty"`
	Bottom     int         `json:"bottom,omitempty"`
	Right      int         `json:"right,omitempty"`
	Width      int         `json:"width,omitempty"`
	Height     int         `json:"height,omitempty"`
	Fit        string      `json:"fit,omitempty"`
	Gravity    interface{} `json:"gravity,omitempty"`
	Background string      `json:"background,omitempty"`
	Rotate     int         `json:"rotate,omitempty"`
}

//easyjson:json
type R2 struct {
	BucketColoId int `json:"bucketColoId,omitempty"`
}

//easyjson:json
type Minify struct {
	Javascript bool `json:"javascript"`
	Css        bool `json:"css"`
	Html       bool `json:"html"`
}

//easyjson:json
type CFImage struct {
	Width       int         `json:"width,omitempty"`
	Height      int         `json:"height,omitempty"`
	Fit         string      `json:"fit,omitempty"`
	Gravity     interface{} `json:"gravity,omitempty"`
	Background  string      `json:"background,omitempty"`
	Rotate      int         `json:"rotate,omitempty"`
	Dpr         float64     `json:"dpr,omitempty"`
	Trim        interface{} `json:"trim,omitempty"`
	Quality     interface{} `json:"quality,omitempty"`
	Format      string      `json:"format,omitempty"`
	Anim        bool        `json:"anim,omitempty"`
	Metadata    string      `json:"metadata,omitempty"`
	Sharpen     float64     `json:"sharpen,omitempty"`
	Blur        float64     `json:"blur,omitempty"`
	Draw        []Draw      `json:"draw,omitempty"`
	Border      interface{} `json:"border,omitempty"`
	Brightness  float64     `json:"brightness,omitempty"`
	Contrast    float64     `json:"contrast,omitempty"`
	Gamma       float64     `json:"gamma,omitempty"`
	Saturation  float64     `json:"saturation,omitempty"`
	Flip        string      `json:"flip,omitempty"`
	Compression string      `json:"compression,omitempty"`
	Minify      Minify      `json:"minify,omitempty"`
	Mirage      bool        `json:"mirage,omitempty"`
	Polish      string      `json:"polish,omitempty"`
	R2          R2          `json:"r2,omitempty"`
}
