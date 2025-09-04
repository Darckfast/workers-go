package httpd1

//easyjson:json
type GenericMap map[string]any

//easyjson:json
type Messages []Message

//easyjson:json
type Message struct {
	Type         int           `json:"type"`
	Content      string        `json:"content,nocopy"`
	Mentions     []interface{} `json:"mentions"`
	MentionRoles []interface{} `json:"mention_roles"`
	Attachments  []interface{} `json:"attachments"`
	Embeds       []struct {
		Type        string `json:"type,nocopy"`
		Title       string `json:"title,nocopy"`
		Description string `json:"description,nocopy"`
		Color       int    `json:"color"`
		Image       struct {
			URL      string `json:"url,nocopy"`
			ProxyURL string `json:"proxy_url,nocopy"`
			Width    int    `json:"width"`
			Height   int    `json:"height"`
			Flags    int    `json:"flags"`
		} `json:"image"`
		Thumbnail struct {
			URL      string `json:"url,nocopy"`
			ProxyURL string `json:"proxy_url,nocopy"`
			Width    int    `json:"width"`
			Height   int    `json:"height"`
			Flags    int    `json:"flags"`
		} `json:"thumbnail"`
		Footer struct {
			Text string `json:"text,nocopy"`
		} `json:"footer"`
		ContentScanVersion int `json:"content_scan_version"`
	} `json:"embeds"`
	Timestamp       string        `json:"timestamp,nocopy"`
	EditedTimestamp string        `json:"edited_timestamp,nocopy"`
	Flags           int           `json:"flags"`
	Components      []interface{} `json:"components"`
	ID              string        `json:"id,nocopy"`
	ChannelID       string        `json:"channel_id,nocopy"`
	Author          struct {
		ID                   string      `json:"id,nocopy"`
		Username             string      `json:"username,nocopy"`
		Avatar               string      `json:"avatar,nocopy"`
		Discriminator        string      `json:"discriminator,nocopy"`
		PublicFlags          int         `json:"public_flags"`
		Flags                int         `json:"flags"`
		Bot                  bool        `json:"bot"`
		Banner               interface{} `json:"banner"`
		AccentColor          interface{} `json:"accent_color"`
		AvatarDecorationData interface{} `json:"avatar_decoration_data"`
		BannerColor          interface{} `json:"banner_color"`
		Clan                 interface{} `json:"clan"`
		PrimaryGuild         interface{} `json:"primary_guild"`
	} `json:"author"`
	Pinned          bool `json:"pinned"`
	MentionEveryone bool `json:"mention_everyone"`
	Tts             bool `json:"tts"`
	Reactions       []struct {
		Emoji struct {
			ID   string `json:"id,nocopy"`
			Name string `json:"name,nocopy"`
		} `json:"emoji"`
		Count        int `json:"count"`
		CountDetails struct {
			Burst  int `json:"burst"`
			Normal int `json:"normal"`
		} `json:"count_details"`
		BurstColors []interface{} `json:"burst_colors"`
		MeBurst     bool          `json:"me_burst"`
		BurstMe     bool          `json:"burst_me"`
		Me          bool          `json:"me"`
		BurstCount  int           `json:"burst_count"`
	} `json:"reactions"`
	ChannelName string `json:"channel_name,nocopy"`
	IsBot       bool   `json:"is_bot"`
}
