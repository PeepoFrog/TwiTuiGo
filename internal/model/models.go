package model

type AuthToTwitch struct {
	ClientID     string
	ClientSecret string
	AccessToken  string
}

// games structs
type Game struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	BoxArtUrl string `json:"box_art_url"`
	IGBDid    string `json:"igdb_id"`
}
type Games struct {
	Data       []Game `json:"data"`
	Pagination struct {
		Cursor string `json:"cursor"`
	} `json:"pagination"`
}

// streamers struct
type Streamer struct {
	ID           string `json:"id"`
	UserID       string `json:"user_id"`
	UserLogin    string `json:"user_login"`
	UserName     string `json:"user_name"`
	GameID       string `json:"game_id"`
	GameName     string `json:"game_name"`
	Type         string `json:"type"`
	Title        string `json:"title"`
	ViewerCount  int    `json:"viewer_count"`
	StartedAt    string `json:"started_at"`
	Language     string `json:"language"`
	ThumbnailURL string `json:"thumbnail_url"`
	TagIDs       string `json:"tag_ids"`
	IsMature     bool   `json:"data"`
}
type Streamers struct {
	Data       []Streamer `json:"data"`
	Pagination struct {
		Cursor string `json:"cursor"`
	} `json:"pagination"`
}
