package memdb

import (
	"GoNews/pkg/storage"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// Хранилище данных.
type Store struct {
	posts        map[int]string
	authors      map[int]string
	lastPostID   int
	lastAuthorID int
	idMute       *sync.Mutex
}

// Конструктор объекта хранилища.
func New() *Store {
	return &Store{
		posts:        make(map[int]string),
		authors:      make(map[int]string),
		lastPostID:   0,
		lastAuthorID: 0,
		idMute:       &sync.Mutex{},
	}
}

func (s *Store) TestData() error {
	id := s.AddAuthor("Mike")
	err := s.AddPost(storage.Post{
		Title:     "fist post",
		Content:   "content of first post",
		AuthorID:  id,
		CreatedAt: time.Now().Unix(),
	})
	if err != nil {
		return err
	}
	id = s.AddAuthor("Ted")
	err = s.AddPost(storage.Post{
		Title:     "second post",
		Content:   "content of second post",
		AuthorID:  id,
		CreatedAt: time.Now().Unix(),
	})
	if err != nil {
		return err
	}
	return nil
}

// Возвращает все посты
func (s *Store) Posts() ([]storage.Post, error) {
	var posts []storage.Post
	var post storage.Post
	// Итерируем по map posts для выбора всех постов
	for _, p := range s.posts {
		err := json.Unmarshal([]byte(p), &post)
		if err != nil {
			return posts, err
		}
		posts = append(posts, post)
	}
	// Итерируем по полученному слайсу постов для добавления имени автора из map authors
	for i, p := range posts {
		posts[i].AuthorName = s.authors[p.AuthorID]
	}
	return posts, nil
}

// Добавление поста
func (s *Store) AddPost(post storage.Post) error {
	s.idMute.Lock()
	defer s.idMute.Unlock()
	p, err := json.Marshal(post)
	if err != nil {
		return err
	}
	s.lastPostID++
	s.posts[s.lastPostID] = string(p)
	return nil
}

// Обновление поста
func (s *Store) UpdatePost(post storage.Post) error {
	if _, exists := s.posts[post.ID]; !exists {
		return fmt.Errorf("post with id %d not exists", post.ID)
	}
	p, err := json.Marshal(post)
	if err != nil {
		return err
	}
	s.posts[post.ID] = string(p)
	return nil
}

// Удаление поста
func (s *Store) DeletePost(post storage.Post) error {
	if _, exists := s.posts[post.ID]; !exists {
		return fmt.Errorf("post with id %d not exists", post.ID)
	}
	delete(s.posts, post.ID)
	return nil
}

// Добавление нового автора в хранилище
func (s *Store) AddAuthor(name string) int {
	s.idMute.Lock()
	defer s.idMute.Unlock()
	s.lastAuthorID++
	s.authors[s.lastAuthorID] = name
	return s.lastAuthorID
}
