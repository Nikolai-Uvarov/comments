//api сервера сервиса комментариев
package api

import (
	"comments/pkg/db/obj"
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// API приложения.
type API struct {
	r  *mux.Router // маршрутизатор запросов
	db obj.DB      // база данных
}

// Конструктор API.
func New(db obj.DB) *API {
	api := API{}
	api.db = db
	api.r = mux.NewRouter()
	api.endpoints()
	return &api
}

// Router возвращает маршрутизатор запросов.
func (api *API) Router() *mux.Router {
	return api.r
}

// Регистрация методов API в маршрутизаторе запросов.
func (api *API) endpoints() {
	api.r.HandleFunc("/comments", api.commentsByPost).Methods(http.MethodGet)
	api.r.HandleFunc("/add", api.addComment).Methods(http.MethodPost)
	api.r.Use(api.HeadersMiddleware)
	api.r.Use(api.RequestIDMiddleware)
	api.r.Use(api.LoggingMiddleware)
}


// commentsByPost возвращает все комментарии по id новости.
func (api *API) commentsByPost(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("postID")

	id, err := strconv.Atoi(idParam)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Получение данных из БД.
	comments, err := api.db.GetComments(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var ans = struct {
		Comments  []obj.Comment
		RequestID any
	}{
		Comments:  comments,
		RequestID: getRequestID(r.Context()),
	}
	// Отправка данных клиенту в формате JSON.
	json.NewEncoder(w).Encode(ans)
}

// addComment создает новый комментарий. В теле request должен быть указан id новости или комментария
func (api *API) addComment(w http.ResponseWriter, r *http.Request) {
	var c obj.Comment
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	err = api.db.SaveComment(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var ans = struct {
		RequestID any
	}{
		RequestID: getRequestID(r.Context()),
	}
	
	json.NewEncoder(w).Encode(ans)

	w.WriteHeader(http.StatusOK)
}


func getRequestID(ctx context.Context) any {
	return ctx.Value(contextKey("requestID"))
}
