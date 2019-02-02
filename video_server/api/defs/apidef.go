package defs

//request
type UserCredential struct {
	UserName string `json:"user_name"`
	Pwd      string `json:"pwd"`
}

type NewComment struct {
	AuthorId int    `json:"author_id"`
	Content  string `json:"content"`
}

type NewVideo struct {
	AuthorId int    `json:"author_id"`
	Name     string `json:"name"`
}

//response
type SignedUp struct {
	Success   bool   `json:"success"`
	SessionId string `json:"session_id"`
}

type UserSession struct {
	UserName  string `json:"user_name"`
	SessionId string `json:"session_id"`
}

type UserInfo struct {
	Id int `json:"id"`
}

type SignedIn struct {
	Success   bool   `json:"success"`
	SessionId string `json:"session_id"`
}

type VideosInfo struct {
	Videos []*VideoInfo `json:"videos"`
}

type Comments struct {
	Comments []*Comment `json:"comments"`
}

//data model
type User struct {
	Id        int
	LoginName string
	Pwd       string
}

type VideoInfo struct {
	Id          string `json:"id"`
	AuthorId    int    `json:"author_id"`
	Name        string `json:"name"`
	DisplayTime string `json:"display_time"`
}

type Comment struct {
	Id      string `json:"id"`
	VideoId string `json:"video_id"`
	Author  string `json:"author"`
	Content string `json:"content"`
}

type SimpleSession struct {
	UserName string
	TTL      int64
}
