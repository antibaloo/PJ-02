package mongodb

import (
	"GoNews/pkg/storage"
	"context"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Константы с именем БД и коллекций
const (
	database = "data"
	posts    = "posts"
	authors  = "authors"
)

// Структура объекта хранилища.
type Store struct {
	Client       *mongo.Client
	lastPostID   int
	lastAuthorID int
	idMute       *sync.Mutex
}

// Дополнительная струтура для работы с коллекцией авторов
type author struct {
	ID   int `bson:"_id"`
	name string
}

// Конструктор, принимает строку подключения к БД.
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
		Client:       client,
		lastPostID:   0,
		lastAuthorID: 0,
		idMute:       &sync.Mutex{},
	}
	return &s, nil
}

// Генератор тестовых данных
func (s *Store) TestData() error {
	return nil
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

// Добавление нового автора в хранилище
func (s *Store) AddAuthor(name string) (int, error) {
	s.idMute.Lock()
	defer s.idMute.Unlock()
	s.lastAuthorID++
	colAuthors := s.Client.Database(database).Collection(authors)
	doc := author{ID: s.lastAuthorID, name: name}
	result, err := colAuthors.InsertOne(context.TODO(), doc)
	if err != nil {
		s.lastAuthorID--
		return 0, err
	}
	return result.InsertedID.(int), nil
}
