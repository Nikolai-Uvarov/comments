package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"comments/pkg/db/obj"

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
}

// HeadersMiddleware устанавливает заголовки ответа сервера.
func (api *API) HeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Manual-Header", "I love you")
		next.ServeHTTP(w, r)
	})
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
	// Отправка данных клиенту в формате JSON.
	json.NewEncoder(w).Encode(comments)
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

	w.WriteHeader(http.StatusOK)
}
