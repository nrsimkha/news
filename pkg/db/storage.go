package storage

// Публикация в БД.
type Post struct {
	ID      uint64 // id публикации в бд
	Title   string // заголовок публикации
	Content string // содержание публикации
	PubTime int64  // время публикации
	Link    string // ссылка на источник
}

// Interface задаёт контракт на работу с БД.
type Interface interface {
	News(int) ([]Post, error) // получение всех публикаций
	AddNews([]Post) error     // создание новой публикации
}
