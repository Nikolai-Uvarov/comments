package api

import (
	"bytes"
	"comments/pkg/db/dbmock"
	"comments/pkg/db/obj"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPI_commentsByPost(t *testing.T) {
	db := dbmock.New()

	for _, c := range comments {
		db.SaveComment(c)
	}

	api := New(db)

	req := httptest.NewRequest(http.MethodGet, "/comments?postID=2", nil)

	// Создаём объект для записи ответа обработчика.
	rr := httptest.NewRecorder()
	// Вызываем маршрутизатор. Маршрутизатор для пути и метода запроса
	// вызовет обработчик. Обработчик запишет ответ в созданный объект.
	api.r.ServeHTTP(rr, req)
	// Проверяем код ответа.
	if !(rr.Code == http.StatusOK) {
		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusOK)
	}
	// Читаем тело ответа.
	b, err := io.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("не удалось раскодировать ответ сервера: %v", err)
	}
	// Раскодируем JSON в массив заказов.
	var data []obj.Comment
	err = json.Unmarshal(b, &data)
	if err != nil {
		t.Fatalf("не удалось раскодировать ответ сервера: %v", err)
	}
	// Проверяем, что в массиве комментарии к нужной новости.
	if data[0].PostID != 2 {
		t.Fatalf("получен id %d , ожидалось %s", data[0].PostID, "2")
	}

}

// комменты для тестов
var comments = []obj.Comment{
	{
		ID:        0,
		PostID:    1,
		CommentID: 0,
		Text:      "коммент к 1 первый",
		Answers:   []obj.Comment{},
	},
	{
		ID:        0,
		PostID:    0,
		CommentID: 1,
		Text:      "коммент к комменту 1",
		Answers:   []obj.Comment{},
	},
	{
		ID:        0,
		PostID:    1,
		CommentID: 0,
		Text:      "коммент к 1 второй",
		Answers:   []obj.Comment{},
	},
	{
		ID:        0,
		PostID:    2,
		CommentID: 0,
		Text:      "коммент к 2 первый",
		Answers:   []obj.Comment{},
	},
	{
		ID:        0,
		PostID:    0,
		CommentID: 2,
		Text:      "коммент к комменту 2",
		Answers:   []obj.Comment{},
	},
}

func TestAPI_addComment(t *testing.T) {

	db := dbmock.New()
	api := New(db)

	b, _ := json.Marshal(comments[0])
	buf := bytes.NewBuffer(b)

	req := httptest.NewRequest(http.MethodPost, "/add", buf)

	// Создаём объект для записи ответа обработчика.
	rr := httptest.NewRecorder()
	// Вызываем маршрутизатор. Маршрутизатор для пути и метода запроса
	// вызовет обработчик. Обработчик запишет ответ в созданный объект.
	api.r.ServeHTTP(rr, req)

	// Проверяем код ответа.
	if !(rr.Code == http.StatusOK) {
		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusOK)
	}
	// Проверяем, что в базе сохранен комментарий к нужной новости.
	c, _:=db.GetComments(1)
	
	if c[0].PostID != 1 {
		t.Fatalf("получен id %d , ожидалось %s", c[0].PostID, "1")
	}
}
