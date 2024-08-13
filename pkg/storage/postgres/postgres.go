package postgres

import (
	"GoNews/pkg/storage"
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Хранилище данных.
type Store struct {
	Pool *pgxpool.Pool
}

// Конструктор, принимает строку подключения к БД.
func New(conStr string) (*Store, error) {
	db, err := pgxpool.Connect(context.Background(), conStr)
	if err != nil {
		return nil, err
	}
	s := Store{
		Pool: db,
	}
	return &s, nil
}

// Генератор тестовых данных
func (s *Store) TestData() error {
	request := `
	DROP TABLE IF EXISTS posts, authors;

	CREATE TABLE authors (
    	id SERIAL PRIMARY KEY,
    	name TEXT NOT NULL
	);

	CREATE TABLE posts (
    	id SERIAL PRIMARY KEY,
    	title TEXT  NOT NULL,
    	content TEXT NOT NULL,
    	author_id INTEGER REFERENCES authors(id) NOT NULL,
    	created_at BIGINT DEFAULT extract(epoch from now()),
    	published_at BIGINT DEFAULT 0
	);
	`
	_, err := s.Pool.Exec(context.Background(), request)
	if err != nil {
		return err
	}

	id, err := s.AddAuthor("Robert")
	if err != nil {
		return err
	}
	err = s.AddPost(storage.Post{
		Title:    "fist post",
		Content:  "content of first post",
		AuthorID: id,
	})
	if err != nil {
		return err
	}
	id, err = s.AddAuthor("John")
	if err != nil {
		return err
	}
	err = s.AddPost(storage.Post{
		Title:    "second post",
		Content:  "content of second post",
		AuthorID: id,
	})
	if err != nil {
		return err
	}
	return nil
}

// Возвращает все посты
func (s *Store) Posts() ([]storage.Post, error) {
	var posts []storage.Post
	rows, err := s.Pool.Query(
		context.Background(),
		`SELECT posts.id, title, content, author_id, name, created_at, published_at FROM posts, authors WHERE posts.author_id = authors.id ORDER BY id`,
	)
	if err != nil {
		return posts, err
	}
	for rows.Next() {
		var post storage.Post
		err = rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.AuthorID,
			&post.AuthorName,
			&post.CreatedAt,
			&post.PublishedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

// Добавление поста
func (s *Store) AddPost(post storage.Post) error {
	_, err := s.Pool.Exec(
		context.Background(),
		`INSERT INTO posts (title, content, author_id) VALUES ($1, $2, $3)`,
		post.Title,
		post.Content,
		post.AuthorID,
	)
	if err != nil {
		return err
	}
	return nil
}

// Обновление поста
func (s *Store) UpdatePost(post storage.Post) error {
	_, err := s.Pool.Exec(
		context.Background(),
		`UPDATE posts SET title = $1, content = $2, author_id = $3, published_at = $4 WHERE id = $5`,
		post.Title,
		post.Content,
		post.AuthorID,
		post.PublishedAt,
		post.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

// Удаление поста
func (s *Store) DeletePost(post storage.Post) error {
	_, err := s.Pool.Exec(
		context.Background(),
		`DELETE FROM posts WHERE id = $1`,
		post.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

// Добавление нового автора в хранилище
func (s *Store) AddAuthor(name string) (int, error) {
	var id int
	err := s.Pool.QueryRow(context.Background(), `INSERT INTO authors (name) VALUES($1) returning id;`, name).Scan(&id)
	if err != nil {
		return id, err
	}
	return id, nil
}
