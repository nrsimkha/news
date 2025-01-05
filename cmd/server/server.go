package main

import (
	"encoding/json"
	"goNews/pkg/api"
	storage "goNews/pkg/db"
	db "goNews/pkg/db/postgres"
	"goNews/pkg/rss"
	"log"
	"net/http"
	"os"
	"time"
)

type configSettings struct {
	RssUrl []string `json:"rss"`
	Period int      `json:"request_period"`
}

func main() {

	// Реляционная БД PostgreSQL.
	database, err := db.New("postgres://postgres:0773@localhost:5432/goNews")
	if err != nil {
		log.Fatal(err)
	}

	// Создание объекта API, использующего БД в памяти.
	api := api.New(database)

	// Чтение файла с конфигом
	fileBytes, _ := os.ReadFile("./config.json")
	var config configSettings
	json.Unmarshal(fileBytes, &config)

	// Канал для обработки постов новостей
	postsChan := make(chan []storage.Post)
	// Канал для обработки получаемых ошибок
	errChan := make(chan error)

	//Парсинг rss потоков по каждому url
	for _, url := range config.RssUrl {
		go getRssData(url, postsChan, errChan, config.Period)
	}

	// запись потока новостей в БД
	go func() {
		for posts := range postsChan {
			database.AddNews(posts)
		}
	}()

	// обработка потока ошибок
	go func() {
		for err := range errChan {
			log.Println("ошибка:", err)
		}
	}()
	// Запуск сетевой службы и HTTP-сервера
	// на всех локальных IP-адресах на порту 80.
	err = http.ListenAndServe(":80", api.Router())
	if err != nil {
		log.Fatal(err)
	}

}

func getRssData(url string, news chan []storage.Post, errors chan error, period int) {
	ticker := time.NewTicker(time.Duration(period) * time.Second)
	for range ticker.C {
		posts, err := rss.ParseRss(url)
		if err != nil {
			errors <- err
		} else {
			news <- posts
		}

	}
}
