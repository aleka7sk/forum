package models

type User struct {
	Id       int
	Username string
	Password string
}

type Post struct {
	Id      int
	Title   string
	Author  string
	Content string
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

// Overall
// Voting kind, uid, sid
