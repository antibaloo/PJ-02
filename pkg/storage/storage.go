package storage

// Post - публикация.
type Post struct {
	ID          int    `json:"id" bson:"_id"` //добавлена нотация bson для работы с mongodb
	Title       string `json:"title" bson:"title"`
	Content     string `json:"content" bson:"content"`
	AuthorID    int    `json:"author_id" bson:"author_id"`
	AuthorName  string `json:"author_name" bson:"-"`
	CreatedAt   int64  `json:"created_at" bson:"created_at"`
	PublishedAt int64  `json:"published_at" bson:"published_at"`
}

// Interface задаёт контракт на работу с БД.
type Interface interface {
	Posts() ([]Post, error) // получение всех публикаций
	AddPost(Post) error     // создание новой публикации
	UpdatePost(Post) error  // обновление публикации
	DeletePost(Post) error  // удаление публикации по ID
}
