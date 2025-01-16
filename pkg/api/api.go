package api

import (
	"encoding/json"
	storage "goNews/pkg/db"
	"goNews/pkg/logger"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type API struct {
	r  *mux.Router
	db storage.Interface
}

const LIMIT = 10

// Конструтор API
func New(db storage.Interface) *API {
	api := API{
		db: db,
	}
	api.r = mux.NewRouter()
	api.endpoints()
	api.r.Use(logger.WrapHandlerWithLogging)
	return &api
}

// Router возвращает маршрутизатор запросов.
func (api *API) Router() *mux.Router {
	return api.r
}

func (api *API) endpoints() {
	api.r.HandleFunc("/news/id/{id}", api.newsByIDHandler).Methods(http.MethodGet)
	api.r.HandleFunc("/news", api.newsHandler).Methods(http.MethodGet)
	api.r.HandleFunc("/news", api.addNewsHandler).Methods(http.MethodPost)
	api.r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./webapp"))))
}

func (api *API) newsHandler(w http.ResponseWriter, r *http.Request) {

	keyString := r.URL.Query().Get("keystring")
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}
	// Получение данных из БД.
	news, err := api.db.News(page, LIMIT, keyString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправка данных клиенту в формате JSON.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(news)
}

func (api *API) newsByIDHandler(w http.ResponseWriter, r *http.Request) {
	s := mux.Vars(r)["id"]
	id, err := strconv.Atoi(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Получение данных из БД.
	post, err := api.db.NewsByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Отправка данных клиенту в формате JSON.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

func (api *API) addNewsHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("failed to read response body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var news []storage.Post
	err = json.Unmarshal(body, &news)
	if err != nil {
		log.Printf("failed to marshal response: %v", err)
		http.Error(w, "Invalid body structure", http.StatusBadRequest)
		return
	}
	id, err := api.db.AddNews(news)
	log.Print("Получили последний ID: ", *id)
	if err != nil {
		log.Printf("failed to add news to db: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Отправка данных клиенту в формате JSON.
	json.NewEncoder(w).Encode(struct {
		Id *int
	}{Id: id})
}
