package types

type ArtistInfo struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

type RequestData struct {
	UserInfo  UserInfo  `json:"userInfo"`
	TrackInfo TrackInfo `json:"trackInfo"`
}

type UserInfo struct {
	Id       string `json:"userId"`
	Username string `json:"userName"`
}

type TrackInfo struct {
	Artists    []ArtistInfo `json:"artists"`
	Image      string       `json:"cover"`
	Title      string       `json:"title"`
	Link       string       `json:"link"`
	Duration   float64      `json:"duration"`
	Liked      bool         `json:"liked"`
	Disliked   bool         `json:"disliked"`
	UpdateTime string
}

type BotConfig struct {
	Token   string
	Debug   bool
	Timeout int
}
