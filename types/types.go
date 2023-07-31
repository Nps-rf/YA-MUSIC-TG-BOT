package types

type ArtistInfo struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

type UserInfo struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

type TrackInfo struct {
	Artists    []ArtistInfo `json:"artists"`
	User       UserInfo     `json:"user"`
	Image      string       `json:"cover"`
	Title      string       `json:"title"`
	Link       string       `json:"link"`
	UpdateTime string
}
