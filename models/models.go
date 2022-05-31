package models

type User struct {
	Id       int
	Username string
	Password string
}

type Post struct {
	Id       int
	Title    string
	Author   string
	Content  string
	AuthorId string
	// date    time.Duration // int64 time.Now().Unix
}

type Category struct {
	Id    int
	Title string
}

type Comments struct {
	Id      int
	Author  string
	Content string
}

type Emotion struct {
	Like    int
	Dislike int
	PostId  int
	UserId  int
}

// Overall
// Voting kind, uid, sid
