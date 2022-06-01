package postmodels

type Post struct {
	Id       int
	Title    string
	Author   string
	Content  string
	AuthorId string
	Likes    int
	Dislikes int
	// date    time.Duration // int64 time.Now().Unix
}
