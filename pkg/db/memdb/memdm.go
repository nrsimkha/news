package memdb

import (
	storage "goNews/pkg/db"
	"time"
)

// Пользовательский тип данных - реализация БД в памяти.
// Т.н. "заглушка".
//type DB []postgres.Post

// Хранилище данных.
type DB struct{}

// Конструктор объекта хранилища.
func New() *DB {
	return new(DB)
}

func (db *DB) News(int) ([]storage.Post, error) {
	return news, nil
}

func (db *DB) AddNews([]storage.Post) error {
	return nil
}

var news = []storage.Post{
	{ID: 1,
		Title:   "Go – компилируемый, многопоточный язык программирования",
		Content: "Одно из ключевых направлений деятельности нашей компании — это аутсорс-разработка цифровых продуктов. При создании очередной системы мы хотим уделять больше времени и сил необходимым фичам для клиентов, а не настройке рутинного взаимодействия с юзерами, и для ускорения проработки основных пользовательских сценариев мы используем технологию",
		PubTime: time.Now().Unix(),
		Link:    "https: //habr.com/ru/companies/kts/articles/870454/?utm_campaign=870454&amp;utm_source=habrahabr&amp;utm_medium=rss",
	},
	{
		ID:      2,
		Title:   "The Go Memory Model",
		Content: "The Go memory model specifies the conditions under which reads of a variable in one goroutine can be guaranteed to observe values produced by writes to the same variable in a different goroutine.",
		Link:    "https: //habr.com/ru/companies/kts/articles/870454/?utm_campaign=870454&amp;utm_source=habrahabr&amp;utm_medium=rss",
	},
}
