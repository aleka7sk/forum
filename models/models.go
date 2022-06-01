package models

type User struct {
	Id       int
	Username string
	Password string
}

type Category struct {
	Id    int
	Title string
}

type Post struct {
	Id         int
	Title      string
	Author     string
	Content    string
	AuthorId   int
	CategoryId int
	Likes      int
	Dislikes   int
	// date    time.Duration // int64 time.Now().Unix
}

type Comment struct {
	Id      int
	Author  string
	Content string
	PostId  int
	UserId  int
}

type Vote struct {
	Id        int
	Condition int
	PostId    int
	UserId    int
}

// Overall
// Voting kind, uid, sid
