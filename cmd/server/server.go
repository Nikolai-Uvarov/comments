// сервер сервиса комментариев
package main

import (
	"comments/pkg/api"
	"comments/pkg/db/postgres"
	"log"
	"net/http"
)

func main() {

	// Инициализация БД
	//db := dbmock.New()
	db := postgres.New()

	// Создание объекта API, использующего БД.
	api := api.New(db)

	// Запуск сетевой службы и HTTP-сервера
	// на всех локальных IP-адресах на порту 8080.
	err := http.ListenAndServe(":8080", api.Router())
	if err != nil {
		log.Fatal(err)
	}
}
