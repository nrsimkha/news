package storage

// Публикация в БД.
type Post struct {
	ID      uint64 // id публикации в бд
	Title   string // заголовок публикации
	Content string // содержание публикации
	PubTime int64  // время публикации
	Link    string // ссылка на источник
}

// Объект пагинации
type Pagination struct {
	TotalPages   int
	CurrentPage  int
	PostsPerPage int
}

// Объект постов с пагинацией
type PostsWithPagination struct {
	Posts      []Post
	Pagination Pagination
}

// Interface задаёт контракт на работу с БД.
type Interface interface {
	News(int, int, string) (*PostsWithPagination, error) // получение всех публикаций
	AddNews([]Post) (*int, error)                        // создание новой публикации
	NewsByID(int) (*Post, error)                         // получение развернутой публикации по ID
}
