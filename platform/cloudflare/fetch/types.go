package fetch

//easyjson:json
type InitOptions struct {
	Redirect    string
	Credentials string
}

//easyjson:json
type RequestInitCF struct {
	CacheTTLByStatus map[string]int `json:"cacheTtlByStatus"`
	Image            *CFImage       `json:"image,omitempty"`
	ResolveOverride  string         `json:"resolveOverride"`
	OriginAuth       string         `json:"origin-auth,omitempty"`
	CacheKey         string         `json:"cacheKey"`
	Polish           string         `json:"polish"`
	CacheTags        []string       `json:"cacheTags"`
	CacheTTL         int            `json:"cacheTtl"`
	Mirage           bool           `json:"mirage"`
	ScrapShield      bool           `json:"scrapShield"`
	Webp             bool           `json:"webp"`
	Apps             bool           `json:"apps"`
	CacheEverything  bool           `json:"cacheEverything"`
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
	Border Border `json:"border,omitempty"`
	Top    int    `json:"top,omitempty"`
	Bottom int    `json:"bottom,omitempty"`
	Left   int    `json:"left,omitempty"`
	Right  int    `json:"right,omitempty"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
}

//easyjson:json
type Draw struct {
	Fit        string  `json:"fit,omitempty"`
	Background string  `json:"background,omitempty"`
	Repeat     string  `json:"repeat,omitempty"`
	URL        string  `json:"url"`
	Gravity    string  `json:"gravity,omitempty"`
	Top        int     `json:"top,omitempty"`
	Right      int     `json:"right,omitempty"`
	Width      int     `json:"width,omitempty"`
	Height     int     `json:"height,omitempty"`
	Bottom     int     `json:"bottom,omitempty"`
	Left       int     `json:"left,omitempty"`
	Opacity    float64 `json:"opacity,omitempty"`
	Rotate     int     `json:"rotate,omitempty"`
}

//easyjson:json
type R2 struct {
	BucketColoID int `json:"bucketColoId,omitempty"`
}

//easyjson:json
type Minify struct {
	Javascript bool `json:"javascript"`
	CSS        bool `json:"css"`
	HTML       bool `json:"html"`
}

//easyjson:json
type CFImage struct {
	Metadata    string  `json:"metadata,omitempty"`
	Polish      string  `json:"polish,omitempty"`
	Fit         string  `json:"fit,omitempty"`
	Gravity     string  `json:"gravity,omitempty"`
	Background  string  `json:"background,omitempty"`
	Compression string  `json:"compression,omitempty"`
	Flip        string  `json:"flip,omitempty"`
	Format      string  `json:"format,omitempty"`
	Draw        []Draw  `json:"draw,omitempty"`
	Border      Border  `json:"border,omitempty"`
	Trim        TrimImg `json:"trim,omitempty"`
	Sharpen     float64 `json:"sharpen,omitempty"`
	Gamma       float64 `json:"gamma,omitempty"`
	Blur        float64 `json:"blur,omitempty"`
	R2          R2      `json:"r2,omitempty"`
	Quality     Quality `json:"quality,omitempty"`
	Brightness  float64 `json:"brightness,omitempty"`
	Contrast    float64 `json:"contrast,omitempty"`
	Width       int     `json:"width,omitempty"`
	Saturation  float64 `json:"saturation,omitempty"`
	Dpr         float64 `json:"dpr,omitempty"`
	Rotate      int     `json:"rotate,omitempty"`
	Height      int     `json:"height,omitempty"`
	Minify      Minify  `json:"minify,omitempty"`
	Mirage      bool    `json:"mirage,omitempty"`
	Anim        bool    `json:"anim,omitempty"`
}

//easyjson:json
type Quality struct {
	Quality int `json:"quality"`
}

//easyjson:json
type TrimImg struct {
	Top    int `json:"top"`
	Right  int `json:"right"`
	Bottom int `json:"bottom"`
	Left   int `json:"left"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

//easyjson:json
type Border struct {
	Color  string `json:"color"`
	Top    int    `json:"top"`
	Right  int    `json:"right"`
	Bottom int    `json:"bottom"`
	Left   int    `json:"left"`
}
