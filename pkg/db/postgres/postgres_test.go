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
		Link:    "https: //habr.com/ru/companies/kts/articles/870454/?utm_campaign=870pp454&amp;utm_fffsource=habrahabr&amp;utm_medium=rssfdgdd",
	}}
	t.Log("Тестируем функцию добавления новостей в БД")
	id, err := store.AddNews(p)
	t.Log("Получили последний ID: ", *id)
	if err != nil {
		t.Errorf("Couldn't add post to db")
		log.Fatal(err)
	}
	t.Log("Тестируем функцию получения всех новостей из БД")
	posts, err := store.News(1, 3, "")
	if err != nil {
		t.Errorf("Couldn't get news from db")
		log.Fatal(err)
	}
	t.Logf("%+v", posts)
	t.Log("Тестируем функцию получения новостей из БД, с фильтром по названию новости")
	posts, err = store.News(1, 3, "Горутины")
	if err != nil {
		t.Errorf("Couldn't get news from db2")
		log.Fatal(err)
	}
	t.Logf("%+v", posts)
	t.Log("Тестируем функцию получения новости по ID")
	post, err := store.NewsByID(*id)
	t.Logf("%+v", post)
	if err != nil {
		t.Errorf("Couldn't get post from db")
		log.Fatal(err)
	}
	t.Logf("%+v", post)

}
