package storage

// Post - публикация.
type Post struct {
	ID          int    `json:"ID"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	AuthorID    int    `json:"author_id"`
	AuthorName  string `json:"-"`
	CreatedAt   int64  `json:"created_at"`
	PublishedAt int64  `json:"published_at"`
}

// Interface задаёт контракт на работу с БД.
type Interface interface {
	Posts() ([]Post, error) // получение всех публикаций
	AddPost(Post) error     // создание новой публикации
	UpdatePost(Post) error  // обновление публикации
	DeletePost(Post) error  // удаление публикации по ID
}
