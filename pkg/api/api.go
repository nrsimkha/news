package api

import (
	"encoding/json"
	storage "goNews/pkg/db"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type API struct {
	r  *mux.Router
	db storage.Interface
}

// Конструтор API
func New(db storage.Interface) *API {
	api := API{
		db: db,
	}
	api.r = mux.NewRouter()

	api.endpoints()
	return &api
}

// Router возвращает маршрутизатор запросов.
func (api *API) Router() *mux.Router {
	return api.r
}

func (api *API) endpoints() {
	api.r.HandleFunc("/news/{amount}", api.newsHandler).Methods(http.MethodGet)
	api.r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./webapp"))))
}

func (api *API) newsHandler(w http.ResponseWriter, r *http.Request) {
	// Считывание параметра {amount} из пути запроса.
	s := mux.Vars(r)["amount"]
	amount, err := strconv.Atoi(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Получение данных из БД.
	news, err := api.db.News(amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Отправка данных клиенту в формате JSON.
	json.NewEncoder(w).Encode(news)
}
