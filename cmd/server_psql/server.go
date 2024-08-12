package main

import (
	"GoNews/pkg/api"
	"GoNews/pkg/storage"
	"GoNews/pkg/storage/postgres"
	"log"
	"net/http"
	"os"
)

// Сервер GoNews.
type server struct {
	db  storage.Interface
	api *api.API
}

func main() {
	// Создаём объект сервера.
	var srv server
	// Реляционная БД PostgreSQL.
	pwd := os.Getenv("postgrespass")
	conString := "postgres://postgres:" + pwd + "@localhost:5432/gonews?sslmode=disable"
	db, err := postgres.New(conString)
	if err != nil {
		log.Fatal(err)
	}
	err = db.TestData()
	if err != nil {
		log.Fatal(err)
	}
	//Освобождаем ресурсы по завершению приложения
	defer db.Pool.Close()
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
