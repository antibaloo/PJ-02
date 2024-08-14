package main

import (
	"GoNews/pkg/api"
	"GoNews/pkg/storage"
	"GoNews/pkg/storage/mongodb"
	"context"
	"log"
	"net/http"
)

// Сервер GoNews.
type server struct {
	db  storage.Interface
	api *api.API
}

func main() {
	// Создаём объект сервера.
	var srv server
	conString := "mongodb://localhost:27017/"
	db, err := mongodb.New(conString)
	if err != nil {
		log.Fatal(err)
	}
	db.Client.Database("data").Drop(context.TODO()) // Удаление БД перед запуском
	err = db.TestData()
	if err != nil {
		log.Fatal(err)
	}
	//Освобождаем ресурсы по завершению приложения
	defer db.Client.Disconnect(context.Background())
	// Инициализируем хранилище сервера конкретной БД.
	srv.db = db
	// Создаём объект API и регистрируем обработчики.
	srv.api = api.New(srv.db)
	// Запускаем веб-сервер на порту 8080 на всех интерфейсах.
	// Предаём серверу маршрутизатор запросов,
	// поэтому сервер будет все запросы отправлять на маршрутизатор.
	// Маршрутизатор будет выбирать нужный обработчик.
	http.ListenAndServe(":8080", srv.api.Router())
}
