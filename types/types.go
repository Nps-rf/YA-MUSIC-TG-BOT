package types

type ArtistInfo struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

type TrackInfo struct {
	Artists    []ArtistInfo `json:"artists"`
	Image      string       `json:"cover"`
	Title      string       `json:"title"`
	Link       string       `json:"link"`
	UpdateTime string
}
