package postmodels

type Post struct {
	Id        int
	Title     string
	Author    string
	Content   string
	AuthorId  int
	Likes     int
	Dislikes  int
	Condition int
	// date    time.Duration // int64 time.Now().Unix
}
