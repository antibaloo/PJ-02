package memdb

import "GoNews/pkg/storage"

// Хранилище данных.
type Store struct {
	posts        map[int]string
	authors      map[int]string
	lastPostID   int
	lastAuthorID int
}

// Конструктор объекта хранилища.
func New() *Store {
	return &Store{
		posts:        make(map[int]string),
		authors:      make(map[int]string),
		lastPostID:   0,
		lastAuthorID: 0,
	}
}

func (s *Store) Posts() ([]storage.Post, error) {
	var posts []storage.Post
	return posts, nil
}

func (s *Store) AddPost(storage.Post) error {
	return nil
}
func (s *Store) UpdatePost(storage.Post) error {
	return nil
}
func (s *Store) DeletePost(storage.Post) error {
	return nil
}
