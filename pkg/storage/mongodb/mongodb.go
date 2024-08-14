package mongodb

import (
	"GoNews/pkg/storage"
	"context"
	"fmt"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Константы с именем БД и коллекций
const (
	database   = "data"
	postsCol   = "posts"
	authorsCol = "authors"
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
	ID         int    `bson:"_id"`
	AuthorName string `bson:"author_name"`
}

// Конструктор, принимает строку подключения к БД.
func New(conStr string) (*Store, error) {
	mongoOpts := options.Client().ApplyURI(conStr)
	client, err := mongo.Connect(context.TODO(), mongoOpts)
	if err != nil {
		return &Store{}, err
	}
	err = client.Ping(context.TODO(), nil)
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

	id, err := s.AddAuthor("Joe")
	if err != nil {
		return err
	}
	err = s.AddPost(storage.Post{
		Title:    "First post title",
		Content:  "First post content",
		AuthorID: id,
	})
	if err != nil {
		return err
	}
	time.Sleep(5 * time.Second)
	id, err = s.AddAuthor("Simon")
	if err != nil {
		return err
	}
	err = s.AddPost(storage.Post{
		Title:    "Second post title",
		Content:  "Second post content",
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
	colPosts := s.Client.Database(database).Collection(postsCol)
	options := options.Find()
	//Дополлнительные параметры .Find() взяты из примеров
	res, err := colPosts.Find(context.TODO(), bson.M{}, options)
	if err != nil {
		fmt.Println(err)
		return posts, err
	}
	for res.Next(context.TODO()) {
		var post storage.Post
		err = res.Decode(&post)
		if err != nil {
			fmt.Println(err)
			return posts, err
		}
		posts = append(posts, post)
	}
	if err = res.Err(); err != nil {
		fmt.Println(err)
		return posts, err
	}
	res.Close(context.TODO())
	colAuthors := s.Client.Database(database).Collection(authorsCol)
	for i, post := range posts {
		var a author
		filter := bson.M{"_id": post.AuthorID}
		err = colAuthors.FindOne(context.TODO(), filter).Decode(&a)
		if err != nil {
			return posts, err
		}
		posts[i].AuthorName = a.AuthorName
	}
	return posts, nil
}

// Добавление поста
func (s *Store) AddPost(post storage.Post) error {
	s.idMute.Lock()
	defer s.idMute.Unlock()
	s.lastPostID++
	colPosts := s.Client.Database(database).Collection(postsCol)
	doc := storage.Post{
		ID:        s.lastPostID,
		Title:     post.Title,
		Content:   post.Content,
		AuthorID:  post.AuthorID,
		CreatedAt: time.Now().Unix(),
	}
	_, err := colPosts.InsertOne(context.TODO(), doc)
	if err != nil {
		s.lastPostID--
		return err
	}
	return nil
}

// Обновление поста
func (s *Store) UpdatePost(post storage.Post) error {
	colPosts := s.Client.Database(database).Collection(postsCol)
	filter := bson.M{"_id": post.ID}
	update := bson.M{
		"$set": bson.M{
			"title":        post.Title,
			"content":      post.Content,
			"author_id":    post.AuthorID,
			"published_at": post.PublishedAt,
		},
	}
	_, err := colPosts.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

// Удаление поста
func (s *Store) DeletePost(post storage.Post) error {
	colPosts := s.Client.Database(database).Collection(postsCol)
	filter := bson.M{"_id": post.ID}
	_, err := colPosts.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	return nil
}

// Добавление нового автора в хранилище
func (s *Store) AddAuthor(name string) (int, error) {
	s.idMute.Lock()
	defer s.idMute.Unlock()
	s.lastAuthorID++
	colAuthors := s.Client.Database(database).Collection(authorsCol)
	doc := author{ID: s.lastAuthorID, AuthorName: name}
	_, err := colAuthors.InsertOne(context.TODO(), doc)
	if err != nil {
		s.lastAuthorID--
		return 0, err
	}
	return s.lastAuthorID, nil
}
