package mongodb

import (
	"GoNews/pkg/storage"
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Конструктор объекта хранилища.
type Store struct {
	Client *mongo.Client
}

func New(conStr string) (*Store, error) {
	mongoOpts := options.Client().ApplyURI(conStr)
	client, err := mongo.Connect(context.Background(), mongoOpts)
	if err != nil {
		return &Store{}, err
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return &Store{}, err
	}
	s := Store{
		Client: client,
	}
	return &s, nil
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
