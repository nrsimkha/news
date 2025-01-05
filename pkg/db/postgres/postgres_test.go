package postgres

import (
	storage "goNews/pkg/db"
	"log"
	"testing"
	"time"
)

func TestStorage_Posts(t *testing.T) {
	store, err := New("postgres://postgres:0773@localhost:5432/goNews")
	if err != nil {
		t.Errorf("Couldn't connect to db")
	}
	//Провека подключения
	err = store.db.Ping(ctx)
	if err != nil {
		log.Fatal(err)
	}
	p := []storage.Post{{
		Title:   "Go – компилируемый, многопоточный язык программирования",
		Content: "Одно из ключевых направлений деятельности нашей компании — это аутсорс-разработка цифровых продуктов. При создании очередной системы мы хотим уделять больше времени и сил необходимым фичам для клиентов, а не настройке рутинного взаимодействия с юзерами, и для ускорения проработки основных пользовательских сценариев мы используем технологию",
		PubTime: time.Now().Unix(),
		Link:    "https: //habr.com/ru/companies/kts/articles/870454/?utm_campaign=870454&amp;utm_source=habrahabr&amp;utm_medium=rss",
	}}
	err = store.AddNews(p)
	if err != nil {
		t.Errorf("Couldn't add post to db")
		log.Fatal(err)
	}
	posts, err := store.News(10)
	if err != nil {
		t.Errorf("Couldn't get news from db")
		log.Fatal(err)
	}
	t.Logf("%+v", posts)

}
